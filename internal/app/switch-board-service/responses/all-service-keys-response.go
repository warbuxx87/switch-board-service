package responses

import "bitbucket.org/project-alkemi/switch-board-service/internal/app/switch-board-service/models"

type AllServiceKeysResponse struct {
	Success      bool                     `json:"success" form:"success" query:"success"`
	ErrorMessage string                   `json:"error_message" form:"error_message" query:"error_message"`
	ServiceKeys  []models.ServiceKeyModel `json:"service_keys" form:"service_keys" query:"service_keys"`
}
