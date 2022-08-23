package main

import (
	"github.com/dominik24c/quizzes_api/answers/core"
	"github.com/dominik24c/quizzes_api/internal/db"
	"log"
	"net/http"
	"os"
)

func main() {
	db.Client = db.InitDatabase(os.Getenv("MONGO_URI"))
	log.Fatal(http.ListenAndServe(":9995", core.GetAnswersRouter()))
}
