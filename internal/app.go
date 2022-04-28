package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/Nisophix/crud_api/internal/config"
	"github.com/Nisophix/crud_api/internal/database"
	"github.com/Nisophix/crud_api/internal/handlers"
	"github.com/Nisophix/crud_api/internal/models"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(config *config.Config) {
	db, _ := database.Connection(config)
	a.DB = models.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) Start(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

func (a *App) setRouters() {
	//get everything
	// a.Get("/events", a.handleRequest(handlers.GetAll))
	a.Get("/{key}/users", a.handleRequest(handlers.GetAllPrivate))

	//post a new event
	// a.Post("/event", a.handleRequest(handlers.Create))
	a.Post("/{key}/createuser", a.handleRequest(handlers.CreatePrivate))

	//get an event by id
	// a.Get("/event/{id}", a.handleRequest(handlers.Read))
	a.Get("/{key}/getuser/{id}", a.handleRequest(handlers.ReadPrivate))

	//update an event by id
	// a.Put("/event/{id}", a.handleRequest(handlers.Update))
	a.Put("/{key}/updateuser/{id}", a.handleRequest(handlers.UpdatePrivate))

	//delet an event by uuid
	// a.Delete("/event/{uuid}", a.handleRequest(handlers.Delete))
	a.Delete("/{key}/deleteuser/{uuid}", a.handleRequest(handlers.DeletePrivate))

	//check status by uuid
	a.Get("/userstatus/{uuid}", a.handleRequest(handlers.CheckSub))
	// a.Get("/{key}/userstatus/{uuid}", a.handleRequest(handlers.CheckSubPrivate))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}
