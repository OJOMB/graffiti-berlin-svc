package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type loginRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (app *App) handleAuthenticate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var reqBodyData loginRequest
		if err := json.Unmarshal(reqBody, &reqBodyData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		user, dErr := app.service.ValidateUserCredentials(r.Context(), reqBodyData.UserName, reqBodyData.Email, reqBodyData.Password)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		// user credentials are valid so now we need to return a token
		tokenString, err := app.tokenAuth.GenerateTokenString(user.ID)
		if err != nil {
			apperr := newAppErr("failed to generate token", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		resp := loginResponse{
			AccessToken: tokenString,
			TokenType:   "Bearer",
		}
		respBodyBytes, err := json.Marshal(resp)
		if err != nil {
			apperr := newAppErr("failed to marsal response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respBodyBytes)
	}
}
