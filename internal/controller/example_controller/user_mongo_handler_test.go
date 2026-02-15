package example_controller

import (
	"context"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	mongoDo "go-starter/internal/models/do/mongo/example_do"
	mysqlDo "go-starter/internal/models/do/mysql/example_do"
	"go-starter/internal/models/vo"
)

type mockUserService struct{}

func (m *mockUserService) GetByID(ctx context.Context, id int64) (mysqlDo.User, error) {
	return mysqlDo.User{}, nil
}

func (m *mockUserService) List(ctx context.Context, page, pageSize int) (vo.UserListResp, error) {
	return vo.UserListResp{}, nil
}

func (m *mockUserService) Create(ctx context.Context, req vo.CreateUserReq) (mysqlDo.User, error) {
	return mysqlDo.User{}, nil
}

func (m *mockUserService) Login(ctx context.Context, req vo.LoginReq) (vo.LoginResp, error) {
	return vo.LoginResp{}, nil
}

func (m *mockUserService) Logout(ctx context.Context, userID int64) (vo.UserIDResp, error) {
	return vo.UserIDResp{}, nil
}

func (m *mockUserService) Update(ctx context.Context, id int64, req vo.UpdateUserReq) (vo.UserIDResp, error) {
	return vo.UserIDResp{}, nil
}

func (m *mockUserService) Delete(ctx context.Context, id int64) (vo.UserIDResp, error) {
	return vo.UserIDResp{}, nil
}

type mockUserMongoService struct {
	available bool
}

func (m *mockUserMongoService) IsAvailable() bool {
	return m.available
}

func (m *mockUserMongoService) GetByID(ctx context.Context, id string) (*mongoDo.User, error) {
	return &mongoDo.User{}, nil
}

func (m *mockUserMongoService) List(ctx context.Context, page, pageSize int) (vo.UserMongoListResp, error) {
	return vo.UserMongoListResp{}, nil
}

func (m *mockUserMongoService) Create(ctx context.Context, req vo.CreateUserReq) (*mongoDo.User, error) {
	return &mongoDo.User{}, nil
}

func (m *mockUserMongoService) Update(ctx context.Context, id string, req vo.UpdateUserReq) (vo.UserMongoIDResp, error) {
	return vo.UserMongoIDResp{}, nil
}

func (m *mockUserMongoService) Delete(ctx context.Context, id string) (vo.UserMongoIDResp, error) {
	return vo.UserMongoIDResp{}, nil
}

func TestInitUserMongoControllerSkipRoutesWhenUnavailable(t *testing.T) {
	e := echo.New()

	InitUserController(e, &mockUserService{})
	InitUserMongoController(e, &mockUserMongoService{available: false})

	var hasMySQLUserRoutes bool
	var hasMongoUserRoutes bool
	for _, route := range e.Routes() {
		if strings.HasPrefix(route.Path, "/mysql/users") {
			hasMySQLUserRoutes = true
		}
		if strings.HasPrefix(route.Path, "/mongo/users") {
			hasMongoUserRoutes = true
		}
	}

	require.True(t, hasMySQLUserRoutes)
	require.False(t, hasMongoUserRoutes)
}
