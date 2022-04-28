package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"

	"log"

	"github.com/jinzhu/gorm"

	app "github.com/Nisophix/crud_api/internal"
	"github.com/Nisophix/crud_api/internal/config"
)

type APIKey struct {
	Prefix string
	Secret string
}

func genkey() APIKey {
	p := Randstring(7)
	s := Randstring(32)
	rawkey := APIKey{
		Prefix: p,
		Secret: s,
	}
	return rawkey
}

func Randstring(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}
func keyToDB(db *gorm.DB) (APIKey, error) {
	rawkey := genkey()
	secret := sha256.Sum256([]byte(rawkey.Secret))
	apikey := APIKey{
		Prefix: rawkey.Prefix,
		Secret: hex.EncodeToString(secret[:]),
	}
	err := db.Create(&apikey).Error
	if err != nil {
		fmt.Println(err)
	}
	return rawkey, nil
}

func CreateAPIKey() (APIKey, error) {
	config, err := config.ReadConfig("./config.toml")
	if err != nil {
		log.Fatal(err)
	}
	app := &app.App{}
	app.Initialize(config)
	key, _ := keyToDB(app.DB)
	return key, nil
}
