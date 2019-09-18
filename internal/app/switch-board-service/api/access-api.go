package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/forms"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// Validate Access Credentials for Exchanges and Services godoc
// @Summary Validates exchanges and service credentials
// @Description Validates exchanges and service credentials
// @Accept  json
// @Produce  json
// @Tags Public
// @Param exchangeAccessCredentialForm body forms.AccessCredentialsForm true "Exchange Access Credentials"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/validate_access [post]
func ValidateAccessAction(c echo.Context) error {
	cc := c.(*util.EndpointContext)

	exchangeModelSubmitted := new(forms.AccessCredentialsForm)

	jsonBindError := cc.Context.Bind(exchangeModelSubmitted)

	if jsonBindError != nil {
		sentry.CaptureException(jsonBindError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if exchangeModelSubmitted.UUID == "" || exchangeModelSubmitted.AccessToken == "" {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	// Hunt for exchange model access
	var exchangeModels []models.ExchangeModel

	findError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findError != nil {
		sentry.CaptureException(findError)
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
		sentry.CaptureException(serviceKeyFindError)
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

// Generate New Temp Token godoc
// @Summary Generate a new temp token for a exchange
// @Description Generate a new temp token for a exchange
// @Accept  json
// @Produce  json
// @Tags Public
// @Param exchangeAccessCredentialsForm body forms.AccessCredentialsForm true "Exchange Access Credentials"
// @Success 200 {object} responses.ExchangeInfoResponse ""
// @Failure 200 {object} responses.ExchangeInfoResponse ""
// @Router /v1/generate_new_temp_token [post]
func GenerateNewTempTokenAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	accessCredentialsForm := new(forms.AccessCredentialsForm)

	jsonBindError := cc.Context.Bind(accessCredentialsForm)

	if jsonBindError != nil {
		sentry.CaptureException(jsonBindError)
		return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if accessCredentialsForm.UUID == "" || accessCredentialsForm.AccessToken == "" {
		return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	var exchangeModels []models.ExchangeModel

	findError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findError != nil {
		sentry.CaptureException(findError)
		return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: util.StrRetrievingExchangesError, Success: false})
	}

	for _, element := range exchangeModels {

		if element.UUID == accessCredentialsForm.UUID {
			if element.AccessToken == accessCredentialsForm.AccessToken {
				//  access token valid lets generate a new temp token save it and submit

				exchangeModel := new(models.ExchangeModel)

				findError := cc.AppContext.Database.Find(&exchangeModel, "uuid = ?", accessCredentialsForm.UUID).Error

				if findError != nil {
					sentry.CaptureException(findError)
					return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: util.StrExchangeNotFoundError, Success: false})
				}

				exchangeModel.TempAccessToken = uuid.NewV4().String()

				cc.AppContext.Database.Save(&exchangeModel)

				exchangeModel.AccessToken = ""

				return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: "", Success: true, Exchange: *exchangeModel})
			} else {
				// Invalid access token
				return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
			}
		}
	}

	// If we got here the exchange uuid is invalid
	sentry.CaptureMessage(util.StrInvalidExchangeUUIDRequest+accessCredentialsForm.UUID)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
}

// Revoke Temp Token godoc
// @Summary Revokes an exchanges temp access token
// @Description Revokes an exchanges temp access token
// @Accept  json
// @Produce  json
// @Tags Public
// @Param exchangeAccessCredentialsForm body forms.AccessCredentialsForm true "Exchange Access Credentials"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/revoke_temp_access_token [post]
func RevokeTempAccessTokenAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	accessCredentialsForm := new(forms.AccessCredentialsForm)

	jsonBindError := cc.Context.Bind(accessCredentialsForm)

	if jsonBindError != nil {
		sentry.CaptureException(jsonBindError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if accessCredentialsForm.UUID == "" || accessCredentialsForm.AccessToken == "" {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	var exchangeModels []models.ExchangeModel

	findError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findError != nil {
		sentry.CaptureException(findError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRetrievingExchangesError, Success: false})
	}

	for _, element := range exchangeModels {

		if element.UUID == accessCredentialsForm.UUID {
			if element.AccessToken == accessCredentialsForm.AccessToken {
				//  access token valid lets generate a new temp token save it and submit

				exchangeModel := new(models.ExchangeModel)

				findError := cc.AppContext.Database.Find(&exchangeModel, "uuid = ?", accessCredentialsForm.UUID).Error

				if findError != nil {
					sentry.CaptureException(findError)
					return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeNotFoundError, Success: false})
				}

				exchangeModel.TempAccessToken = ""

				cc.AppContext.Database.Save(&exchangeModel)

				exchangeModel.AccessToken = ""

				return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: "", Success: true, Exchange: *exchangeModel})
			} else {
				// Invalid access token
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
			}
		}
	}

	// If we got here the exchange uuid is invalid
	sentry.CaptureMessage(util.StrInvalidExchangeUUIDRequest+accessCredentialsForm.UUID)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
}
