package service

import (
	"context"
	"go-starter/internal/models/do"
	"go-starter/internal/repository/mysql"
	"time"
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
