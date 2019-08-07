package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"github.com/getsentry/raven-go"
	uuid "github.com/satori/go.uuid"
	"net/http"

	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/labstack/echo"
)

func AddTradePairAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	if loggedError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	_ = util.CheckForValidSession(cc.Context)

	tradePairModel := new(models.TradePairModel)

	jsonBindError := cc.Context.Bind(tradePairModel)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})

	}

	tradePairModel.UUID = uuid.NewV4().String()

	createError := cc.AppContext.Database.Create(&tradePairModel).Error

	if createError != nil {
		raven.CaptureError(createError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrAddTradePairError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})

}

func RemoveTradePairAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	tradePairUUID := cc.Context.Param("tradePairUUID")

	if loggedError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	_ = util.CheckForValidSession(cc.Context)

	tradePairModel := new(models.TradePairModel)

	findError := cc.AppContext.Database.Find(&tradePairModel, "uuid = ?", tradePairUUID).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrTradePairNotFoundError, Success: false})
	}

	var exchangeModels []models.ExchangeModel

	findAllError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findAllError != nil {
		raven.CaptureError(findAllError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangesCouldNotBeRetrieved, Success: false})
	}

	for _, element := range exchangeModels {

		if element.UUID == tradePairModel.UUID {
			return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRemoveTradePairInUseError, Success: false})
		}

	}

	deleteError := cc.AppContext.Database.Unscoped().Delete(&tradePairModel).Error

	if deleteError != nil {
		raven.CaptureError(deleteError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRemoveTradePairFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})
}

func RetrieveAllTradePairs(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)


	if loggedError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	_ = util.CheckForValidSession(cc.Context)

	var tradePairModels []models.TradePairModel

	findError := cc.AppContext.Database.Find(&tradePairModels).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRetrievingTradePairsError, Success: false})
	}


	return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ TradePairs: tradePairModels, ErrorMessage: "", Success: true})
}