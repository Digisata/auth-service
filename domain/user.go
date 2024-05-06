package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole int8

const (
	CollectionUser string = "users"

	Admin UserRole = iota + 1
	Customer
	Committee
)

type (
	User struct {
		ID       primitive.ObjectID `bson:"_id"`
		Name     string             `bson:"name"`
		Role     int8               `bson:"role"`
		Email    string             `bson:"email"`
		Password string             `bson:"password"`
	}

	UserProfile struct {
		Name  string
		Email string
	}

	Login struct {
		AccessToken  string
		RefreshToken string
	}

	RefreshTokenRequest struct {
		AccessToken  string
		RefreshToken string
	}
)
