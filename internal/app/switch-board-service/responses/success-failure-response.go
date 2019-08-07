package responses

type SuccessFailureResponse struct {
	Success      bool   `json:"success" form:"success" query:"success"`
	ErrorMessage string `json:"error_message" form:"error_message" query:"error_message"`
}