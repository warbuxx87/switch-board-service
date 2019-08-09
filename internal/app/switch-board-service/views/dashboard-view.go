package views

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DashboardView(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	if loggedError != nil {
		return loggedError
	}

	return cc.Context.Render(http.StatusOK, "dashboard.html", map[string]interface{}{})
}