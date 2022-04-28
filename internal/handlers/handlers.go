package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Nisophix/crud_api/internal/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllPrivate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := ValidateKey(db, key, w)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}
	events := []models.Event{}
	db.Find(&events)
	respondJSON(w, http.StatusOK, events)
}

func CreatePrivate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	event := models.Event{}
	err := ValidateKey(db, key, w)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&event).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, event)
}

func ReadPrivate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])
	key := vars["key"]

	event, err := checkIDprivate(db, key, id, w, r)
	if err != nil {
		return
	}
	respondJSON(w, http.StatusOK, event)

}

// func CheckSubPrivate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	uuid := vars["uuid"]
// 	key := vars["key"]
// 	event, err := checkUUID(db, key, uuid, w, r)
// 	if err != nil {
// 		return
// 	}
// 	respondJSON(w, http.StatusOK, event.Sub)
// }

func CheckSub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	event := checkUUID(db, uuid, w, r)
	if event == nil {
		return
	}
	respondJSON(w, http.StatusOK, event.Sub)
}

func UpdatePrivate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])
	key := vars["key"]
	event, err := checkIDprivate(db, key, id, w, r)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&event).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, event)
}

func DeletePrivate(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	key := vars["key"]
	event, err := checkUUIDprivate(db, key, uuid, w, r)
	if err != nil {
		return
	}
	if err := db.Delete(&event).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func checkIDprivate(db *gorm.DB, key string, id int, w http.ResponseWriter, r *http.Request) (*models.Event, error) {
	err := ValidateKey(db, key, w)
	if err != nil {
		return nil, err

	}
	event := models.Event{}
	if err := db.First(&event, models.Event{ID: id}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil, err
	}
	return &event, nil
}
func checkUUIDprivate(db *gorm.DB, key string, uuid string, w http.ResponseWriter, r *http.Request) (*models.Event, error) {
	err := ValidateKey(db, key, w)
	if err != nil {
		return nil, err

	}
	event := models.Event{}
	if err := db.First(&event, models.Event{UUID: uuid}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil, err
	}
	return &event, nil
}

func checkUUID(db *gorm.DB, uuid string, w http.ResponseWriter, r *http.Request) *models.Event {
	event := models.Event{}
	if err := db.First(&event, models.Event{UUID: uuid}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &event
}

func checkID(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *models.Event {
	event := models.Event{}
	if err := db.First(&event, models.Event{ID: id}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &event
}

func ValidateKey(db *gorm.DB, key string, w http.ResponseWriter) error {
	if len(key) != 39 {
		err := errors.New("not a valid api key")
		respondError(w, http.StatusNotFound, err.Error())
		return err
	}

	apikey := models.APIKeys{}
	p := key[0:7]
	s := key[7:]
	hashedSecret := sha256.Sum256([]byte(s))

	if err := db.First(&apikey, models.APIKeys{Prefix: p}).Error; err != nil {
		err = errors.New("not a valid api key")
		respondError(w, http.StatusNotFound, err.Error())
		return err
	}
	if apikey.Secret != hex.EncodeToString(hashedSecret[:]) {
		err := errors.New("not a valid api key")
		respondError(w, http.StatusNotFound, err.Error())
		return err
	}
	return nil
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// func Delete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	uuid := vars["uuid"]
// 	event := checkUUID_or404(db, uuid, w, r)
// 	if event == nil {
// 		return
// 	}
// 	if err := db.Delete(&event).Error; err != nil {
// 		respondError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondJSON(w, http.StatusNoContent, nil)
// }

// func Update(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)

// 	id, _ := strconv.Atoi(vars["id"])
// 	event := checkID_or404(db, id, w, r)
// 	if event == nil {
// 		return
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&event); err != nil {
// 		respondError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	defer r.Body.Close()

// 	if err := db.Save(&event).Error; err != nil {
// 		respondError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondJSON(w, http.StatusOK, event)
// }

//возвращает просто true или false

// func Read(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, _ := strconv.Atoi(vars["id"])
// 	event := checkID_or404(db, id, w, r)
// 	if event == nil {
// 		return
// 	}
// 	respondJSON(w, http.StatusOK, event)
// }

// func Create(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	event := models.Event{}

// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&event); err != nil {
// 		respondError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	defer r.Body.Close()

// 	if err := db.Save(&event).Error; err != nil {
// 		respondError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondJSON(w, http.StatusCreated, event)
// }

// func GetAll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// events := []models.Event{}
// db.Find(&events)
// respondJSON(w, http.StatusOK, events)
// }
