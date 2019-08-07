package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"github.com/getsentry/raven-go"
	uuid "github.com/satori/go.uuid"

	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/labstack/echo"
)

func AddExchangeAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	if loggedError != nil {
		return loggedError
	}

	_ = util.CheckForValidSession(cc.Context)

	exchangeModel := new(models.ExchangeModel)

	jsonBindError := cc.Context.Bind(exchangeModel)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return jsonBindError
	}

	exchangeModel.UUID = uuid.NewV4().String()

	createError := cc.AppContext.Database.Create(&exchangeModel).Error

	if createError != nil {
		raven.CaptureError(createError, nil)
		return createError
	}

	return nil

}

func RemoveExchangeAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	exchangeUUID := cc.Context.Param("exchangeUUID")

	if loggedError != nil {
		return loggedError
	}

	_ = util.CheckForValidSession(cc.Context)

	exchangeModel := new(models.ExchangeModel)

	findError := cc.AppContext.Database.Find(&exchangeModel, "uuid = ?", exchangeUUID).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return findError
	}

	deleteError := cc.AppContext.Database.Unscoped().Delete(&exchangeModel).Error

	if deleteError != nil {
		raven.CaptureError(deleteError, nil)
		return deleteError
	}

	return nil
}
