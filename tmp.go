package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Payload struct {
	Name       string    `json:"Name"`
	UUID       string    `json:"UUID"`
	Subbed     bool      `json:"Subbed"`
	TimeSubbed time.Time `json:"TimeSubbed"`
	Expires    time.Time `json:"Expires"`
}

func CreateUser(Name, UUID string, Subbed bool, t time.Time, weeks int) {
	data := Payload{
		Name:       Name,
		UUID:       UUID,
		Subbed:     Subbed,
		TimeSubbed: t,
		Expires:    t.Add(time.Hour * 24 * 7 * time.Duration(weeks)),
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/event", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

}

func main() {
	CreateUser("TestUser", "1234-5678-90AB-CDEF", true, time.Now(), 8)
}
