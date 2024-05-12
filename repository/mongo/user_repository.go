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

func (r UserRepository) Create(ctx context.Context, req domain.User) error {
	collection := r.db.Collection(r.collection)
	user := req

	now := time.Now().Local().Unix()
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetAll(ctx context.Context, req domain.GetAllUserRequest) ([]domain.User, error) {
	collection := r.db.Collection(r.collection)
	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})

	filter := bson.M{}
	if req.Search != "" {
		pattern := ".*" + req.Search + ".*"
		filter["$or"] = []interface{}{
			bson.M{"name": bson.M{"$regex": pattern, "$options": "i"}},
			bson.M{"email": bson.M{"$regex": pattern, "$options": "i"}},
		}
	}

	filter["is_active"] = req.IsActive
	filter["role"] = bson.M{"$ne": domain.ADMIN}
	filter["deleted_at"] = 0

	cursor, err := collection.Find(ctx, filter, opts)
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

func (r UserRepository) Update(ctx context.Context, req domain.UpdateUser) error {
	collection := r.db.Collection(r.collection)

	updateUser := req
	updateUser.UpdatedAt = time.Now().Local().Unix()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": req.ID}, bson.M{"$set": updateUser})
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) Delete(ctx context.Context, req domain.DeleteUser) error {
	collection := r.db.Collection(r.collection)

	deleteUser := req
	now := time.Now().Local().Unix()
	deleteUser.IsActive = false
	deleteUser.UpdatedAt = now
	deleteUser.DeletedAt = now

	_, err := collection.UpdateOne(ctx, bson.M{"_id": req.ID}, bson.M{"$set": deleteUser})
	if err != nil {
		return err
	}

	return nil
}
