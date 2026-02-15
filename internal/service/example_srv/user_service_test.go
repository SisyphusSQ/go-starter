package example_srv

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go-starter/config"
	gormv2 "go-starter/internal/lib/gorm"
	redisv9 "go-starter/internal/lib/redis"
	do "go-starter/internal/models/do/mysql/example_do"
	"go-starter/internal/models/vo"
	"go-starter/utils"
)

type mockMySQLUserRepo struct {
	lastConditions *gormv2.DBConditions
	listResult     []do.User
	listErr        error
	total          int64
}

func (m *mockMySQLUserRepo) GetByID(ctx context.Context, id int64) (do.User, error) {
	return do.User{}, nil
}

func (m *mockMySQLUserRepo) GetByEmail(ctx context.Context, email string) (do.User, error) {
	return do.User{}, nil
}

func (m *mockMySQLUserRepo) ListByConditions(ctx context.Context, conditions *gormv2.DBConditions) ([]do.User, error) {
	m.lastConditions = conditions
	if conditions != nil {
		conditions.Count = m.total
	}
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.listResult, nil
}

func (m *mockMySQLUserRepo) Create(ctx context.Context, user *do.User) error {
	return nil
}

func (m *mockMySQLUserRepo) UpdateByID(ctx context.Context, id int64, updates map[string]any) error {
	return nil
}

func (m *mockMySQLUserRepo) DeleteByID(ctx context.Context, id int64) error {
	return nil
}

func TestUserServiceListValidateBaseList(t *testing.T) {
	repo := &mockMySQLUserRepo{}
	srv := NewUserService(repo, time.Second, config.Config{}, &redisv9.Client{})

	_, err := srv.List(context.Background(), 0, 10)
	require.ErrorIs(t, err, utils.ErrBadParamInput)

	_, err = srv.List(context.Background(), 1, vo.MaxPageSize+1)
	require.ErrorIs(t, err, utils.ErrBadParamInput)
}

func TestUserServiceListUsesDBConditions(t *testing.T) {
	repo := &mockMySQLUserRepo{
		listResult: []do.User{{ID: 1, Name: "alice"}},
		total:      3,
	}
	srv := NewUserService(repo, time.Second, config.Config{}, &redisv9.Client{})

	resp, err := srv.List(context.Background(), 2, 10)
	require.NoError(t, err)
	require.Len(t, resp.List, 1)
	require.Equal(t, int64(3), resp.Total)

	require.NotNil(t, repo.lastConditions)
	require.True(t, repo.lastConditions.NeedCount)
	require.Equal(t, "id DESC", repo.lastConditions.Order)
	require.Equal(t, 10, repo.lastConditions.Limit)
	require.Equal(t, 10, repo.lastConditions.Offset)
}
