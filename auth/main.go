package main

import (
	"github.com/dominik24c/quizzes_api/auth/core"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":9998", core.GetAuthRouter()))
}
