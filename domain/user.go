package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole int8

const (
	ADMIN UserRole = iota + 1
	CUSTOMER
	COMMITTEE

	USER_COLLECTION string = "users"
)

type (
	// User
	User struct {
		ID        primitive.ObjectID `bson:"_id"`
		Name      string             `bson:"name"`
		Role      int8               `bson:"role"`
		Email     string             `bson:"email"`
		Password  string             `bson:"password"`
		IsActive  bool               `bson:"is_active"`
		Note      string             `bson:"note"`
		CreatedAt int64              `bson:"created_at"`
		UpdatedAt int64              `bson:"updated_at"`
		DeletedAt int64              `bson:"deleted_at"`
	}

	AuthResponse struct {
		AccessToken  string
		RefreshToken string
	}

	RefreshTokenRequest struct {
		AccessToken  string
		RefreshToken string
	}

	GetAllUserRequest struct {
		Search   string
		IsActive bool
	}

	UpdateUser struct {
		ID        primitive.ObjectID `bson:"_id"`
		Name      string             `bson:"name,omitempty"`
		IsActive  bool               `bson:"is_active"`
		Note      string             `bson:"note,omitempty"`
		UpdatedAt int64              `bson:"updated_at,omitempty"`
		DeletedAt int64              `bson:"deleted_at,omitempty"`
	}

	DeleteUser struct {
		ID        primitive.ObjectID `bson:"_id"`
		IsActive  bool               `bson:"is_active"`
		UpdatedAt int64              `bson:"updated_at,omitempty"`
		DeletedAt int64              `bson:"deleted_at,omitempty"`
	}
)
