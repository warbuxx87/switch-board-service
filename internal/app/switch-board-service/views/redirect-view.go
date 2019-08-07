package views

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/labstack/echo"
	"net/http"
)

func RedirectView(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	return cc.Context.Redirect(http.StatusFound, "/dashboard")
}