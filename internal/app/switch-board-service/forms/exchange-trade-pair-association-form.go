package forms

type ExchangeTradePairAssociationForm struct {
	ExchangeUUID  string `json:"exchange_uuid" form:"exchange_uuid" query:"exchange_uuid"`
	TradePairUUID string `json:"trade_pair_uuid" form:"trade_pair_uuid" query:"trade_pair_uuid"`
}
