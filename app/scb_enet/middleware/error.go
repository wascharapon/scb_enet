package middleware

import (
	"SCBEasyNetScraper/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func EchoErrorHandler(err error, c echo.Context) {
	// default error 500 with internal server error
	var title interface{}
	httpCode := http.StatusInternalServerError
	errCode := domain.INTERNAL_SERVER_ERROR
	description := domain.Failed
	title = domain.Close
	appErr, ok := err.(domain.Error)
	if ok {
		errCode = appErr.Status
		title = appErr.Title
		switch appErr.Category {
		case domain.UNAUTHORIZED:
			httpCode = http.StatusUnauthorized
		case domain.FORBIDDEN:
			httpCode = http.StatusForbidden
		case domain.NOT_FOUND:
			description = domain.NotFound
			httpCode = http.StatusNotFound
		case domain.CONFLICT:
			httpCode = http.StatusConflict
		default:
			httpCode = http.StatusBadRequest
		}
	}
	if err := c.JSON(httpCode, domain.ErrorResponse{Title: title, Status: httpCode, Description: description, Result: errCode}); err != nil {
		c.Logger().Error(err)
		return
	}
	c.Logger().Error(err)
}
