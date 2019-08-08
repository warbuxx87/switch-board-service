package responses

import "bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"

type ExchangeInfoResponse struct {
	Success      bool                 `json:"success" form:"success" query:"success"`
	ErrorMessage string               `json:"error_message" form:"error_message" query:"error_message"`
	Exchange     models.ExchangeModel `json:"exchange" form:"exchange" query:"exchange"`
}
