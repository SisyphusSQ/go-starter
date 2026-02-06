package example_srv

import (
	"context"
	"time"

	do "go-starter/internal/models/do/mysql/example"
	"go-starter/internal/repository/mysql/my_example"
)

type (
	AuditClusterService interface {
		GetByID(ctx context.Context, id int64) (cluster do.AuditCluster, err error)
	}

	AuditClusterServiceImpl struct {
		repo           my_example.AuditClusterRepository
		contextTimeout time.Duration
	}
)

func NewAuditClusterService(auditClusterRepo my_example.AuditClusterRepository, timeout time.Duration) AuditClusterService {
	if auditClusterRepo == nil {
		panic("AuditClusterRepository is nil")
	}
	if timeout == 0 {
		panic("Timeout is empty")
	}
	return &AuditClusterServiceImpl{
		repo:           auditClusterRepo,
		contextTimeout: timeout,
	}
}

func (s *AuditClusterServiceImpl) GetByID(ctx context.Context, id int64) (cluster do.AuditCluster, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	cluster, err = s.repo.GetByID(ctx, id)
	return
}
