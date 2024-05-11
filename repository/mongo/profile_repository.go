package repository

import (
	"context"
	"time"

	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileRepository struct {
	db         mongo.Database
	collection string
}

func NewProfileRepository(db mongo.Database, collection string) *ProfileRepository {
	return &ProfileRepository{
		db:         db,
		collection: collection,
	}
}

func (r ProfileRepository) GetByID(ctx context.Context, id string) (domain.Profile, error) {
	collection := r.db.Collection(r.collection)

	var profile domain.Profile

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return profile, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&profile)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func (r ProfileRepository) ChangePassword(ctx context.Context, id, newPassword string) error {
	collection := r.db.Collection(r.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateProfile := domain.UpdateProfile{
		Password:  newPassword,
		UpdatedAt: time.Now().Local().Unix(),
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": idHex}, bson.M{"$set": updateProfile})
	if err != nil {
		return err
	}

	return nil
}
