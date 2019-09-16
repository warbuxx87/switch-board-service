package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/forms"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"

	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/labstack/echo/v4"
)

// Create Exchange godoc
// @Summary Create a new exchange entry.
// @Description Create a new exchange entry.
// @Accept  json
// @Produce  json
// @Tags Exchange
// @Param Authorization header string true "Authentication header"
// @Param createExchangeModel body forms.CreateExchangeForm true "Add Exchange Form"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/action/create_exchange [post]
func CreateExchangeAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	exchangeModel := new(models.ExchangeModel)
	addExchangeForm := new(forms.CreateExchangeForm)

	jsonBindError := cc.Context.Bind(addExchangeForm)

	if jsonBindError != nil {
		sentry.CaptureException(jsonBindError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	exchangeModel.Name = addExchangeForm.Name
	exchangeModel.UUID = uuid.NewV4().String()
	exchangeModel.AccessToken = uuid.NewV4().String()

	createError := cc.AppContext.Database.Create(&exchangeModel).Error

	if createError != nil {
		sentry.CaptureException(createError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrAddExchangeFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})

}

// Remove Exchange Entry godoc
// @Summary Removes a existing exchange entry.
// @Description Removes a existing exchange entry
// @Accept  json
// @Produce  json
// @Tags Exchange
// @Param Authorization header string true "Authentication header"
// @Param exchangeUUID query string true "Exchange UUID to Remove"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/action/remove_exchange [get]
func RemoveExchangeAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	exchangeUUID := cc.Context.QueryParam("exchangeUUID")

	exchangeModel := new(models.ExchangeModel)

	findError := cc.AppContext.Database.Find(&exchangeModel, "uuid = ?", exchangeUUID).Error

	if findError != nil {
		sentry.CaptureException(findError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeNotFoundError, Success: false})
	}

	deleteError := cc.AppContext.Database.Unscoped().Delete(&exchangeModel).Error

	if deleteError != nil {
		sentry.CaptureException(deleteError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRemoveExchangeFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})
}

// Authorize Exchange godoc
// @Summary Authorizes an exchange
// @Description  Authorizes an exchange
// @Accept  json
// @Produce  json
// @Tags Exchange
// @Param Authorization header string true "Authentication header"
// @Param exchangeUUID query string true "Exchange UUID"
// @Param exchangeAuthorized query bool true "Should the exchange be authorized"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/action/authorize_exchange [get]
func AuthorizeExchangeAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	exchangeUUID := cc.Context.QueryParam("exchangeUUID")
	exchangeAuthorized := cc.Context.QueryParam("exchangeAuthorized")

	exchangeModel := new(models.ExchangeModel)

	findError := cc.AppContext.Database.Find(&exchangeModel, "uuid = ?", exchangeUUID).Error

	if findError != nil {
		sentry.CaptureException(findError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeNotFoundError, Success: false})
	}

	convBool, convBoolError := strconv.ParseBool(exchangeAuthorized)

	if convBoolError != nil {
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
	}

	exchangeModel.Authorized = convBool

	cc.AppContext.Database.Save(&exchangeModel)

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})
}

// Retrieve All Exchange Entries godoc
// @Summary Retrieve all exchanges .
// @Description Retrieves all exchanges and exchange information
// @Accept  json
// @Produce  json
// @Tags Exchange
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} responses.AllExchangesResponse ""
// @Failure 200 {object} responses.AllExchangesResponse ""
// @Router /v1/action/retrieve_all_exchanges [get]
func RetrieveAllExchangesAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	var exchangeModels []models.ExchangeModel

	findError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findError != nil {
		sentry.CaptureException(findError)
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
	return c.JSON(http.StatusOK, &responses.AllExchangesResponse{Exchanges: sanitizedExchangeModels, ErrorMessage: "", Success: true})
}

// Retrieve Exchange Info By Exchange Temp Token Entries godoc
// @Summary Retrieve exchange info by temp token .
// @Description Retrieves all exchanges and exchange information
// @Accept  json
// @Produce  json
// @Tags Exchange
// @Param exchangeAccessCredentialsForm body forms.AccessCredentialsForm true "Exchange Access Credentials"
// @Success 200 {object} responses.ExchangeInfoResponse ""
// @Failure 200 {object} responses.ExchangeInfoResponse ""
// @Router /v1/action/retrieve_exchange_by_temp_token [post]
func RetrieveExchangeInfoByExchangeTempTokenAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	accessCredentialsForm := new(forms.AccessCredentialsForm)

	jsonBindError := cc.Context.Bind(accessCredentialsForm)

	if jsonBindError != nil {
		sentry.CaptureException(jsonBindError)
		return c.JSON(http.StatusOK, &responses.ExchangeInfoResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})
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
			if element.TempAccessToken == accessCredentialsForm.AccessToken {
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
	sentry.CaptureMessage(util.StrInvalidExchangeUUIDRequest+accessCredentialsForm.UUID)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
}
