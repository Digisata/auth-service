package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Profile struct {
		ID        primitive.ObjectID `bson:"_id"`
		Name      string             `bson:"name"`
		Email     string             `bson:"email"`
		Password  string             `bson:"password"`
		CreatedAt int64              `bson:"created_at"`
		UpdatedAt int64              `bson:"updated_at"`
		DeletedAt int64              `bson:"deleted_at"`
	}

	UpdateProfile struct {
		Name      string `bson:"name,omitempty"`
		Password  string `bson:"password,omitempty"`
		UpdatedAt int64  `bson:"updated_at,omitempty"`
		DeletedAt int64  `bson:"deleted_at,omitempty"`
	}

	ChangePasswordRequest struct {
		OldPassword string
		NewPassword string
	}
)
