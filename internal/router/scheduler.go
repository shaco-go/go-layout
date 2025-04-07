package router

import (
	"context"
	"fmt"
	"forum/g"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NewScheduler(s gocron.Scheduler) *Scheduler {
	return &Scheduler{
		server: s,
	}
}

type Scheduler struct {
	server gocron.Scheduler
}

// RegisterScheduler 注册调度
func (s *Scheduler) RegisterScheduler() {
}

// 单线程调度,如果前面有任务会跳过
func (s *Scheduler) singleTask(cron, jobName string, fn func(context context.Context) error) {
	rid := uuid.New()
	ctx := context.WithValue(context.Background(), g.RID, rid.String())
	_, err := s.server.NewJob(
		gocron.CronJob(cron, true),
		gocron.NewTask(func() error {
			return fn(ctx)
		}),
		gocron.WithName(jobName),
		gocron.WithIdentifier(rid),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		panic(fmt.Sprintf("%s 启动失败", jobName))
	}
	log.Debug().
		Str("脚本名称", jobName).
		Str("cron", cron).
		Msg("success")
}
