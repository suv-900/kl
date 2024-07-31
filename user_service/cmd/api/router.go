package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) SetupRouter() {

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(app.routeNotFound)

	r.HandleFunc("/v1/register", app.RegisterUser)
}
