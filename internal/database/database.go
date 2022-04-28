package database

import (
	"fmt"
	"log"

	"github.com/Nisophix/crud_api/internal/config"
	"github.com/jinzhu/gorm"
)

func Connection(config *config.Config) (*gorm.DB, error) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	}
	return db, err
}
