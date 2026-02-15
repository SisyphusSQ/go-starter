package example_do

import "github.com/qiniu/qmgo/field"

type User struct {
	field.DefaultField `bson:",inline"`
	Name               string `json:"name" bson:"name"`
	Email              string `json:"email" bson:"email"`
	IsDelete           bool   `json:"isDelete" bson:"isDelete"`
}

func (User) Collection() string {
	return "users"
}
