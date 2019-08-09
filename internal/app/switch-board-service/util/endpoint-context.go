package util

import (
	"github.com/labstack/echo/v4"
)

type EndpointContext struct {
	echo.Context
	AppContext *DataContext
}