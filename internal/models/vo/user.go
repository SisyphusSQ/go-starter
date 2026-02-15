package vo

import (
	mongoDo "go-starter/internal/models/do/mongo/example_do"
	mysqlDo "go-starter/internal/models/do/mysql/example_do"
)

type CreateUserReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire int    `json:"expire"`
}

type UserListResp struct {
	Total int64          `json:"total"`
	List  []mysqlDo.User `json:"list"`
}

type UserMongoListResp struct {
	Total int64           `json:"total"`
	List  []*mongoDo.User `json:"list"`
}

type UserIDResp struct {
	ID int64 `json:"id"`
}

type UserMongoIDResp struct {
	ID string `json:"id"`
}
