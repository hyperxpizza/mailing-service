package main

import (
	"context"
	"flag"

	"github.com/hyperxpizza/mailing-service/pkg/config"
	"github.com/hyperxpizza/mailing-service/pkg/impl"
	"github.com/sirupsen/logrus"
)

var configPathOpt = flag.String("config", "", "path to config file")

func main() {
	flag.Parse()
	if *configPathOpt == "" {
		panic("config path empty")
	}

	cfg, err := config.NewConfig(*configPathOpt)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	if level, err := logrus.ParseLevel(cfg.MailingService.Loglevel); err == nil {
		logger.Level = level
	}

	ctx := context.Background()
	server, err := impl.NewMailingServiceServer(ctx, logger, cfg)
	if err != nil {
		panic(err)
	}

	if err := server.Run(); err != nil {
		panic(err)
	}

}
