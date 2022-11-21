package domain

import (
	"context"
	"crypto/sha1"
	"strconv"
	"strings"

	"github.com/patrickmn/go-cache"
)

type ScbEnetResponse struct {
	Title       string      `json:"title"`
	Status      int         `json:"status"`
	Description string      `json:"description"`
	Result      interface{} `json:"result"`
}
type ScbEnetSignIn struct {
	SessionId string `json:"sessionId"`
}
type ScbEnetTransaction struct {
	Date        string   `json:"date"`
	Time        string   `json:"time"`
	Transaction string   `json:"transaction"`
	Channel     string   `json:"channel"`
	Withdrawal  *float64 `json:"withdrawal"`
	Deposits    *float64 `json:"deposits"`
	Description string   `json:"description"`
}

type ScbEnetAccountBalance struct {
	AccountNo           string   `json:"accountNo"`
	AccountName         string   `json:"accountName"`
	AccountBalance      *float64 `json:"accountBalance"`
	AvailableBalance    *float64 `json:"availableBalance"`
	OverdraftAccount    *float64 `json:"overdraftAccount"`
	AccruedInterest     *float64 `json:"accruedInterest"`
	LastTransactionDate string   `json:"lastTransactionDate"`
}
type ScbEnetLoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ScbEnetUseCase interface {
	SignIn(ctx context.Context, Dto *ScbEnetLoginDto, cache *cache.Cache) (*ScbEnetResponse, error)
	GetTransaction(ctx context.Context, Dto *ScbEnetLoginDto, cache *cache.Cache) (*ScbEnetResponse, error)
	GetAccountBalance(ctx context.Context, Dto *ScbEnetLoginDto, cache *cache.Cache) (*ScbEnetResponse, error)
}

func SubString(Text string, SubText string) string {
	return strings.TrimSpace(Text)[strings.Index(strings.TrimSpace(Text), SubText)+1 : len([]rune(strings.TrimSpace(Text)))]
}

func HashSha1(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	return string(h.Sum(nil))
}

func StringConvertToFloat64(text string) (*float64, error) {
	var number *float64
	res, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return nil, err
	}
	number = &res
	return number, nil
}
