package server

import (
	"context"
	"forum/internal/router"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

var _ Server = (*SchedulerServer)(nil)

type SchedulerServer struct {
	server gocron.Scheduler
}

func (t *SchedulerServer) Start(ctx context.Context) error {
	scheduler, err := gocron.NewScheduler(
		gocron.WithStopTimeout(10*time.Second),
		gocron.WithGlobalJobOptions(
			gocron.WithEventListeners(
				gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
					log.Error().Str("脚本名称", jobID.String()).
						Err(err).
						Msg("脚本错误")
				}),
			),
		),
	)
	if err != nil {
		return err
	}
	t.server = scheduler
	s := router.NewScheduler(t.server)
	s.RegisterScheduler()
	t.server.Start()
	log.Info().Msg("sheduler server start ...")
	return nil
}

func (t *SchedulerServer) End(ctx context.Context) error {
	return t.server.Shutdown()
}
