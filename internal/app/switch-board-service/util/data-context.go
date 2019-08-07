package util

import "github.com/jinzhu/gorm"

type DataContext struct {
	Database *gorm.DB
	Config   *Configuration
}
