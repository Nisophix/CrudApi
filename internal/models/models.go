package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Event is a model for a database table
type Event struct {
	ID int
}

//DBMigrate creates a table acording to the model
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Event{})
	return db
}
