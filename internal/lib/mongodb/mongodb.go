package mongodb

import (
	"context"

	"github.com/qiniu/qmgo"

	"go-starter/config"
)

func New(c config.MongoDB, coll string) (*qmgo.QmgoClient, error) {
	client, err := qmgo.Open(context.Background(), &qmgo.Config{
		Uri:              "mongodb://" + c.URI,
		Database:         c.Database,
		Coll:             coll,
		MaxPoolSize:      &c.MaxPoolSize,
		MinPoolSize:      &c.MinPoolSize,
		ConnectTimeoutMS: &c.ConnectTimeoutMS,
		SocketTimeoutMS:  &c.SocketTimeoutMS,
		Auth: &qmgo.Credential{
			AuthSource: c.AuthSource,
			Username:   c.User,
			Password:   c.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	err = client.Ping(5)
	if err != nil {
		return nil, err
	}

	return client, nil
}
