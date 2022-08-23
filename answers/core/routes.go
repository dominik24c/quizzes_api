package core

import (
	"context"
	"github.com/dominik24c/quizzes_api/internal"
	"github.com/dominik24c/quizzes_api/internal/db"
	"github.com/dominik24c/quizzes_api/internal/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func getAnswers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	quidId := vars["id"]

	_id, err := primitive.ObjectIDFromHex(quidId)

	if err != nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Not found quiz")
		return
	}

	collection := db.GetCollection(db.Client)
	var quiz internal.AQuiz

	err = collection.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&quiz)
	if err != nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Not found quiz")
		return
	}

	utils.JsonResponseHandler(w, http.StatusOK, quiz)
}

func GetAnswersRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/quiz/{id}/answers", getAnswers).Methods("POST")
	return r
}
