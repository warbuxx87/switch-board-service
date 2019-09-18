package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/forms"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Login godoc
// @Summary Login and receive JWT token.
// @Description Login endpoint to receive a JWT token to authorize against protected endpoints.
// @Accept  json
// @Produce  json
// @Tags Public
// @Param user_credentials body forms.LoginForm true "User Credentials"
// @Success 200 {object} responses.LoginResponse ""
// @Failure 200 {object} responses.LoginResponse ""
// @Router /v1/login [post]
func LoginAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loginForm := new(forms.LoginForm)

	jsonBindError := cc.Context.Bind(loginForm)

	if jsonBindError != nil {
		sentry.CaptureException(jsonBindError)
		return c.JSON(http.StatusOK, &responses.LoginResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	userModel, userModelError := util.RetrieveUserModelByUsername(cc.AppContext.Database, loginForm.Username)

	if userModelError != nil {
		sentry.CaptureException(userModelError)
		return c.JSON(http.StatusOK, &responses.LoginResponse{ErrorMessage: util.StrUserNotFoundError, Success: false})
	}
	if err := verifyPassword(loginForm.Password, userModel.PasswordHash); err != nil {
		sentry.CaptureException(err)
		sentry.CaptureMessage(fmt.Sprintf("Invalid authentication credentials; username=%s&password=%s&ip=%s",
			loginForm.Username, loginForm.Password, c.RealIP()))
		return c.JSON(http.StatusOK, &responses.LoginResponse{
			Success:      false,
			ErrorMessage: util.StrInvalidLoginCredential,
		})
	}
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = userModel.FirstName + " " + userModel.LastName
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["userModel"] = userModel

	// Generate encoded token and send it as response.
	createdToken, createdTokenError := token.SignedString([]byte(cc.AppContext.Config.JWTSecret))
	if createdTokenError != nil {
		sentry.CaptureException(errors.Wrap(createdTokenError, "failed to create JWT token"))
		return createdTokenError
	}
	return c.JSON(http.StatusOK, &responses.LoginResponse{ErrorMessage: "", Success: true, AuthToken: "Bearer " + createdToken})
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
// verifyPassword will return nil if inputPassword is same as hashedPassword
func verifyPassword(inputPassword, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err != nil {
		return err
	}
	return nil
}
