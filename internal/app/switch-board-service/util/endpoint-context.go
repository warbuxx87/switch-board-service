package util

import (
	"github.com/labstack/echo"
)

type EndpointContext struct {
	echo.Context
	AppContext *DataContext
}