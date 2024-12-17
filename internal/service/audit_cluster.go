package service

import (
	"context"
	"time"

	do "go-starter/internal/models/do/mysql/audit"
	"go-starter/internal/repository/mysql"
)

type (
	AuditClusterService interface {
		GetByID(ctx context.Context, id int64) (cluster do.AuditCluster, err error)
	}

	AuditClusterServiceImpl struct {
		repo           mysql.AuditClusterRepository
		contextTimeout time.Duration
	}
)

func NewAuditClusterService(auditClusterRepo mysql.AuditClusterRepository, timeout time.Duration) AuditClusterService {
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
