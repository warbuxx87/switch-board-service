package models

import "github.com/jinzhu/gorm"

type ExchangeModel struct {
	gorm.Model

	UUID            string `json:"uuid" form:"uuid" query:"uuid"`
	Name            string `json:"name" form:"name" query:"name"`
	Authorized      bool   `json:"authorized" form:"authorized" query:"authorized"`
	AccessToken     string `json:"access_token" form:"access_token" query:"access_token"`
	TempAccessToken string `json:"temp_access_token" form:"temp_access_token" query:"temp_access_token"`

	AuthorizedTradePairs []*TradePairModel `gorm:"many2many:exchange_authorized_trade_pairs;" json:"authorized_trade_pairs" form:"authorized_trade_pairs" query:"authorized_trade_pairs"`
}
