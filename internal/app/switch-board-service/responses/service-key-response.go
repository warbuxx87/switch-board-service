package responses

type ServiceKeyResponse struct {
	Success      bool   `json:"success" form:"success" query:"success"`
	ErrorMessage string `json:"error_message" form:"error_message" query:"error_message"`
	ServiceUUID   string `json:"service_uuid" form:"service_uuid" query:"service_uuid"`
	ServiceKey   string `json:"service_key" form:"service_key" query:"service_key"`
}
