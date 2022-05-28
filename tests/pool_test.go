package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hyperxpizza/mailing-service/pkg/config"
	job_pool "github.com/hyperxpizza/mailing-service/pkg/jobPool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func sampleJob() (*job_pool.MailingJob, error) {
	job, err := job_pool.NewMailingJob("*/5 * * * *", job_pool.Recurrent)
	if err != nil {
		return nil, err
	}
	execFN := func(ctx context.Context, logger logrus.FieldLogger, done chan string) {
		for i := 1; i < 6; i++ {
			logger.Infof("%d/5 executing job: %s", i, job.GetID())
			time.Sleep(5 * time.Second)
		}
		logger.Infof("finished job: %s", job.GetID())
		done <- job.GetID()
	}
	job.SetExecFN(execFN)

	return job, nil
}

// go test -v ./tests --run TestPool --config=/home/hyperxpizza/dev/golang/reusable-microservices/mailing-service/config.dev.json --loglevel=debug
func TestPool(t *testing.T) {

	c, err := config.NewConfig(*configPathOpt)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	if level, err := logrus.ParseLevel(*loglevelOpt); err == nil {
		logger.Level = level
	}

	rdc := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		DB:   int(c.Redis.DB),
	})

	ctx, cancel := context.WithCancel(context.Background())
	pool, err := job_pool.NewPool(ctx, logger, rdc)
	assert.NoError(t, err)
	go pool.Run()

	t.Run("Test add job", func(t *testing.T) {
		job, err := sampleJob()
		assert.NoError(t, err)

		err = pool.AddJob(job)
		assert.NoError(t, err)
	})

	t.Run("Test delete job", func(t *testing.T) {
		job, err := sampleJob()
		assert.NoError(t, err)

		err = pool.AddJob(job)
		assert.NoError(t, err)

		err = pool.RemoveJob(job.GetID())
		assert.NoError(t, err)

	})

	cancel()
}
