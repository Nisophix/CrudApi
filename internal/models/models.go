package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Sub struct {
	Subbed bool `gorm:"type:bool;default:false"`
}

//Event is a model for a database table
type Event struct {
	ID         int       `gorm:"primaryKey";auto_increment;not_null"`
	UUID       string    `gorm:"type:VARCHAR(36);unique;not_null"`
	Name       string    `gorm:"type:VARCHAR(20);not_null"`
	Sub                  //bool `gorm:"type:bool;default:false"`
	TimeSubbed time.Time `gorm:"index"`
	Expires    time.Time
}

type APIKeys struct {
	ID     int    `gorm:"primaryKey";auto_increment;not_null"`
	Prefix string `gorm:"type:VARCHAR(7);unique;not_null"`
	Secret string `gorm:"type:VARCHAR(64);unique;not_null"`
}

//DBMigrate creates a table acording to the model
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Event{}, &APIKeys{})
	return db
}
