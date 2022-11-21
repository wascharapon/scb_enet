package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"

	"SCBEasyNetScraper/app/scb_enet/config"
	"SCBEasyNetScraper/app/scb_enet/handler"
	"SCBEasyNetScraper/app/scb_enet/middleware"
	"SCBEasyNetScraper/module/scb_enet"
)

func init() {

}

func main() {
	c := config.Init()
	e := echo.New()
	caching := cache.New(15*time.Minute, 15*time.Minute)
	seu := scb_enet.NewUseCase()
	e.HTTPErrorHandler = middleware.EchoErrorHandler
	handler.InitscbEnetHandler(e, seu, caching)
	e.Logger.Fatal(e.Start(":" + c.Port))
}
