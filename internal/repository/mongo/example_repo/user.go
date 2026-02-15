package example_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-starter/config"
	"go-starter/internal/lib/log"
	"go-starter/internal/lib/mongodb"
	do "go-starter/internal/models/do/mongo/example_do"
)

var ErrMongoUnavailable = errors.New("mongo user repository unavailable")

type UserRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*do.User, error)
	GetByCondAndPage(ctx context.Context, cond bson.M, skip, limit int64) ([]*do.User, int64, error)
	Create(ctx context.Context, user *do.User) error
	UpdateByID(ctx context.Context, id primitive.ObjectID, updates bson.M) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	IsAvailable() bool
}

type mongoUserRepo struct {
	client  *qmgo.QmgoClient
	initErr error
}

func NewUserRepository(c config.Config) UserRepository {
	cli, err := mongodb.New(c.MongoDB, do.User{}.Collection())
	if err != nil {
		if log.Logger != nil {
			log.Logger.Warnf("mongo user repository init failed, fallback to degraded mode: %v", err)
		}
		return &mongoUserRepo{
			initErr: err,
		}
	}
	return &mongoUserRepo{client: cli}
}

func (m *mongoUserRepo) IsAvailable() bool {
	return m.client != nil && m.initErr == nil
}

func (m *mongoUserRepo) ensureAvailable() error {
	if m.IsAvailable() {
		return nil
	}
	if m.initErr != nil {
		return fmt.Errorf("%w: %v", ErrMongoUnavailable, m.initErr)
	}
	return ErrMongoUnavailable
}

func (m *mongoUserRepo) GetByID(ctx context.Context, id primitive.ObjectID) (user *do.User, err error) {
	if err = m.ensureAvailable(); err != nil {
		return nil, err
	}

	user = &do.User{}
	err = m.client.Find(ctx, bson.M{
		"_id":      id,
		"isDelete": bson.M{"$ne": true},
	}).One(user)
	return
}

func (m *mongoUserRepo) GetByCondAndPage(ctx context.Context, cond bson.M, skip, limit int64) (users []*do.User, count int64, err error) {
	if err = m.ensureAvailable(); err != nil {
		return nil, 0, err
	}

	baseQuery := m.client.Find(ctx, cond)
	count, err = baseQuery.Count()
	if err != nil {
		return
	}

	err = baseQuery.Sort("-_id").Skip(skip).Limit(limit).All(&users)
	return
}

func (m *mongoUserRepo) Create(ctx context.Context, user *do.User) (err error) {
	if err = m.ensureAvailable(); err != nil {
		return err
	}

	_, err = m.client.InsertOne(ctx, user)
	return
}

func (m *mongoUserRepo) UpdateByID(ctx context.Context, id primitive.ObjectID, updates bson.M) (err error) {
	if err = m.ensureAvailable(); err != nil {
		return err
	}

	err = m.client.UpdateId(ctx, id, bson.M{"$set": updates})
	return
}

func (m *mongoUserRepo) DeleteByID(ctx context.Context, id primitive.ObjectID) (err error) {
	if err = m.ensureAvailable(); err != nil {
		return err
	}

	err = m.client.UpdateId(ctx, id, bson.M{
		"$set": bson.M{
			"isDelete": true,
		},
	})
	return
}
