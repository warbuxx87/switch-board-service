package models

import "github.com/jinzhu/gorm"

type UserModel struct {
	gorm.Model

	UserUUID             string `json:"user_uuid" form:"user_uuid" query:"user_uuid"`
	Username             string `json:"username" form:"username" query:"username"`
	PasswordHash         string `json:"password_hash" form:"password_hash" query:"password_hash"`
	FirstName            string `json:"first_name" form:"first_name" query:"first_name"`
	LastName             string `json:"first_name" form:"first_name" query:"first_name"`
	EmailAddress         string `json:"email_address" form:"email_address" query:"email_address"`
}