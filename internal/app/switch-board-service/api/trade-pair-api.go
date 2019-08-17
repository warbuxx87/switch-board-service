package api

import (
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/forms"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"
	"bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/responses"
	"github.com/getsentry/raven-go"
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
		raven.CaptureError(jsonBindError, nil)
		return c.JSON(http.StatusOK, &responses.SuccessFailureResponse{ErrorMessage: util.StrIncorrectParametersPassed, Success: false})

	}

	var tradePairModel models.TradePairModel

	tradePairModel.UUID = uuid.NewV4().String()
	tradePairModel.Name = tradePairForm.Name
	tradePairModel.BaseAsset = tradePairForm.BaseAsset
	tradePairModel.SecondaryAsset = tradePairForm.SecondaryAsset

	createError := cc.AppContext.Database.Create(&tradePairForm).Error

	if createError != nil {
		raven.CaptureError(createError, nil)
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
func RetrieveAllTradePairs(c echo.Context) error {

	cc := c.(*util.EndpointContext)

	var tradePairModels []models.TradePairModel

	findError := cc.AppContext.Database.Find(&tradePairModels).Error

	if findError != nil {
		raven.CaptureError(findError, nil)
		return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{ErrorMessage: util.StrRetrievingTradePairsError, Success: false})
	}

	return c.JSON(http.StatusOK, &responses.AllTradePairsResponse{TradePairs: tradePairModels, ErrorMessage: "", Success: true})
}
