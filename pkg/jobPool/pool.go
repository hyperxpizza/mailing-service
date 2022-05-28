package job_pool

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/adhocore/gronx"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Pool struct {
	ctx     context.Context
	wg      sync.WaitGroup
	rwMutex sync.RWMutex
	logger  logrus.FieldLogger
	execJob chan *MailingJob
	done    chan string
	jobs    map[string]*MailingJob
	gron    gronx.Gronx
	rdc     *redis.Client
	ticker  *time.Ticker
}

func NewPool(ctx context.Context, logger logrus.FieldLogger, rdc *redis.Client) (*Pool, error) {
	pool := Pool{
		ctx:     ctx,
		wg:      sync.WaitGroup{},
		rwMutex: sync.RWMutex{},
		logger:  logger,
		execJob: make(chan *MailingJob),
		done:    make(chan string),
		jobs:    make(map[string]*MailingJob),
		gron:    gronx.New(),
		rdc:     rdc,
		ticker:  time.NewTicker(time.Minute),
	}

	err := pool.LoadJobsFromDB()
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

func (p *Pool) Run() {
	p.wg.Add(2)
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				p.wg.Done()
				return
			case t := <-p.ticker.C:
				p.searchForJobs(t)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				p.wg.Done()
				return
			case job := <-p.execJob:
				p.wg.Add(1)
				p.logger.Infof("executing job: %s", job.id)
				go job.Exec(p.ctx, p.logger, p.done)
			case id := <-p.done:
				p.wg.Done()
				go p.cleanup(id)
			}
		}
	}()
	p.logger.Infoln("Running pool...")
	p.wg.Wait()
	p.logger.Infoln("Shutting down...")
}

func (p *Pool) AddJob(job *MailingJob) error {

	if err := p.rdc.Set(p.ctx, job.id, job, 0).Err(); err != nil {
		p.logger.Infof("failed to insert job %s into the database: %s", job.id, err.Error())
		return err
	}

	p.rwMutex.Lock()
	p.jobs[job.id] = job
	p.rwMutex.Unlock()

	p.logger.Infof("Added new job: %s", job.id)

	return nil
}

func (p *Pool) RemoveJob(id string) error {

	if err := p.rdc.Del(p.ctx, id).Err(); err != nil {
		p.logger.Infof("failed to detele job %s from the database: %s", id, err.Error())
		return err
	}

	p.rwMutex.Lock()
	delete(p.jobs, id)
	p.rwMutex.Unlock()

	return nil
}

func (p *Pool) LoadJobsFromDB() error {
	p.logger.Infoln("Loading jobs from the database...")
	counter := 0
	iter := p.rdc.Scan(p.ctx, 0, "prefix:*", 0).Iterator()
	for iter.Next(p.ctx) {
		fmt.Println(iter.Val())

		result, err := p.rdc.Get(p.ctx, iter.Val()).Result()
		if err != nil {
			p.logger.Debugf("getting value for key %s failed: %s", iter.Val(), err.Error())
			continue
		}

		var job MailingJob
		err = json.Unmarshal([]byte(result), &job)
		if err != nil {
			p.logger.Debugf("unmarshaling failed: %s", err.Error())
			continue
		}

		p.rwMutex.Lock()
		p.jobs[job.id] = &job
		p.rwMutex.Unlock()
		counter++
	}

	if err := iter.Err(); err != nil {
		return err
	}

	p.logger.Infof("loaded %d jobs from the database", counter)
	return nil

}

func (p *Pool) searchForJobs(t time.Time) {
	p.rwMutex.RLock()
	for _, job := range p.jobs {
		if job.status == Waiting {

			due, err := p.gron.IsDue(job.cron, t)
			if err != nil {
				p.logger.Debugf("job with id: %s has an invalid cron: %s", job.id, err.Error())
				continue
			}

			if due {
				job.SetStatus(Running)
				p.execJob <- job
			}

		}
	}
	p.rwMutex.RUnlock()
}

func (p *Pool) GetJobs() map[string]*MailingJob {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return p.jobs
}

func (p *Pool) cleanup(id string) {
	p.rwMutex.Lock()
	job := p.jobs[id]
	job.status = Waiting
	if job.jobType == OneTime {
		p.RemoveJob(id)
	}
	p.rwMutex.Unlock()
}
