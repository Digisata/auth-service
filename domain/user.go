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
