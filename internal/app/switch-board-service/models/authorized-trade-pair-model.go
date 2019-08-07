package models

import "github.com/jinzhu/gorm"

type AuthorizedTradePairModel struct {
	gorm.Model

	Exchange  ExchangeModel  `gorm:"foreignkey:AuthorizedTradePairModelRefer"`
	TradePair TradePairModel `gorm:"foreignkey:AuthorizedTradePairModelRefer"`
}
