package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"github.com/getsentry/raven-go"
	uuid "github.com/satori/go.uuid"
	"net/http"

	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/labstack/echo/v4"
)

func AddExchangeAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	if loggedError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	_ = util.CheckForValidSession(cc.Context)

	exchangeModel := new(models.ExchangeModel)

	jsonBindError := cc.Context.Bind(exchangeModel)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	exchangeModel.UUID = uuid.NewV4().String()
	exchangeModel.AccessToken = uuid.NewV4().String()

	createError := cc.AppContext.Database.Create(&exchangeModel).Error

	if createError != nil {
		raven.CaptureError(createError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrAddExchangeFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})

}

func RemoveExchangeAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	exchangeUUID := cc.Context.Param("exchangeUUID")

	if loggedError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	_ = util.CheckForValidSession(cc.Context)

	exchangeModel := new(models.ExchangeModel)

	findError := cc.AppContext.Database.Find(&exchangeModel, "uuid = ?", exchangeUUID).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeNotFoundError, Success: false})
	}

	deleteError := cc.AppContext.Database.Unscoped().Delete(&exchangeModel).Error

	if deleteError != nil {
		raven.CaptureError(deleteError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRemoveExchangeFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})
}

func AuthorizeExchangeAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	if loggedError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	_ = util.CheckForValidSession(cc.Context)

	exchangeModel := new(models.ExchangeModel)
	exchangeModelForm := new(models.ExchangeModel)

	jsonBindError := cc.Context.Bind(exchangeModelForm)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	findError := cc.AppContext.Database.Find(&exchangeModel, "uuid = ?", exchangeModel.UUID).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeNotFoundError, Success: false})
	}

	exchangeModel.Authorized = exchangeModelForm.Authorized

	cc.AppContext.Database.Save(&exchangeModel)

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})
}

func RetrieveAllExchangesAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	if loggedError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	_ = util.CheckForValidSession(cc.Context)

	var exchangeModels []models.ExchangeModel

	findError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRetrievingExchangesError, Success: false})
	}

	var sanitizedExchangeModels []models.ExchangeModel

	for _, element := range exchangeModels {
		element.AccessToken = ""
		var tradePairs []models.TradePairModel
		cc.AppContext.Database.Model(&element).Related(&tradePairs, "AuthorizedTradePairs")
		element.AuthorizedTradePairs = tradePairs
		sanitizedExchangeModels = append(sanitizedExchangeModels, element)
	}

	return c.JSON(http.StatusOK, &responses.AllExchangesResponse{Exchanges: sanitizedExchangeModels, ErrorMessage: "", Success: true})
}

func RetrieveExchangeInfoByExchangeTempTokenAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	exchangeModelSubmitted := new(models.ExchangeModel)

	jsonBindError := cc.Context.Bind(exchangeModelSubmitted)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if exchangeModelSubmitted.UUID == "" || exchangeModelSubmitted.TempAccessToken == "" {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	var exchangeModels []models.ExchangeModel

	findError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRetrievingExchangesError, Success: false})
	}

	for _, element := range exchangeModels {

		if element.UUID == exchangeModelSubmitted.UUID {
			if element.TempAccessToken == exchangeModelSubmitted.TempAccessToken {
				// Temp access token valid lets return the information
				element.AccessToken = ""
				var tradePairs []models.TradePairModel
				cc.AppContext.Database.Model(&element).Related(&tradePairs, "AuthorizedTradePairs")
				element.AuthorizedTradePairs = tradePairs
				return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: "", Success: true, Exchange: element})
			} else {
				// Invalid temp access token
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeTempAccessTokenInvalid, Success: false})
			}
		}
	}

	// If we got here the exchange uuid is invalid

	raven.CaptureMessage("Invalid Exchange UUID requested: "+exchangeModelSubmitted.UUID, nil)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
}
