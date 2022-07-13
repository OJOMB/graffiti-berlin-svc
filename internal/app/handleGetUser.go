package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

const handleGetUser = "handleGetUser"

func (app *App) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars[urlVarUserID]

		user, dErr := app.service.GetUser(r.Context(), userID)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		} else if user == nil {
			apperr := newAppErr("user not found", http.StatusNotFound)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBytes, err := json.Marshal(user)
		if err != nil {
			app.logger.WithField(appHandler, handleGetUser).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.Write(respBytes)
	}
}
