package server

import (
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"time"
)

type Server interface {
	Start(ctx context.Context) error
	End(ctx context.Context) error
}

func Run() {
	var group = []Server{
		new(SchedulerServer),
		new(HttpServer),
	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		for _, item := range group {
			go func() {
				err := item.Start(ctx)
				if err != nil {
					log.Fatal().Err(err).Send()
				}
			}()
		}
	}

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Info().Msgf("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var eg errgroup.Group
	for _, item := range group {
		eg.Go(func() error {
			return item.End(ctx)
		})
	}
	err := eg.Wait()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msg("Server stop ...")
}
