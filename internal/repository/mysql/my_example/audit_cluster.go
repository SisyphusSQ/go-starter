package my_example

import (
	"context"

	gormv2 "go-starter/internal/lib/gorm"
	do "go-starter/internal/models/do/mysql/example"
)

type AuditClusterRepository interface {
	GetByID(ctx context.Context, id int64) (cluster do.AuditCluster, err error)
}

type mysqlAuditClusterRepository struct {
	engine *gormv2.Engine
}

func NewAuditClusterRepository(engine *gormv2.Engine) AuditClusterRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &mysqlAuditClusterRepository{engine: engine}
}

func (m *mysqlAuditClusterRepository) GetByID(ctx context.Context, id int64) (cluster do.AuditCluster, err error) {
	err = m.engine.Connect().Table(do.AuditCluster{}.TableName()).Where("id = ?", id).First(&cluster).Error
	return
}
