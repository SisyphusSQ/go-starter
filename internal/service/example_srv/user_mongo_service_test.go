package example_srv

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap/zapcore"

	"go-starter/config"
	"go-starter/internal/lib/log"
	do "go-starter/internal/models/do/mongo/example_do"
	"go-starter/internal/models/vo"
	"go-starter/internal/repository/mongo/example_repo"
	"go-starter/utils"
)

func TestMain(m *testing.M) {
	cfg := config.Config{Log: config.Log{LogLevel: zapcore.InfoLevel}}
	log.New(cfg)
	code := m.Run()
	if log.Logger != nil {
		log.Logger.Sync()
	}
	os.Exit(code)
}

type mockMongoUserRepo struct {
	available bool
	users     []*do.User
	total     int64
	listErr   error

	lastCond  bson.M
	lastSkip  int64
	lastLimit int64
}

func (m *mockMongoUserRepo) IsAvailable() bool {
	return m.available
}

func (m *mockMongoUserRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*do.User, error) {
	if !m.available {
		return nil, example_repo.ErrMongoUnavailable
	}
	return &do.User{}, nil
}

func (m *mockMongoUserRepo) GetByCondAndPage(ctx context.Context, cond bson.M, skip, limit int64) ([]*do.User, int64, error) {
	if !m.available {
		return nil, 0, example_repo.ErrMongoUnavailable
	}
	m.lastCond = cond
	m.lastSkip = skip
	m.lastLimit = limit
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.users, m.total, nil
}

func (m *mockMongoUserRepo) Create(ctx context.Context, user *do.User) error {
	if !m.available {
		return example_repo.ErrMongoUnavailable
	}
	return nil
}

func (m *mockMongoUserRepo) UpdateByID(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	if !m.available {
		return example_repo.ErrMongoUnavailable
	}
	return nil
}

func (m *mockMongoUserRepo) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	if !m.available {
		return example_repo.ErrMongoUnavailable
	}
	return nil
}

func TestUserMongoServiceUnavailable(t *testing.T) {
	repo := &mockMongoUserRepo{available: false}
	srv := NewUserMongoService(repo, time.Second)

	require.False(t, srv.IsAvailable())
	_, err := srv.List(context.Background(), 1, 10)
	require.ErrorIs(t, err, example_repo.ErrMongoUnavailable)
}

func TestUserMongoServiceListValidateBaseList(t *testing.T) {
	repo := &mockMongoUserRepo{available: true}
	srv := NewUserMongoService(repo, time.Second)

	_, err := srv.List(context.Background(), 0, 10)
	require.ErrorIs(t, err, utils.ErrBadParamInput)

	_, err = srv.List(context.Background(), 1, vo.MaxPageSize+1)
	require.ErrorIs(t, err, utils.ErrBadParamInput)
}

func TestUserMongoServiceListPassesPagination(t *testing.T) {
	repo := &mockMongoUserRepo{
		available: true,
		users:     []*do.User{{Name: "alice"}},
		total:     2,
	}
	srv := NewUserMongoService(repo, time.Second)

	resp, err := srv.List(context.Background(), 3, 10)
	require.NoError(t, err)
	require.Len(t, resp.List, 1)
	require.Equal(t, int64(2), resp.Total)

	require.Equal(t, int64(20), repo.lastSkip)
	require.Equal(t, int64(10), repo.lastLimit)
	require.Equal(t, bson.M{"isDelete": bson.M{"$ne": true}}, repo.lastCond)
}
