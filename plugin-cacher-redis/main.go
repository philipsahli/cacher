package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/go-plugin"
	"github.com/philipsahli/cacher/shared"

	"github.com/redis/go-redis/v9"
)

// Here is a real implementation of KV that writes to redis with
// the key name.
type KV struct {
	Logger hclog.Logger
}

func (kv KV) Get(key string) ([]byte, error) {
	ctx := context.Background()

	ip := "localhost"
	port := "6379"

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", ip, port),
		Password: "",
		DB:       0,
	})

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return []byte{}, err
	}

	kv.Logger.Info(fmt.Sprintf("Key %s read from backend (redis)", key))

	return []byte(val), nil
}

func (kv KV) Put(key string, value []byte) error {
	ctx := context.Background()

	ip := "localhost"
	port := "6379"

	rdb := redis.NewClient(&redis.Options{
		//Addr:     fmt.Sprintf("localhost:%s", port),
		Addr:     fmt.Sprintf("%s:%s", ip, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	kv.Logger.Info(fmt.Sprintf("Key %s written to backend (redis)", key))

	return nil

}

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Info,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	logger.Debug("message from plugin-cacher-redis")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"kv": &shared.KVGRPCPlugin{Impl: &KV{
				Logger: logger,
			}},
		},

		GRPCServer: plugin.DefaultGRPCServer,
	})

}
