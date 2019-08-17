package forms

type CreateServiceForm struct {
	Name string `json:"name" form:"name" query:"name"`
}