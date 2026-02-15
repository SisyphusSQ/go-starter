package example_repo

import (
	"context"
	"errors"

	gormv2 "go-starter/internal/lib/gorm"
	do "go-starter/internal/models/do/mysql/example_do"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (do.User, error)
	GetByEmail(ctx context.Context, email string) (do.User, error)
	ListByConditions(ctx context.Context, conditions *gormv2.DBConditions) ([]do.User, error)
	Create(ctx context.Context, user *do.User) error
	UpdateByID(ctx context.Context, id int64, updates map[string]any) error
	DeleteByID(ctx context.Context, id int64) error
}

type mysqlUserRepository struct {
	engine *gormv2.Engine
}

func NewUserRepository(engine *gormv2.Engine) UserRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &mysqlUserRepository{engine: engine}
}

func (m *mysqlUserRepository) GetByID(ctx context.Context, id int64) (user do.User, err error) {
	err = m.engine.Connect().WithContext(ctx).Table(do.User{}.TableName()).
		Where("id = ?", id).First(&user).Error
	return
}

func (m *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (user do.User, err error) {
	err = m.engine.Connect().WithContext(ctx).Table(do.User{}.TableName()).
		Where("email = ?", email).First(&user).Error
	return
}

func (m *mysqlUserRepository) ListByConditions(ctx context.Context, conditions *gormv2.DBConditions) (users []do.User, err error) {
	if conditions == nil {
		return nil, errors.New("conditions is nil")
	}

	err = conditions.Fill(m.engine.Connect().WithContext(ctx).Table(do.User{}.TableName())).
		Find(&users).Error
	return
}

func (m *mysqlUserRepository) Create(ctx context.Context, user *do.User) (err error) {
	err = m.engine.Connect().WithContext(ctx).Table(do.User{}.TableName()).
		Create(user).Error
	return
}

func (m *mysqlUserRepository) UpdateByID(ctx context.Context, id int64, updates map[string]any) (err error) {
	err = m.engine.Connect().WithContext(ctx).Table(do.User{}.TableName()).
		Where("id = ?", id).Updates(updates).Error
	return
}

func (m *mysqlUserRepository) DeleteByID(ctx context.Context, id int64) (err error) {
	err = m.engine.Connect().WithContext(ctx).Table(do.User{}.TableName()).
		Where("id = ?", id).Delete(&do.User{}).Error
	return
}
