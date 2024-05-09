package repository

import (
	"context"
	"time"

	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db         mongo.Database
	collection string
}

func NewUserRepository(db mongo.Database, collection string) *UserRepository {
	return &UserRepository{
		db:         db,
		collection: collection,
	}
}

func (r UserRepository) Create(ctx context.Context, user domain.User) error {
	collection := r.db.Collection(r.collection)

	now := time.Now().Local().Unix()
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	collection := r.db.Collection(r.collection)
	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var users []domain.User

	err = cursor.All(ctx, &users)
	if users == nil {
		return []domain.User{}, err
	}

	return users, nil
}

func (r UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	collection := r.db.Collection(r.collection)

	var user domain.User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	collection := r.db.Collection(r.collection)

	var user domain.User

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}
