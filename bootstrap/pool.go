package bootstrap

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
)

func InitPool() *ants.Pool {
	pool, err := ants.NewPool(
		ants.DefaultAntsPoolSize,
		ants.WithLogger(&poolLogger{}),
	)
	if err != nil {
		panic(fmt.Errorf("协程池创建失败:%w", err))
	}
	return pool
}

type poolLogger struct {
}

func (i *poolLogger) Printf(format string, args ...any) {
	log.Error().Msgf("默认协程池:%s", fmt.Sprintf(format, args...))
}
