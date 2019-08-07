package models

import "github.com/jinzhu/gorm"

type ExchangeModel struct {
	gorm.Model

	UUID       string `json:"uuid" form:"uuid" query:"uuid"`
	Name       string `json:"name" form:"name" query:"name"`
	Authorized bool   `json:"authorized" form:"authorized" query:"authorized"`

	AuthorizedTradePairs []AuthorizedTradePairModel `gorm:"foreignkey:ExchangeModelRefer"`
}
