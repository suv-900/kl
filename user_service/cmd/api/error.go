package api

import "net/http"

func (app *application) sendErrorResponse(w http.ResponseWriter, statusCode int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, env, statusCode)
	if err != nil {
		//err while sending
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) invalidAuthenticationHeader(w http.ResponseWriter) {
	message := "invalid authentication headers."
	app.sendErrorResponse(w, http.StatusExpectationFailed, message)
}

func (app *application) invalidToken(w http.ResponseWriter) {
	message := "token invalid"
	app.sendErrorResponse(w, http.StatusUnauthorized, message)
}

func (app *application) invalidTokenDeletedUser(w http.ResponseWriter) {
	message := "user is deleted[dead token] please create new account."
	app.sendErrorResponse(w, http.StatusNotFound, message)
}

func (app *application) internalServerError(w http.ResponseWriter) {
	message := "server error"
	app.sendErrorResponse(w, http.StatusInternalServerError, message)
}
