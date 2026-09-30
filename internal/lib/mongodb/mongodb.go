package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/SisyphusSQ/go-starter/v2/config"
	"github.com/SisyphusSQ/go-starter/v2/internal/health"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/fx"
)

type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

func New(lifecycle fx.Lifecycle, c config.Config) (*Client, error) {
	cfg := c.MongoDB
	result := &Client{}
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			opts := options.Client().ApplyURI(cfg.URI).SetConnectTimeout(cfg.ConnectTimeout).SetServerSelectionTimeout(cfg.ConnectTimeout).SetMaxPoolSize(cfg.MaxPoolSize).SetMinPoolSize(cfg.MinPoolSize)
			if cfg.User != "" {
				opts.SetAuth(options.Credential{AuthSource: cfg.AuthSource, Username: cfg.User, Password: cfg.Password})
			}
			client, err := mongo.Connect(opts)
			if err != nil {
				return fmt.Errorf("configure mongodb client: %w", err)
			}
			if err := client.Ping(ctx, nil); err != nil {
				closeCtx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
				defer cancel()
				_ = client.Disconnect(closeCtx)
				return fmt.Errorf("ping mongodb: %w", err)
			}
			result.client = client
			result.database = client.Database(cfg.Database)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if result.client == nil {
				return nil
			}
			return result.client.Disconnect(ctx)
		},
	})
	return result, nil
}
func Readiness(c *Client) health.Check {
	return health.Check{Name: "mongodb", Enabled: true, Ping: func(ctx context.Context) error {
		if c.client == nil {
			return errors.New("mongodb not started")
		}
		return c.client.Ping(ctx, nil)
	}}
}
func (c *Client) Collection(name string) *mongo.Collection { return c.database.Collection(name) }
