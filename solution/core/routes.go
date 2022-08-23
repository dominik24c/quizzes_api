package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dominik24c/quizzes_api/internal"
	"github.com/dominik24c/quizzes_api/internal/db"
	"github.com/dominik24c/quizzes_api/internal/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func addSolution(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// is authorized
	payload, err := utils.IsAuthorized(r.Header.Get("Authorization"))
	if err != nil || !payload.Authorized {
		utils.ErrorHandler(w, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	// get answers in quiz
	var quiz internal.AQuiz
	url := fmt.Sprintf("http://answers_service:9995/quiz/%s/answers", id)
	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Answers service error!")
		return
	}

	if response.StatusCode == http.StatusNotFound {
		utils.ErrorHandler(w, http.StatusNotFound, "Not Found Quiz!")
		return
	}

	var questionsBody QuestionBody
	err = json.NewDecoder(response.Body).Decode(&quiz)
	err1 := json.NewDecoder(r.Body).Decode(&questionsBody)
	if err != nil || err1 != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Internal error!")
		return
	}

	errMsg, err := utils.ValidateRequestBody(questionsBody)
	if err != nil {
		utils.JsonResponseHandler(w, http.StatusBadRequest, errMsg)
		return
	}

	//check answers and calculate totalPoints
	solvedQuiz, err := ValidateAndCalculateSolvedQuiz(questionsBody.Questions, quiz.Questions)
	if err != nil {
		utils.ErrorHandler(w, http.StatusUnprocessableEntity, "Invalid data you sent!")
		return
	}

	solvedQuiz.QuizId = quiz.ID
	solvedQuiz.UserId, _ = primitive.ObjectIDFromHex(payload.ID)
	solvedQuiz.Title = quiz.Title

	// insert
	collection := db.GetCollection(db.Client)
	result, err := collection.InsertOne(context.TODO(), solvedQuiz)
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Internal server error!")
		return
	}

	idStr, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		utils.ErrorHandler(w, http.StatusInternalServerError, "Cannot insert solution!")
		return
	}

	solutionResponse := SolutionResponse{
		ID:          idStr,
		Title:       quiz.Title,
		TotalPoints: solvedQuiz.TotalPoints,
		MaxPoints:   solvedQuiz.MaxPoints,
	}
	//msg["message"] = "Solution of quiz was added!"
	utils.JsonResponseHandler(w, http.StatusCreated, solutionResponse)
	return
}

func GetSolutionRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/quiz/{id}/solution", addSolution).Methods("POST")
	return r
}
