package handler

import (
	"SCBEasyNetScraper/domain"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

type scbEnetHandler struct {
	ScbEnetUseCase domain.ScbEnetUseCase
	caching        *cache.Cache
}

func InitscbEnetHandler(e *echo.Echo, ScbEnetUseCase domain.ScbEnetUseCase, caching *cache.Cache) {
	handler := &scbEnetHandler{
		ScbEnetUseCase,
		caching,
	}
	se := e.Group("/scb-enet")
	se.POST("/sign-in", handler.SighIn)
	se.POST("/transaction", handler.GetTransaction)
	se.POST("/account-balance", handler.GetAccountBalance)

}

func (seh *scbEnetHandler) SighIn(c echo.Context) error {
	var dto domain.ScbEnetLoginDto
	if err := c.Bind(&dto); err != nil {
		return domain.ErrorBindStructFailed.SetMessage(domain.SignIn)
	}
	res, err := seh.ScbEnetUseCase.SignIn(c.Request().Context(), &dto, seh.caching)
	if err != nil {
		return err
	}
	return c.JSON(res.Status, res)
}

func (seh *scbEnetHandler) GetTransaction(c echo.Context) error {
	var dto domain.ScbEnetLoginDto
	if err := c.Bind(&dto); err != nil {
		return domain.ErrorBindStructFailed.SetMessage(domain.GetTransaction)
	}
	res, err := seh.ScbEnetUseCase.GetTransaction(c.Request().Context(), &dto, seh.caching)
	if err != nil {
		return err
	}
	return c.JSON(res.Status, res)
}

func (seh *scbEnetHandler) GetAccountBalance(c echo.Context) error {
	var dto domain.ScbEnetLoginDto
	if err := c.Bind(&dto); err != nil {
		return domain.ErrorBindStructFailed.SetMessage(domain.GetAccountBalance)
	}
	res, err := seh.ScbEnetUseCase.GetAccountBalance(c.Request().Context(), &dto, seh.caching)
	if err != nil {
		return err
	}
	return c.JSON(res.Status, res)
}
