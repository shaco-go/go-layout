package server

import (
	"context"
	"fmt"
	"forum/g"
	"forum/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

var _ Server = (*HttpServer)(nil)

type HttpServer struct {
	server *http.Server
}

func (h *HttpServer) Start(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)
	r := router.RegisterHttpRouter()
	port := g.Conf.GetInt64("Http.Port")
	// 跨域,日志,错误处理的中间件
	h.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
	domain := "http://127.0.0.1"
	if g.Conf.GetString("Http.Domain") != "" {
		domain = g.Conf.GetString("Http.Domain")
	}
	log.Info().Msgf("http server start: %s:%d", domain, port)
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "http listen err")
	}
	return nil
}

func (h *HttpServer) End(ctx context.Context) error {
	err := h.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
