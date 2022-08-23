package core

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	Email     string             `bson:"email"`
}

type UserIn struct {
	FirstName string `bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	Email     string `bson:"email" json:"email" validate:"required,email"`
	Password  string `bson:"password" json:"password" validate:"required,min=8,max=32"`
}

type UserList struct {
	Users []User `json:"users"`
}

type UserLogin struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type UserPayload struct {
	ID       string `bson:"_id" json:"id"`
	Password string `bson:"password" json:"password"`
}
