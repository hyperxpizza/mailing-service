package job_pool

import (
	"context"
	"encoding/json"
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

func NewPool(ctx context.Context, logger logrus.FieldLogger, rdc *redis.Client) *Pool {
	return &Pool{
		ctx:     ctx,
		wg:      sync.WaitGroup{},
		rwMutex: sync.RWMutex{},
		logger:  logger,
		execJob: make(chan *MailingJob),
		done:    make(chan string),
		jobs:    make(map[string]*MailingJob),
		gron:    gronx.New(),
		rdc:     rdc,
		ticker:  time.NewTicker(time.Second),
	}
}

func (p *Pool) Run() {
	p.wg.Add(2)
	go func() {
		for {
			select {

			case <-p.ctx.Done():
				p.wg.Done()
				break

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
				break
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
	p.wg.Wait()
}

func (p *Pool) AddJob(job *MailingJob) error {

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	go p.rdc.Set(p.ctx, job.id, data, 0)

	p.rwMutex.Lock()
	p.jobs[job.id] = job
	p.rwMutex.Unlock()

	return nil
}

func (p *Pool) RemoveJob(id string) error {
	p.rwMutex.Lock()
	delete(p.jobs, id)
	p.rwMutex.Unlock()

	go p.rdc.Del(p.ctx, id)

	return nil
}

func (p *Pool) LoadJobsFromDB() {
	var jobs []*MailingJob

	for _, j := range jobs {
		p.jobs[j.id] = j
	}
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

func (p *Pool) cleanup(id string) {
	p.rwMutex.Lock()
	job := p.jobs[id]
	job.status = Waiting
	if job.jobType == OneTime {
		p.RemoveJob(id)
	}
	p.rwMutex.Unlock()
}
