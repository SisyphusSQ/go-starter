package repository

import (
	"go.uber.org/fx"

	"go-starter/internal/repository/mysql"
)

var Module = fx.Provide(
	mysql.NewAuditClusterRepository,
	mysql.NewTaskResultRepository,
)
