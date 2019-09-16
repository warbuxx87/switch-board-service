package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/forms"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"net/http"

	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/util"
	"github.com/labstack/echo/v4"
)

// Create Trade Pair godoc
// @Summary Create a trade pair entry.
// @Description Create a trade pair entry.
// @Accept  json
// @Produce  json
// @Tags TradePair
// @Param Authorization header string true "Authentication header"
// @Param createTradePairModel body forms.CreateTradePairForm true "Trade Pair Form"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/action/create_trade_pair [post]
func CreateTradePairAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	tradePairForm := new(forms.CreateTradePairForm)

	jsonBindError := cc.Context.Bind(tradePairForm)

	if jsonBindError != nil {
		sentry.CaptureException(jsonBindError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})

	}

	var tradePairModel models.TradePairModel

	tradePairModel.UUID = uuid.NewV4().String()
	tradePairModel.Name = tradePairForm.Name
	tradePairModel.BaseAsset = tradePairForm.BaseAsset
	tradePairModel.SecondaryAsset = tradePairForm.SecondaryAsset

	createError := cc.AppContext.Database.Create(&tradePairModel).Error

	if createError != nil {
		sentry.CaptureException(createError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrAddTradePairError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})

}

// Remove Trade Pair Entry godoc
// @Summary Removes a trade pair entry.
// @Description Removes a trade pair entry
// @Accept  json
// @Produce  json
// @Tags TradePair
// @Param Authorization header string true "Authentication header"
// @Param tradePairUUID query string true "Trade Pair UUID"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/action/remove_trade_pair [get]
func RemoveTradePairAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	tradePairUUID := cc.Context.QueryParam("tradePairUUID")

	tradePairModel := new(models.TradePairModel)

	findError := cc.AppContext.Database.Find(&tradePairModel, "uuid = ?", tradePairUUID).Error

	if findError != nil {
		sentry.CaptureException(findError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrTradePairNotFoundError, Success: false})
	}

	var exchangeModels []models.ExchangeModel

	findAllError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findAllError != nil {
		sentry.CaptureException(findAllError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangesCouldNotBeRetrieved, Success: false})
	}

	for _, element := range exchangeModels {

		var tradePairs []models.TradePairModel
		cc.AppContext.Database.Model(&element).Related(&tradePairs, "AuthorizedTradePairs")

		for _, tradePairElement := range tradePairs {
			if tradePairElement.UUID == tradePairModel.UUID {
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRemoveTradePairInUseError, Success: false})
			}

		}

	}

	deleteError := cc.AppContext.Database.Unscoped().Delete(&tradePairModel).Error

	if deleteError != nil {
		sentry.CaptureException(deleteError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRemoveTradePairFailedError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})
}

// Retrieve All Trade Pair Entries godoc
// @Summary Retrieve all trade pairs .
// @Description Retrieves all trade pair entries
// @Accept  json
// @Produce  json
// @Tags TradePair
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} responses.AllTradePairsResponse ""
// @Failure 200 {object} responses.AllTradePairsResponse ""
// @Router /v1/action/retrieve_all_trade_pairs [get]
func RetrieveAllTradePairsAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	var tradePairModels []models.TradePairModel

	findError := cc.AppContext.Database.Find(&tradePairModels).Error

	if findError != nil {
		sentry.CaptureException(findError)
		return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: util.StrRetrievingTradePairsError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{TradePairs: tradePairModels, ErrorMessage: "", Success: true})
}

// Create a  Associate Exchange and Trade Pair godoc
// @Summary Associates a exchange and trade pair
// @Description Associates a exchange and trade pair
// @Accept  json
// @Produce  json
// @Tags TradePair
// @Param Authorization header string true "Authentication header"
// @Param exchangeTradePairAssociationForm body forms.ExchangeTradePairAssociationForm true "Exchange Trade Pair Association"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/action/create_exchange_trade_pair_association [post]
func CreateExchangeTradePairsAssociationAction(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	tradePairAssociationForm := new(forms.ExchangeTradePairAssociationForm)

	jsonBindError := cc.Context.Bind(tradePairAssociationForm)

	if jsonBindError != nil {
		sentry.CaptureException(jsonBindError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})

	}

	var exchangeModels []models.ExchangeModel

	findError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findError != nil {
		sentry.CaptureException(findError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRetrievingExchangesError, Success: false})
	}

	for _, element := range exchangeModels {

		if element.UUID == tradePairAssociationForm.ExchangeUUID {

			var tradePairs []models.TradePairModel
			cc.AppContext.Database.Model(&element).Related(&tradePairs, "AuthorizedTradePairs")
			element.AuthorizedTradePairs = tradePairs

			for _, tradePairElement := range tradePairs {
				if tradePairElement.UUID == tradePairAssociationForm.TradePairUUID {
					return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrTradePairAlreadyAssociatedError, Success: false})
				}
			}

			tradePairModel := new(models.TradePairModel)

			findError := cc.AppContext.Database.Find(&tradePairModel, "uuid = ?", tradePairAssociationForm.TradePairUUID).Error

			if findError != nil {
				sentry.CaptureException(findError)
				return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrTradePairNotFoundError, Success: false})
			}

			element.AuthorizedTradePairs = append(element.AuthorizedTradePairs, *tradePairModel)

			cc.AppContext.Database.Save(&element)

			return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})

		}
	}

	// If we got here the exchange uuid is invalid

	sentry.CaptureMessage("Invalid Exchange UUID requested: "+tradePairAssociationForm.ExchangeUUID)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
}

// Removes a  Associate Exchange and Trade Pair godoc
// @Summary Removes a exchange and trade pair association
// @Description Removes a exchange and trade pair association
// @Accept  json
// @Produce  json
// @Tags TradePair
// @Param Authorization header string true "Authentication header"
// @Param exchangeTradePairAssociationForm body forms.ExchangeTradePairAssociationForm true "Exchange Trade Pair Association"
// @Success 200 {object} responses.SuccessFailureResponse ""
// @Failure 200 {object} responses.SuccessFailureResponse ""
// @Router /v1/action/remove_exchange_trade_pair_association [post]
func RemoveExchangeTradePairsAssociationAction(c echo.Context) error {
	cc := c.(*util.EndpointContext)

	tradePairAssociationForm := new(forms.ExchangeTradePairAssociationForm)

	jsonBindError := cc.Context.Bind(tradePairAssociationForm)

	if jsonBindError != nil {
		sentry.CaptureException(jsonBindError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})

	}

	var exchangeModels []models.ExchangeModel

	findError := cc.AppContext.Database.Find(&exchangeModels).Error

	if findError != nil {
		sentry.CaptureException(findError)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrRetrievingExchangesError, Success: false})
	}

	for _, element := range exchangeModels {

		if element.UUID == tradePairAssociationForm.ExchangeUUID {

			var tradePairs []models.TradePairModel
			cc.AppContext.Database.Model(&element).Related(&tradePairs, "AuthorizedTradePairs")

			for tradePairIndex, tradePairElement := range tradePairs {
				if tradePairElement.UUID == tradePairAssociationForm.TradePairUUID {

					tradePairs = append(tradePairs[:tradePairIndex], tradePairs[tradePairIndex+1:]...)

					element.AuthorizedTradePairs = tradePairs

					var deleteReturn []string

					// Note: This is the current way to do this in gorm we need to be careful to avoid security issues.
					cc.AppContext.Database.Raw("DELETE FROM `exchange_authorized_trade_pairs` WHERE `exchange_model_id` = ? AND `trade_pair_model_id` = ?", element.ID, tradePairElement.ID).Scan(&deleteReturn)

					return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: "", Success: true})
				}
			}

			return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrTradePairAssociationDoesNotExistError, Success: false})

		}
	}

	// If we got here the exchange uuid is invalid

	sentry.CaptureMessage("Invalid Exchange UUID requested: "+tradePairAssociationForm.ExchangeUUID)
	return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrExchangeAccessTokenInvalid, Success: false})
}
