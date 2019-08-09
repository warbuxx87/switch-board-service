package views

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-contrib/session"
	"net/http"
)

func LoginView(c echo.Context) error {
	cc := c.(*util.EndpointContext)
	_ = cc

	dateSession, dateSessionError := session.Get("date-session", cc.Context)

	if dateSessionError != nil {
		return dateSessionError
	}

	if value, ok := dateSession.Values["loggedIn"]; ok {

		if value == true {

			return cc.Context.Redirect(http.StatusFound, "/dashboard")
		}
	}

	return cc.Context.Render(http.StatusOK, "login.html", map[string]interface{}{})
}
