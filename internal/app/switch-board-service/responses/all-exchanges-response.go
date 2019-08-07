package responses

import "bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"

type AllExchangesResponse struct {
	Success      bool                   `json:"success" form:"success" query:"success"`
	ErrorMessage string                 `json:"error_message" form:"error_message" query:"error_message"`
	Exchanges    []models.ExchangeModel `json:"exchanges" form:"exchanges" query:"exchanges"`
}
