package core

import (
	"encoding/json"
	"github.com/dominik24c/quizzes_api/internal"
	"github.com/dominik24c/quizzes_api/internal/utils"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var SecretKey = os.Getenv("SECRET_KEY")

func getAccessToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tokenPayload internal.TokenPayload
	err := json.NewDecoder(r.Body).Decode(&tokenPayload)
	if err != nil {
		utils.ErrorHandler(w, http.StatusBadRequest, "Invalid data!")
		return
	}

	accessToken := GenerateJWT(tokenPayload.ID, tokenPayload.Email, SecretKey)
	if accessToken == "" {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Cannot generate access token!")
		return
	}

	token := internal.Token{AccessToken: accessToken}
	json.NewEncoder(w).Encode(&token)
	return
}

func authorize(w http.ResponseWriter, r *http.Request) {
	if authHeader := r.Header["Authorization"]; authHeader == nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Not found auth token!")
		return
	}

	claims, err := Authorize(r.Header.Get("Authorization"), SecretKey)
	if err != nil {
		utils.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	//fmt.Println(claims)
	utils.JsonResponseHandler(w, 200, claims)
	return

}

func GetAuthRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/auth/token", getAccessToken).Methods("POST")
	r.HandleFunc("/auth", authorize).Methods("POST")
	return r
}
