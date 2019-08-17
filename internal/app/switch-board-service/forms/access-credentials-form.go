package forms

type AccessCredentialsForm struct {
	UUID        string `json:"uuid" form:"uuid" query:"uuid"`
	AccessToken string `json:"access_token" form:"access_token" query:"access_token"`
}
