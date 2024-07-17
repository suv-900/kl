package api

import "net/http"

func (app *application) sendResponse(w http.ResponseWriter, statusCode int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, env, statusCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) invalidAuthenticationHeader(w http.ResponseWriter, r *http.Request) {
}
