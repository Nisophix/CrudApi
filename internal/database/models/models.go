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
	ID   int    `gorm:"primaryKey";auto_increment;not_null"`
	UUID string `gorm:"type:VARCHAR(36);unique;not_null"`
	Name string `gorm:"type:VARCHAR(20);not_null"`
	Sub
	TimeSubbed time.Time `gorm:"index"`
	Expires    time.Time
}

//DBMigrate creates a table acording to the model
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Event{})
	return db
}
