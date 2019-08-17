package forms

type CreateExchangeForm struct {
	Name string `json:"name" form:"name" query:"name"`
}