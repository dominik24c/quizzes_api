package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dominik24c/quizzes_api/internal/db"
	"github.com/dominik24c/quizzes_api/internal/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func createQuiz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload, err := utils.IsAuthorized(r.Header.Get("Authorization"))
	if err != nil || !payload.Authorized {
		utils.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	var quiz Quiz
	err = json.NewDecoder(r.Body).Decode(&quiz)
	if err != nil {
		utils.ErrorHandler(w, http.StatusUnprocessableEntity, "Cannot process data!")
		return
	}

	errMsg, err := utils.ValidateRequestBody(quiz)
	if err != nil {
		utils.JsonResponseHandler(w, http.StatusBadRequest, errMsg)
		return
	}
	quiz.UserId, _ = primitive.ObjectIDFromHex(payload.ID)

	collection := db.GetCollection(db.Client)

	result, err := collection.InsertOne(context.TODO(), quiz)

	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Cannot insert data to db!")
		return
	}

	msg := make(map[string]string)
	msg["message"] = fmt.Sprintf("Quiz %v was created!", result.InsertedID.(primitive.ObjectID).Hex())
	utils.JsonResponseHandler(w, http.StatusCreated, msg)
	return
}

func listQuiz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := db.GetCollection(db.Client)

	res, err := collection.Find(context.TODO(), bson.M{}, utils.Pagination(r, 10))
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Cannot find data!")
		return
	}

	var quizzes []QuizInfo
	err = res.All(context.TODO(), &quizzes)
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Internal error!")
		return
	}

	utils.JsonResponseHandler(w, 0, QuizList{Quizzes: quizzes})
	return
}

func getQuiz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload, err := utils.IsAuthorized(r.Header.Get("Authorization"))
	if err != nil || !payload.Authorized {
		utils.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized!")
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Not found quiz")
		return
	}

	collection := db.GetCollection(db.Client)
	var quiz QuizOut
	err = collection.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&quiz)
	if err != nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Not found quiz")
		return
	}

	utils.JsonResponseHandler(w, 0, quiz)
	return
}

func deleteQuiz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload, err := utils.IsAuthorized(r.Header.Get("Authorization"))
	if err != nil || !payload.Authorized {
		utils.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Not found quiz")
		return
	}

	collection := db.GetCollection(db.Client)
	var quiz QuizOut
	userId, _ := primitive.ObjectIDFromHex(payload.ID)

	err = collection.FindOne(context.TODO(), bson.M{"_id": _id, "user_id": userId}).Decode(&quiz)
	if err != nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Not found quiz")
		return
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": _id})
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Cannot delete quiz")
		return
	}

	utils.JsonResponseHandler(w, http.StatusNoContent, "")
	return
}

func GetQuizRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/quiz", listQuiz).Methods("GET")
	r.HandleFunc("/quiz", createQuiz).Methods("POST")
	r.HandleFunc("/quiz/{id}", getQuiz).Methods("GET")
	r.HandleFunc("/quiz/{id}", deleteQuiz).Methods("DELETE")
	return r
}
