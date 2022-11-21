package scb_enet

import (
	"SCBEasyNetScraper/domain"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/html"
)

type useCase struct {
}

func NewUseCase() domain.ScbEnetUseCase {
	return &useCase{}
}

func CacheManage(caching *cache.Cache, key string, data interface{}, cmd string) (interface{}, bool) {
	key = domain.HashSha1(key)
	switch cmd {
	case "set":
		caching.Set(key, data, cache.DefaultExpiration)
	case "get":
		if res, error := caching.Get(key); error {
			return res, error
		}
	case "del":
		caching.Delete(key)
		<-time.NewTimer(30 * time.Second).C
	}
	// fmt.Println("status:", cmd, "key:", key, "data:", data)
	return data, false
}

func (u *useCase) SignIn(ctx context.Context, dto *domain.ScbEnetLoginDto, caching *cache.Cache) (*domain.ScbEnetResponse, error) {
	if res, err := CacheManage(caching, dto.Username+dto.Password, "", "get"); err {
		if res != domain.SignInTokenFailed {
			resp := &domain.ScbEnetResponse{
				Title: domain.SignIn, Status: http.StatusOK, Description: domain.Success, Result: domain.ScbEnetSignIn{
					SessionId: res.(string),
				},
			}
			return resp, nil
		}
	}
	req, err := http.NewRequest("POST", "https://www.scbeasy.com/online/easynet/page/lgn/login.aspx", strings.NewReader("LOGIN="+dto.Username+"&PASSWD="+dto.Password))
	client := &http.Client{}
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if strings.Contains(string(body), "ERR_CASE") {
		CacheManage(caching, dto.Username+dto.Password, domain.SignInTokenFailed, "set")
		return nil, domain.ErrorSignInTokenFailed.SetMessage(domain.SignIn)
	}
	reader := strings.NewReader(string(body))
	tokenizer := html.NewTokenizer(reader)
	for {
		tokenizer.Next()
		tag, hasAttr := tokenizer.TagName()
		if string(tag) == "input" { //get first tag input sessionId
			if hasAttr {
				for {
					attrKey, attrValue, moreAttr := tokenizer.TagAttr()
					if string(attrKey) == "value" {
						res := &domain.ScbEnetResponse{
							Title: domain.SignIn, Status: http.StatusOK, Description: domain.Success, Result: domain.ScbEnetSignIn{
								SessionId: string(attrValue),
							},
						}
						CacheManage(caching, dto.Username+dto.Password, string(attrValue), "set")
						return res, nil
					}
					if !moreAttr {
						return nil, domain.ErrorSignInTokenFailed.SetMessage(domain.SignIn)
					}
				}
			}
		}
	}
}
func (u *useCase) GetAccountBalance(ctx context.Context, dto *domain.ScbEnetLoginDto, caching *cache.Cache) (*domain.ScbEnetResponse, error) {
	si, err := u.SignIn(ctx, dto, caching)
	if err != nil {
		return nil, err
	}
	if si.Result == domain.SignInTokenFailed {
		return nil, domain.ErrorSignInTokenFailed.SetMessage(domain.SignIn)
	}
	var scbsi domain.ScbEnetSignIn
	respsi, err := json.Marshal(si.Result)
	if err != nil {
		return nil, domain.ErrorSignInTokenFailed.SetMessage(domain.SignIn)
	}
	json.Unmarshal([]byte(respsi), &scbsi)

	var res *domain.ScbEnetResponse
	goColly := colly.NewCollector()
	goColly.OnResponse(func(resgc *colly.Response) {
		if strings.Contains(string(resgc.Body), "Account Balance") {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(resgc.Body)))
			scab := domain.ScbEnetAccountBalance{}
			if err == nil {
				doc.Find("tr td.bd_en_blk11").Each(func(eid int, td *goquery.Selection) {
					switch eid {
					case 0:
						scab.AccountNo = domain.SubString(td.Text(), " X")
					case 1:
						scab.AccountName = strings.TrimSpace(td.Text())
					}
				})
				doc.Find("tr td.hd_th_blk11_bld").Each(func(eid int, td *goquery.Selection) {
					switch eid {
					case 1:
						res, err := domain.StringConvertToFloat64(domain.SubString(td.Text(), "+"))
						if err != nil {
							scab.AccountBalance = nil
						}
						scab.AccountBalance = res
					case 3:
						res, err := domain.StringConvertToFloat64(domain.SubString(td.Text(), "+"))
						if err != nil {
							scab.AvailableBalance = nil
						}
						scab.AvailableBalance = res
					case 5:
						res, err := domain.StringConvertToFloat64(domain.SubString(td.Text(), "+"))
						if err != nil {
							scab.OverdraftAccount = nil
						}
						scab.OverdraftAccount = res
					case 7:
						res, err := domain.StringConvertToFloat64(domain.SubString(td.Text(), "+"))
						if err != nil {
							scab.AccruedInterest = nil
						}
						scab.AccruedInterest = res
					case 9:
						scab.LastTransactionDate = strings.TrimSpace(td.Text())
					}
				})
				if resgc.StatusCode == http.StatusOK {
					if scab.AccountNo != "" && scab.AccountName != "" {
						res = &domain.ScbEnetResponse{
							Title: domain.GetAccountBalance, Result: scab, Status: http.StatusOK, Description: domain.Success,
						}
					} else {
						res = &domain.ScbEnetResponse{
							Title: domain.GetAccountBalance, Result: nil, Status: http.StatusNotFound, Description: domain.NotFound,
						}
					}
				} else {
					res = &domain.ScbEnetResponse{
						Title: domain.GetAccountBalance, Result: nil, Status: http.StatusNotFound, Description: domain.Failed,
					}
				}
			} else if strings.Contains(string(resgc.Body), "login to SCB Easy Net again") {
				CacheManage(caching, dto.Username+dto.Password, "", "del")
				res, _ = u.GetTransaction(ctx, dto, caching)
			} else if strings.Contains(string(resgc.Body), "SCB Maintenance Page") {
				CacheManage(caching, dto.Username+dto.Password, "", "del")
				res, _ = u.GetTransaction(ctx, dto, caching)
			} else {
				res = &domain.ScbEnetResponse{
					Title: domain.GetAccountBalance, Result: nil, Status: http.StatusBadRequest, Description: domain.Suspend,
				}
			}
		}
	})
	goColly.Request("POST",
		"https://www.scbeasy.com/online/easynet/page/acc/acc_bnk_bln.aspx",
		strings.NewReader("SESSIONEASY="+scbsi.SessionId),
		nil,
		http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}})
	if res == nil {
		CacheManage(caching, dto.Username+dto.Password, "", "del")
		res, _ = u.GetAccountBalance(ctx, dto, caching)
	}
	return res, nil
}

func (u *useCase) GetTransaction(ctx context.Context, dto *domain.ScbEnetLoginDto, caching *cache.Cache) (*domain.ScbEnetResponse, error) {
	si, err := u.SignIn(ctx, dto, caching)
	if err != nil {
		return nil, err
	}
	if si.Result == domain.SignInTokenFailed {
		return nil, domain.ErrorSignInTokenFailed.SetMessage(domain.SignIn)
	}
	var scbsi domain.ScbEnetSignIn
	respsi, err := json.Marshal(si.Result)
	if err != nil {
		return nil, domain.ErrorSignInTokenFailed.SetMessage(domain.SignIn)
	}
	json.Unmarshal([]byte(respsi), &scbsi)
	var transaction []domain.ScbEnetTransaction
	var res *domain.ScbEnetResponse
	goColly := colly.NewCollector()
	goColly.OnResponse(func(resgc *colly.Response) {
		if strings.Contains(string(resgc.Body), "Statement - Today") {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(resgc.Body)))
			if err == nil {
				doc.Find("#DataProcess_GridView tbody tr").Each(func(j int, tr *goquery.Selection) {
					scbet := domain.ScbEnetTransaction{}
					tr.Find("td").Each(func(eid int, td *goquery.Selection) {
						switch eid {
						case 0:
							scbet.Date = strings.TrimSpace(td.Text())
						case 1:
							scbet.Time = strings.TrimSpace(td.Text())
						case 2:
							scbet.Transaction = strings.TrimSpace(td.Text())
						case 3:
							scbet.Channel = strings.TrimSpace(td.Text())
						case 4:
							res, err := domain.StringConvertToFloat64(domain.SubString(td.Text(), "-"))
							if err != nil {
								scbet.Withdrawal = nil
							}
							scbet.Withdrawal = res
						case 5:
							res, err := domain.StringConvertToFloat64(domain.SubString(td.Text(), "+"))
							if err != nil {
								scbet.Deposits = nil
							}
							scbet.Deposits = res
						case 6:
							scbet.Description = td.Text()
						}
					})
					if scbet.Transaction != "" {
						transaction = append(transaction, scbet)
					}
				})
				if resgc.StatusCode == http.StatusOK {
					if transaction != nil {
						res = &domain.ScbEnetResponse{
							Title: domain.GetTransaction, Result: transaction, Status: http.StatusOK, Description: domain.Success,
						}
					} else {
						res = &domain.ScbEnetResponse{
							Title: domain.GetTransaction, Result: nil, Status: http.StatusNotFound, Description: domain.NotFound,
						}
					}
				} else {
					res = &domain.ScbEnetResponse{
						Title: domain.GetTransaction, Result: nil, Status: http.StatusNotFound, Description: domain.Failed,
					}
				}
			}
		} else if strings.Contains(string(resgc.Body), "login to SCB Easy Net again") {
			CacheManage(caching, dto.Username+dto.Password, "", "del")
			res, _ = u.GetTransaction(ctx, dto, caching)
		} else if strings.Contains(string(resgc.Body), "SCB Maintenance Page") {
			CacheManage(caching, dto.Username+dto.Password, "", "del")
			res, _ = u.GetTransaction(ctx, dto, caching)
		} else {
			res = &domain.ScbEnetResponse{
				Title: domain.GetTransaction, Result: nil, Status: http.StatusBadRequest, Description: domain.Suspend,
			}
		}
	})
	goColly.Request("POST",
		"https://www.scbeasy.com/online/easynet/page/acc/acc_bnk_tst.aspx",
		strings.NewReader("SESSIONEASY="+scbsi.SessionId),
		nil,
		http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}})

	if res == nil {
		CacheManage(caching, dto.Username+dto.Password, "", "del")
		res, _ = u.GetTransaction(ctx, dto, caching)
	}
	return res, nil
}
