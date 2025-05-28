package my_common

import (
	"context"

	gormv2 "go-starter/internal/lib/gorm"

	"github.com/SisyphusSQ/golib/models/do/base_do"
)

type ConfigKVRepository interface {
	GetByKey(ctx context.Context, k string) (kv base_do.ConfigKV, err error)
}

type mysqlConfigKVRepo struct {
	engine *gormv2.Engine
}

func NewConfigKVRepository(engine *gormv2.Engine) ConfigKVRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &mysqlConfigKVRepo{
		engine: engine,
	}
}

func (m *mysqlConfigKVRepo) GetByKey(ctx context.Context, k string) (kv base_do.ConfigKV, err error) {
	err = m.engine.Connect().WithContext(ctx).Table(base_do.ConfigKV{}.TableName()).
		Where("config_key = ?", k).
		Find(&kv).Error
	return
}
