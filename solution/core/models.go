package core

import (
	"github.com/dominik24c/quizzes_api/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizSolution struct {
	ID          primitive.ObjectID     `json:"id" bson:"id"`
	UserId      primitive.ObjectID     `json:"user_id" bson:"user_id"`
	QuizId      primitive.ObjectID     `json:"quiz_id" bson:"quiz_id"`
	Title       string                 `json:"title" bson:"title"`
	Questions   []internal.QuestionOut `json:"questions" bson:"questions"`
	TotalPoints float32                `json:"total_points" bson:"total_points"`
	MaxPoints   float32                `json:"max_points" bson:"max_points"`
}

type QuestionBody struct {
	Questions []internal.QuestionOut `json:"questions" bson:"questions" validate:"required,min=4,max=20,dive"`
}

type SolutionResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"quiz_title"`
	TotalPoints float32            `json:"total_points"`
	MaxPoints   float32            `json:"max_points"`
}
