package user

import (
	"context"
	"fmt"
	"time"
	"url-shortener/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mgo struct {
	db *mongo.Database
}

func NewMongoDB(db *mongo.Database) *mgo {
	return &mgo{db: db}
}

func (m *mgo) col() *mongo.Collection {
	return m.db.Collection("users")
}

func (m *mgo) CreateOrUpdate(ctx context.Context, user model.User) (*model.User, error) {
	const op = "user.mgo.CreateOrUpdate"

	if err := m.update(ctx, user, true); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (m *mgo) Update(ctx context.Context, user model.User) error {
	const op = "user.mgo.Update"

	if err := m.update(ctx, user, false); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *mgo) GetByGithubLogin(ctx context.Context, ghLogin string) (*model.User, error) {
	const op = "user.mgo.GetByGithubLogin"

	var u mgoUser
	if err := m.col().FindOne(ctx, bson.M{"_id": ghLogin}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s: %w", op, model.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return modelUserFromMgo(u), nil
}

func (m *mgo) Deactivate(ctx context.Context, ghLogin string) error {
	const op = "user.mgo.Deactivate"

	user, err := m.GetByGithubLogin(ctx, ghLogin)
	if err != nil {
		return fmt.Errorf("%s: %w", op, model.ErrNotFound)
	}

	user.IsActive = false

	if err := m.update(ctx, *user, false); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *mgo) update(ctx context.Context, user model.User, upsert bool) error {
	var (
		query       = bson.M{"_id": user.GithubLogin}
		replacement = mgoUserFromModel(user)
		opts        = &options.ReplaceOptions{Upsert: &upsert}
	)

	_, err := m.col().ReplaceOne(ctx, query, replacement, opts)
	if err != nil {
		return err
	}

	return nil
}

type mgoUser struct {
	IsActive    bool   `bson:"is_verified, omitempty"`
	GithubLogin string `bson:"_id"`

	// TODO: maybe store it in something like vault
	GithubAccessKey string    `bson:"gh_access_key,omitempty"`
	CreatedAt       time.Time `bson:"created_at,omitempty"`
}

func mgoUserFromModel(m model.User) mgoUser {
	return mgoUser{
		IsActive:        m.IsActive,
		GithubLogin:     m.GithubLogin,
		GithubAccessKey: m.GithubAccessKey,
		CreatedAt:       m.CreatedAt,
	}
}

func modelUserFromMgo(m mgoUser) *model.User {
	return &model.User{
		IsActive:        m.IsActive,
		GithubLogin:     m.GithubLogin,
		GithubAccessKey: m.GithubAccessKey,
		CreatedAt:       m.CreatedAt,
	}
}
