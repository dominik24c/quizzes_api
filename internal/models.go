package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

type Token struct {
	AccessToken string `json:"access_token"`
}

type TokenPayload struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AuthResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Authorized bool   `json:"authorized"`
}

type AQuiz struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	UserId      primitive.ObjectID `bson:"user_id" json:"userId"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Categories  []string           `bson:"categories" json:"categories"`
	Questions   []Question         `bson:"questions" json:"questions"`
}

type Question struct {
	Question string   `bson:"question" json:"question"`
	Points   int      `bson:"points" json:"points"`
	Answers  []Answer `bson:"answers" json:"answers"`
}

type Answer struct {
	Answer    string `bson:"answer" json:"answer"`
	IsCorrect bool   `bson:"is_correct" json:"is_correct"`
}

type QuestionOut struct {
	Question string      `bson:"question" json:"question" validate:"required"`
	Answers  []AnswerOut `bson:"answers" json:"answers" validate:"required,min=1,max=4,dive"`
}

type AnswerOut struct {
	Answer string `bson:"answer" json:"answer" validate:"required"`
}

type ErrorMessage struct {
	Errors []ErrorField `json:"errors"`
}

type ErrorField struct {
	Field string `json:"field"`
	Value string `json:"value"`
}
