package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nisophix/crud_api/internal/database/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAll(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	events := []models.Event{}
	db.Find(&events)
	respondJSON(w, http.StatusOK, events)
}

func Create(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	event := models.Event{}

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

func Read(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])
	event := checkID_or404(db, id, w, r)
	if event == nil {
		return
	}
	respondJSON(w, http.StatusOK, event)
}

//возвращает просто true или false
func CheckSub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	event := checkUUID_or404(db, uuid, w, r)
	if event == nil {
		return
	}
	respondJSON(w, http.StatusOK, event.Sub)
}

func Update(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])
	event := checkID_or404(db, id, w, r)
	if event == nil {
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

func Delete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	event := checkUUID_or404(db, uuid, w, r)
	if event == nil {
		return
	}
	if err := db.Delete(&event).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func checkID_or404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *models.Event {
	event := models.Event{}
	if err := db.First(&event, models.Event{ID: id}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &event
}

func checkUUID_or404(db *gorm.DB, uuid string, w http.ResponseWriter, r *http.Request) *models.Event {
	event := models.Event{}
	if err := db.First(&event, models.Event{UUID: uuid}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &event
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
