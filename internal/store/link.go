package store

import (
	"context"

	"github.com/akithepriest/chisai.click/internal/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Link struct {
	ID           primitive.ObjectID `bson:"_id"`
	Keyword      string             `bson:"keyword"`
	RedirectLink string             `bson:"redirect_Link"`
	CreatedAt    primitive.DateTime `bson:"created_at"`
}

type LinkStore interface {
	GetByKeyword(ctx context.Context, keyword string) (*Link, error)
	Create(ctx context.Context, keyword string, redirectLink string) (*Link, error)
}

type LinkDBStore struct {
	collection *mongo.Collection
}

func NewLinkDBStore(db *mongo.Database) *LinkDBStore {
	return &LinkDBStore{
		collection: db.Collection("Links"),
	}
}

func (s *LinkDBStore) GetByKeyword(ctx context.Context, keyword string) (*Link, error) {
	var Link Link

	query := bson.M{"keyword": keyword}
	err := s.collection.FindOne(ctx, query).Decode(&Link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrNotFound
		}
		return nil, errs.ErrInternal(err)
	}

	return &Link, nil
}

func (s *LinkDBStore) Create(ctx context.Context, keyword string, redirectLink string) (*Link, error) {
	if _, err := s.GetByKeyword(ctx, keyword); err == nil {
		return nil, errs.ErrAlreadyExists
	}

	Link := &Link{
		ID:           primitive.NewObjectID(),
		Keyword:      keyword,
		RedirectLink: redirectLink,
	}
	_, err := s.collection.InsertOne(ctx, Link)
	if err != nil {
		return nil, errs.ErrInternal(err)
	}
	return Link, nil
}