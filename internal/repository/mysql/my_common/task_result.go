package my_common

import (
	"context"

	"github.com/SisyphusSQ/golib/models/do/base_do"
	"gorm.io/gorm"

	gormv2 "go-starter/internal/lib/gorm"
	"go-starter/internal/lib/log"
)

type TaskResultRepository interface {
	CreateTaskResult(ctx context.Context, res base_do.TaskResult) (err error)
	UpdateByUUID(ctx context.Context, res base_do.TaskResult) (err error)
}

type mysqlTaskResultRepository struct {
	engine *gormv2.Engine
}

func NewTaskResultRepository(engine *gormv2.Engine) TaskResultRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &mysqlTaskResultRepository{
		engine: engine,
	}
}

func (m *mysqlTaskResultRepository) CreateTaskResult(ctx context.Context, res base_do.TaskResult) (err error) {
	err = m.engine.Connect().Session(&gorm.Session{Logger: log.SilentLogger{}}).
		WithContext(ctx).
		Table(base_do.TaskResult{}.TableName()).
		Create(&res).Error
	return
}

func (m *mysqlTaskResultRepository) UpdateByUUID(ctx context.Context, res base_do.TaskResult) (err error) {
	err = m.engine.Connect().Session(&gorm.Session{Logger: log.SilentLogger{}}).
		WithContext(ctx).
		Table(base_do.TaskResult{}.TableName()).
		Where("uuid = ?", res.UUID).
		Updates(&res).Error
	return
}
