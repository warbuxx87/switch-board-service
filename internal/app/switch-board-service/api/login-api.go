package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/forms"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"crypto/sha256"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"net/http"
)

func LoginAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loginForm := new(forms.LoginForm)

	jsonBindError := cc.Context.Bind(loginForm)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	userModel, userModelError := util.RetrieveUserModelByUsername(cc.AppContext.Database, loginForm.Username)

	if userModelError != nil {
		raven.CaptureError(userModelError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrUserNotFoundError, Success: false})
	}

	if userModel.PasswordHash == fmt.Sprintf("%x", sha256.Sum256([]byte(loginForm.Password))) {

		dateSession, dateSessionError := session.Get("date-session", cc.Context)

		if dateSessionError != nil {
			return dateSessionError
		}

		dateSession.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		dateSession.Values["loggedIn"] = true
		dateSession.Values["userModel"] = userModel
		saveFailed := dateSession.Save(cc.Context.Request(), cc.Context.Response())

		if saveFailed != nil {
			raven.CaptureError(saveFailed, nil)
			return saveFailed
		}

		return cc.Context.Redirect(http.StatusFound, "/dashboard")
	}

	return cc.Context.Redirect(http.StatusFound, "/login")
}

func LogoutAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	dateSession, dateSessionError := session.Get("date-session", cc.Context)

	if dateSessionError != nil {
		return dateSessionError
	}

	dateSession.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	dateSession.Values["loggedIn"] = false
	dateSession.Values["userModel"] = nil
	saveFailed := dateSession.Save(cc.Context.Request(), cc.Context.Response())

	if saveFailed != nil {
		raven.CaptureError(saveFailed, nil)
		return saveFailed
	}

	return cc.Context.Redirect(http.StatusFound, "/login")
}
