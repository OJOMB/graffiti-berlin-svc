package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const handlerCreateUser = "handleCreateUser"

type createUserReq struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *App) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var userData createUserReq
		if err := json.Unmarshal(reqBodyBytes, &userData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		user, dErr := app.service.CreateUser(r.Context(), userData.UserName, userData.Email, userData.Password)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBodyBytes, err := json.Marshal(user)
		if err != nil {
			app.logger.WithField(appHandler, handlerCreateUser).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(respBodyBytes)
	}
}
