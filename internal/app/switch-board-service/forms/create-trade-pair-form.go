package forms

type CreateTradePairForm struct {
	Name           string `json:"name" form:"name" query:"name"`
	BaseAsset      string `json:"base_asset" form:"base_asset" query:"base_asset"`
	SecondaryAsset string `json:"secondary_asset" form:"secondary_asset" query:"secondary_asset"`
}
