package models

import "github.com/jinzhu/gorm"

type ServiceKeyModel struct {
	gorm.Model

	UUID            string `gorm:"unique_index" json:"uuid" form:"uuid" query:"uuid"`
	Name            string `gorm:"unique_index" json:"name" form:"name" query:"name"`
	AccessToken     string `json:"access_token" form:"access_token" query:"access_token"`
}
