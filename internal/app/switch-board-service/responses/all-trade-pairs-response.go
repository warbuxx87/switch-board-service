package responses

import "bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"

type AllTradePairsResponse struct {
	Success      bool                    `json:"success" form:"success" query:"success"`
	ErrorMessage string                  `json:"error_message" form:"error_message" query:"error_message"`
	TradePairs   []models.TradePairModel `json:"trade_pairs" form:"trade_pairs" query:"trade_pairs"`
}
