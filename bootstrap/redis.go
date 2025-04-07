package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"forum/g"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"net"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", g.Conf.GetString("Redis.Default.Host"), g.Conf.GetInt("Redis.Default.Port")),
		Password: g.Conf.GetString("Redis.Default.Password"),
		DB:       g.Conf.GetInt("Redis.Default.DB"),
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Errorf("redis conn fail:%w", err))
	}
	rdb.AddHook(redisHook{})
	return rdb
}

type redisHook struct{}

func (redisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}
func (redisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		err := next(ctx, cmd)
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Error().Err(err).Msgf("redis err")
			return err
		}
		return nil
	}
}
func (redisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
