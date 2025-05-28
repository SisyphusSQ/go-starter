package repository

import (
	"go.uber.org/fx"

	"go-starter/internal/repository/mysql/my_common"
	"go-starter/internal/repository/mysql/my_example"
)

var Module = fx.Provide(
	my_common.NewTaskResultRepository,
	my_common.NewConfigKVRepository,
	my_common.NewLarkMsgLogRepository,
	my_example.NewAuditClusterRepository,
)
