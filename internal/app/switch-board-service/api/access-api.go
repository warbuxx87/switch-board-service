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

func ValidateAccessAction(c echo.Context) error {
	cc := c.(*util.EndpointContext)

	exchangeModelSubmitted := new(models.ExchangeModel)

	jsonBindError := cc.Context.Bind(exchangeModelSubmitted)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if exchangeModelSubmitted.UUID == "" || exchangeModelSubmitted.AccessToken == "" {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	// Hunt for exchange model access
	var exchangeModels []models.ExchangeModel

	findError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRetrievingExchangesError, Success: false})
	}

	for _, element := range exchangeModels {

		if element.UUID == exchangeModelSubmitted.UUID {
			if element.AccessToken == exchangeModelSubmitted.AccessToken {
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})

			}
		}
	}

	// Hunt service key models
	var serviceKeyModels []models.ServiceKeyModel

	serviceKeyFindError := cc.AppContext.Database.Find(&serviceKeyModels).Error

	if serviceKeyFindError != nil {
		raven.CaptureError(serviceKeyFindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrServiceKeysCouldNotBeRetrieved, Success: false})
	}

	for _, element := range serviceKeyModels {

		if element.UUID == exchangeModelSubmitted.UUID {
			if element.AccessToken == exchangeModelSubmitted.AccessToken {

				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})

			} else {
				// Invalid access token
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrServiceAccessTokenInvalid, Success: false})
			}
		}
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})

}

func GenerateNewTempTokenAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	exchangeModelSubmitted := new(models.ExchangeModel)

	jsonBindError := cc.Context.Bind(exchangeModelSubmitted)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if exchangeModelSubmitted.UUID == "" || exchangeModelSubmitted.AccessToken == "" {
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
			if element.AccessToken == exchangeModelSubmitted.AccessToken {
				//  access token valid lets generate a new temp token save it and submit

				exchangeModel := new(models.ExchangeModel)

				findError := cc.AppContext.Database.Find(&exchangeModel, "uuid = ?", exchangeModel.UUID).Error

				if findError != nil {
					raven.CaptureError(findError, nil)
					return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeNotFoundError, Success: false})
				}

				exchangeModel.AccessToken = ""
				exchangeModel.TempAccessToken = uuid.NewV4().String()

				cc.AppContext.Database.Save(&exchangeModel)

				return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: "", Success: true, Exchange: *exchangeModel})
			} else {
				// Invalid access token
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
			}
		}
	}

	// If we got here the exchange uuid is invalid
	raven.CaptureMessage("Invalid Exchange UUID requested: "+exchangeModelSubmitted.UUID, nil)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
}

func RevokeTempAccessTokenAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	exchangeModelSubmitted := new(models.ExchangeModel)

	jsonBindError := cc.Context.Bind(exchangeModelSubmitted)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if exchangeModelSubmitted.UUID == "" || exchangeModelSubmitted.AccessToken == "" {
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
			if element.AccessToken == exchangeModelSubmitted.AccessToken {
				//  access token valid lets generate a new temp token save it and submit

				exchangeModel := new(models.ExchangeModel)

				findError := cc.AppContext.Database.Find(&exchangeModel, "uuid = ?", exchangeModel.UUID).Error

				if findError != nil {
					raven.CaptureError(findError, nil)
					return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeNotFoundError, Success: false})
				}

				exchangeModel.TempAccessToken = ""

				cc.AppContext.Database.Save(&exchangeModel)

				return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: "", Success: true, Exchange: *exchangeModel})
			} else {
				// Invalid access token
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
			}
		}
	}

	// If we got here the exchange uuid is invalid
	raven.CaptureMessage("Invalid Exchange UUID requested: "+exchangeModelSubmitted.UUID, nil)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
}
