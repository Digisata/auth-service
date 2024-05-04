package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser string = "users"
)

type (
	User struct {
		ID       primitive.ObjectID `bson:"_id"`
		Name     string             `bson:"name"`
		Role     string             `bson:"role"`
		Email    string             `bson:"email"`
		Password string             `bson:"password"`
	}

	UserProfile struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)
