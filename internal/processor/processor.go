package processor

import (
	"github.com/hibiken/asynq"
	"github.com/hiennguyen9874/go-boilerplate/config"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
)

type RedisTaskProcessor struct {
	Server *asynq.Server
	Cfg    *config.Config
	Logger logger.Logger
}

func NewRedisTaskProcessor(server *asynq.Server, cfg *config.Config, logger logger.Logger) RedisTaskProcessor {
	return RedisTaskProcessor{Server: server, Cfg: cfg, Logger: logger}
}
