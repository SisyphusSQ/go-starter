package repository

import (
	"go.uber.org/fx"

	mongo_example_repo "go-starter/internal/repository/mongo/example_repo"
	mysql_example_repo "go-starter/internal/repository/mysql/example_repo"
	"go-starter/internal/repository/mysql/my_common"
)

var Module = fx.Provide(
	my_common.NewConfigKVRepository,
	my_common.NewLarkMsgLogRepository,
	mysql_example_repo.NewUserRepository,
	mongo_example_repo.NewUserRepository,
)
