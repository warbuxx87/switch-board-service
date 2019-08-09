package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/getsentry/raven-go"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func AddServiceKeyAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	if loggedError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})
	}

	_ = util.CheckForValidSession(cc.Context)

	serviceKeyModel := new(models.ServiceKeyModel)

	jsonBindError := cc.Context.Bind(serviceKeyModel)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	serviceKeyModel.UUID = uuid.NewV4().String()
	serviceKeyModel.AccessToken = uuid.NewV4().String()

	createError := cc.AppContext.Database.Create(&serviceKeyModel).Error

	if createError != nil {
		raven.CaptureError(createError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrAddServiceKeyFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})

}

func RemoveServiceKeyAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	loggedError := util.CheckIfLoggedIn(cc.Context)

	if loggedError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrNotCurrentlyLoggedIn, Success: false})

	}

	_ = util.CheckForValidSession(cc.Context)

	serviceKeyModel := new(models.ServiceKeyModel)

	jsonBindError := cc.Context.Bind(serviceKeyModel)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	deleteError := cc.AppContext.Database.Unscoped().Delete(&serviceKeyModel).Error

	if deleteError != nil {
		raven.CaptureError(deleteError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRemoveServiceKeyFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})
}

func ServiceRetrieveAllExchangesAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	serviceKeyModelSubmitted := new(models.ServiceKeyModel)

	jsonBindError := cc.Context.Bind(serviceKeyModelSubmitted)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if serviceKeyModelSubmitted.UUID == "" || serviceKeyModelSubmitted.AccessToken == "" {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	var serviceKeyModels []models.ServiceKeyModel

	findError := cc.AppContext.Database.Find(&serviceKeyModels).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrServiceKeysCouldNotBeRetrieved, Success: false})
	}

	for _, element := range serviceKeyModels {

		if element.UUID == serviceKeyModelSubmitted.UUID {
			if element.AccessToken == serviceKeyModelSubmitted.AccessToken {
				//  access token valid lets generate a new temp token save it and submit

				var exchangeModels []models.ExchangeModel

				findError := cc.AppContext.Database.Find(&exchangeModels).Error

				if findError != nil {
					raven.CaptureError(findError, nil)
					return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRetrievingExchangesError, Success: false})
				}

				return c.JSON(http.StatusOK, &responses.AllExchangesResponse{ErrorMessage: "", Success: true, Exchanges: exchangeModels})
			} else {
				// Invalid access token
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrServiceAccessTokenInvalid, Success: false})
			}
		}
	}

	// If we got here the exchange uuid is invalid
	raven.CaptureMessage("Invalid Service UUID requested: "+serviceKeyModelSubmitted.UUID, nil)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrServiceAccessTokenInvalid, Success: false})
}

func ServiceRetrieveAllTradePairsAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	serviceKeyModelSubmitted := new(models.ServiceKeyModel)

	jsonBindError := cc.Context.Bind(serviceKeyModelSubmitted)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if serviceKeyModelSubmitted.UUID == "" || serviceKeyModelSubmitted.AccessToken == "" {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	var serviceKeyModels []models.ServiceKeyModel

	findError := cc.AppContext.Database.Find(&serviceKeyModels).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrServiceKeysCouldNotBeRetrieved, Success: false})
	}

	for _, element := range serviceKeyModels {

		if element.UUID == serviceKeyModelSubmitted.UUID {
			if element.AccessToken == serviceKeyModelSubmitted.AccessToken {
				//  access token valid lets generate a new temp token save it and submit

				var tradePairModels []models.TradePairModel

				findError := cc.AppContext.Database.Find(&tradePairModels).Error

				if findError != nil {
					raven.CaptureError(findError, nil)
					return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRetrievingTradePairsError, Success: false})
				}

				return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: "", Success: true, TradePairs: tradePairModels})
			} else {
				// Invalid access token
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrServiceAccessTokenInvalid, Success: false})
			}
		}
	}

	// If we got here the exchange uuid is invalid
	raven.CaptureMessage("Invalid Service UUID requested: "+serviceKeyModelSubmitted.UUID, nil)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrServiceAccessTokenInvalid, Success: false})
}
