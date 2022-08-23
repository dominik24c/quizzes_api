package core

import (
	"github.com/dominik24c/quizzes_api/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	UserId      primitive.ObjectID `bson:"user_id" json:"userId"`
	Title       string             `bson:"title" json:"title" validate:"required,min=3,max=100"`
	Description string             `bson:"description" json:"description" validate:"required,max=1000"`
	Categories  []string           `bson:"categories" json:"categories" validate:"required,dive,min=2,max=60"`
	Questions   []Question         `bson:"questions" json:"questions" validate:"required,min=4,max=20,dive"`
}

type Question struct {
	Question string   `bson:"question" json:"question" validate:"required,min=2,max=100"`
	Points   int      `bson:"points" json:"points" validate:"required,gte=1,lte=6"`
	Answers  []Answer `bson:"answers" json:"answers" validate:"required,min=2,max=4"`
}

type Answer struct {
	Answer    string `bson:"answer" json:"answer" validate:"required,min=1,max=50"`
	IsCorrect bool   `bson:"is_correct" json:"is_correct" validate:"required"`
}

type QuizList struct {
	Quizzes []QuizInfo `json:"quizzes"`
}

type QuizInfo struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Categories  []string           `bson:"categories" json:"categories"`
}
type QuizOut struct {
	ID          primitive.ObjectID     `bson:"_id" json:"id"`
	Title       string                 `bson:"title" json:"title"`
	Description string                 `bson:"description" json:"description"`
	Questions   []internal.QuestionOut `bson:"questions" json:"questions"`
}
