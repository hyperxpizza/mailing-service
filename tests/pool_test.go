package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hyperxpizza/mailing-service/pkg/config"
	job_pool "github.com/hyperxpizza/mailing-service/pkg/jobPool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {

	c, err := config.NewConfig(*configPathOpt)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	if level, err := logrus.ParseLevel(*loglevelOpt); err == nil {
		logger.Level = level
	}

	execFN := func(ctx context.Context, wg sync.WaitGroup, logger logrus.FieldLogger, id string, done chan string) {
		defer wg.Done()
		for i := 1; i < 7; i++ {
			logger.Infof("%d/5 executing job: %s", i, id)
			time.Sleep(10 * time.Second)
		}
		logger.Infof("finished job: %s", id)
		done <- id
	}

	rdc := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		DB:   int(c.Redis.DB),
	})

	job, err := job_pool.NewMailingJob("*/2 * * * *", job_pool.Recurrent, execFN)
	assert.NoError(t, err)

	ctx := context.Background()
	pool := job_pool.NewPool(ctx, logger, rdc)
	go pool.Run()
	pool.AddJob(job)

	time.Sleep(5 * time.Minute)

}
