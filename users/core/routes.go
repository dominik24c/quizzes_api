package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dominik24c/quizzes_api/internal"
	"github.com/dominik24c/quizzes_api/internal/db"
	"github.com/dominik24c/quizzes_api/internal/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user UserIn
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.ErrorHandler(w, http.StatusUnprocessableEntity, "Cannot process data!")
		return
	}

	errMsg, err := utils.ValidateRequestBody(user)
	if err != nil {
		utils.JsonResponseHandler(w, http.StatusBadRequest, errMsg)
		return
	}

	user.Password = HashPassword(user.Password)
	collection := db.GetCollection(db.Client)
	var u User
	err = collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&u)
	if err == nil {
		utils.ErrorHandler(w, http.StatusBadRequest, "This email exists in db!")
		return
	}

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Cannot insert data!")
		return
	}

	msg := make(map[string]string)
	msg["message"] = fmt.Sprintf("User was created %s!", result.InsertedID.(primitive.ObjectID).Hex())
	utils.JsonResponseHandler(w, http.StatusCreated, msg)
	return
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := db.GetCollection(db.Client)

	cur, err := collection.Find(context.TODO(), bson.M{}, utils.Pagination(r, 10))
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Cannot find data!")
		return
	}

	var users []User
	err = cur.All(context.TODO(), &users)
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Internal error!")
		return
	}

	userResponse := UserList{Users: users}
	utils.JsonResponseHandler(w, 0, userResponse)
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.ErrorHandler(w, http.StatusUnprocessableEntity, "Cannot process data!")
		return
	}

	// get hashed password from db
	collection := db.GetCollection(db.Client)
	var u UserPayload
	err = collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&u)
	if err != nil || u.Password == "" {
		utils.ErrorHandler(w, http.StatusNotFound, "Not found user!")
		return
	}

	// verify password
	if !VerifyPassword(u.Password, user.Password) {
		utils.ErrorHandler(w, http.StatusBadRequest, "Invalid credentials!")
		return
	}

	// get token
	authPayload := internal.TokenPayload{ID: u.ID, Email: user.Email}
	data, err := json.Marshal(authPayload)
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Internal error!")
		return
	}
	response, err := http.Post("http://auth_service:9998/auth/token", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		utils.ErrorHandler(w, http.StatusInternalServerError, "Auth service is not respond!")
		return
	}

	var token internal.Token
	err = json.NewDecoder(response.Body).Decode(&token)
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Auth service is not respond!")
		return
	}

	utils.JsonResponseHandler(w, 0, token)
	return
}

func GetUsersRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", listUsers).Methods("GET")
	r.HandleFunc("/register", createUser).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	return r
}
