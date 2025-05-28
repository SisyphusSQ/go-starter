package my_common

import (
	"context"

	"github.com/SisyphusSQ/golib/models/do/base_do"

	gormv2 "go-starter/internal/lib/gorm"
)

type LarkMsgLogRepository interface {
	CreateRecord(ctx context.Context, msg base_do.LarkMsgLog) (err error)
}

type mysqlLarkMsgLogRepository struct {
	engine *gormv2.Engine
}

func NewLarkMsgLogRepository(engine *gormv2.Engine) LarkMsgLogRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &mysqlLarkMsgLogRepository{
		engine: engine,
	}
}

func (m *mysqlLarkMsgLogRepository) CreateRecord(ctx context.Context, msg base_do.LarkMsgLog) (err error) {
	err = m.engine.Connect().WithContext(ctx).Table(base_do.LarkMsgLog{}.TableName()).
		Create(&msg).Error
	return
}
