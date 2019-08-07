package util

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/satori/go.uuid"
	"net/http"
)

func CheckForDefaultUser(database *gorm.DB, config *Configuration) {

	var userModel models.UserModel

	userModelError := database.First(&userModel, 1).Error

	if userModelError != nil {

		database.Create(&models.UserModel{
			UUID:     uuid.NewV4().String(),
			Username:     "admin",
			PasswordHash: "***REMOVED***",
			FirstName:    "admin",
			LastName:     "admin",
			EmailAddress: "admin@admin.com",
		})
	}
}

func CheckIfLoggedIn(c echo.Context) error {
	dateSession, dateSessionError := session.Get("date-session", c)
	if dateSessionError != nil {
		return c.Redirect(http.StatusFound, "/login")
	} else {
		if dateSession.Values["loggedIn"] == true {
			return nil
		} else {
			return c.Redirect(http.StatusFound, "/login")
		}
	}
}

func CheckForValidSession(c echo.Context) error {

	_, dateSessionError := session.Get("date-session", c)

	if dateSessionError != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	return nil
}

func RetrieveUserModelByUUID(database *gorm.DB, uuid string) (*models.UserModel, error) {

	var userModel models.UserModel

	userModelError := database.First(&userModel, "user_uuid = ?", uuid).Error

	if userModelError != nil {
		return nil, userModelError
	}

	return &userModel, nil

}

func RetrieveUserModelByUsername(database *gorm.DB, username string) (*models.UserModel, error) {

	var userModel models.UserModel

	userModelError := database.First(&userModel, "username = ?", username).Error

	if userModelError != nil {
		return nil, userModelError
	}

	return &userModel, nil

}
