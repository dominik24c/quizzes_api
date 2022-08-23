package main

import (
	"github.com/dominik24c/quizzes_api/internal/db"
	"github.com/dominik24c/quizzes_api/quiz/core"
	"log"
	"net/http"
	"os"
)

func main() {
	db.Client = db.InitDatabase(os.Getenv("MONGO_URI"))
	log.Fatal(http.ListenAndServe(":9997", core.GetQuizRouter()))
}
