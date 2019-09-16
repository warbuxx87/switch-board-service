package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/forms"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/getsentry/raven-go"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// Create Service Key Entry godoc
// @Summary Create a service key entry .
// @Description Creates a new service entry
// @Accept  json
// @Produce  json
// @Tags Service
// @Param Authorization header string true "Authentication header"
// @Param createServiceKeyModel body forms.CreateServiceForm true "Add Exchange Form"
// @Success 200 {object} responses.ServiceKeyResponse ""
// @Failure 200 {object} responses.ServiceKeyResponse ""
// @Router /v1/action/create_service_key [post]
func CreateServiceKeyAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	serviceKeyForm := new(forms.CreateServiceForm)

	jsonBindError := cc.Context.Bind(serviceKeyForm)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.ServiceKeyResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	var serviceKeyModel models.ServiceKeyModel

	newUUID := uuid.NewV4().String()
	newKey := uuid.NewV4().String()

	serviceKeyModel.UUID = newUUID
	serviceKeyModel.AccessToken = newKey
	serviceKeyModel.Name = serviceKeyForm.Name

	createError := cc.AppContext.Database.Create(&serviceKeyModel).Error

	if createError != nil {
		raven.CaptureError(createError, nil)
		return c.JSON(http.StatusOK, &responses.ServiceKeyResponse{ErrorMessage: util.StrAddServiceKeyFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.ServiceKeyResponse{ErrorMessage: "", Success: true, ServiceUUID: newUUID, ServiceKey: newKey})

}

// Remove Service key Entry godoc
// @Summary Removes a service key entry
// @Description  Removes a service key entry
// @Accept  json
// @Produce  json
// @Tags Service
// @Param Authorization header string true "Authentication header"
// @Param serviceKeyUUID query string true "Service Key UUID to Remove"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/action/remove_service_key [get]
func RemoveServiceKeyAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	serviceKeyUUID := cc.Context.QueryParam("serviceKeyUUID")

	serviceKeyModel := new(models.ServiceKeyModel)

	findError := cc.AppContext.Database.Find(&serviceKeyModel, "uuid = ?", serviceKeyUUID).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeNotFoundError, Success: false})
	}

	deleteError := cc.AppContext.Database.Unscoped().Delete(&serviceKeyModel).Error

	if deleteError != nil {
		raven.CaptureError(deleteError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRemoveServiceKeyFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})
}

// Retrieve All Exchange Entries godoc
// @Summary Retrieve all exchanges .
// @Description Retrieves all exchanges and exchange information
// @Accept  json
// @Produce  json
// @Tags Service-Endpoints
// @Param serviceKeyCredentialForm body forms.AccessCredentialsForm true "Service Key Credentials"
// @Success 200 {object} responses.AllExchangesResponse ""
// @Failure 200 {object} responses.AllExchangesResponse ""
// @Router /v1/service/service_retrieve_all_exchanges [post]
func ServiceRetrieveAllExchangesAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	serviceKeyModelSubmitted := new(forms.AccessCredentialsForm)

	jsonBindError := cc.Context.Bind(serviceKeyModelSubmitted)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.AllExchangesResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if serviceKeyModelSubmitted.UUID == "" || serviceKeyModelSubmitted.AccessToken == "" {
		return c.JSON(http.StatusOK, &responses.AllExchangesResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	var serviceKeyModels []models.ServiceKeyModel

	findError := cc.AppContext.Database.Find(&serviceKeyModels).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.AllExchangesResponse{ErrorMessage: util.StrServiceKeysCouldNotBeRetrieved, Success: false})
	}

	for _, element := range serviceKeyModels {

		if element.UUID == serviceKeyModelSubmitted.UUID {
			if element.AccessToken == serviceKeyModelSubmitted.AccessToken {
				//  access token valid lets generate a new temp token save it and submit

				var exchangeModels []models.ExchangeModel

				findError := cc.AppContext.Database.Find(&exchangeModels).Error

				if findError != nil {
					raven.CaptureError(findError, nil)
					return c.JSON(http.StatusOK, &responses.AllExchangesResponse{ErrorMessage: util.StrRetrievingExchangesError, Success: false})
				}

				var sanitizedExchangeModels []models.ExchangeModel

				for _, element := range exchangeModels {
					element.AccessToken = ""
					var tradePairs []models.TradePairModel
					cc.AppContext.Database.Model(&element).Related(&tradePairs, "AuthorizedTradePairs")
					element.AuthorizedTradePairs = tradePairs
					sanitizedExchangeModels = append(sanitizedExchangeModels, element)
				}

				return c.JSON(http.StatusOK, &responses.AllExchangesResponse{ErrorMessage: "", Success: true, Exchanges: sanitizedExchangeModels})
			} else {
				// Invalid access token
				return c.JSON(http.StatusOK, &responses.AllExchangesResponse{ErrorMessage: util.StrServiceAccessTokenInvalid, Success: false})
			}
		}
	}

	// If we got here the exchange uuid is invalid
	raven.CaptureMessage("Invalid Service UUID requested: "+serviceKeyModelSubmitted.UUID, nil)
	return c.JSON(http.StatusOK, &responses.AllExchangesResponse{ErrorMessage: util.StrServiceAccessTokenInvalid, Success: false})
}

// Retrieve All Trade Pair Entries godoc
// @Summary Retrieve all trade pair entries .
// @Description Retrieves all trade pair entries
// @Accept  json
// @Produce  json
// @Tags Service-Endpoints
// @Param accessCredentialForm body forms.AccessCredentialsForm true "Service Key Credentials"
// @Success 200 {object} responses.AllExchangesResponse ""
// @Failure 200 {object} responses.AllExchangesResponse ""
// @Router /v1/service/service_retrieve_all_trade_pairs [post]
func ServiceRetrieveAllTradePairsAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	serviceKeyModelSubmitted := new(forms.AccessCredentialsForm)

	jsonBindError := cc.Context.Bind(serviceKeyModelSubmitted)

	if jsonBindError != nil {
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	if serviceKeyModelSubmitted.UUID == "" || serviceKeyModelSubmitted.AccessToken == "" {
		return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	var serviceKeyModels []models.ServiceKeyModel

	findError := cc.AppContext.Database.Find(&serviceKeyModels).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: util.StrServiceKeysCouldNotBeRetrieved, Success: false})
	}

	for _, element := range serviceKeyModels {

		if element.UUID == serviceKeyModelSubmitted.UUID {
			if element.AccessToken == serviceKeyModelSubmitted.AccessToken {
				//  access token valid lets generate a new temp token save it and submit

				var tradePairModels []models.TradePairModel

				findError := cc.AppContext.Database.Find(&tradePairModels).Error

				if findError != nil {
					raven.CaptureError(findError, nil)
					return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: util.StrRetrievingTradePairsError, Success: false})
				}

				return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: "", Success: true, TradePairs: tradePairModels})
			} else {
				// Invalid access token
				return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: util.StrServiceAccessTokenInvalid, Success: false})
			}
		}
	}

	// If we got here the exchange uuid is invalid
	raven.CaptureMessage("Invalid Service UUID requested: "+serviceKeyModelSubmitted.UUID, nil)
	return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: util.StrServiceAccessTokenInvalid, Success: false})
}
