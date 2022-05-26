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
	go func() {
		for {
			select {
			case t := <-p.ticker.C:
				p.searchForJobs(t)
			}
		}
	}()

	go func() {
		for {
			select {
			case job := <-p.execJob:
				p.wg.Add(1)
				go job.Exec(p.ctx, p.logger, p.done)
			case id := <-p.done:
				p.wg.Done()
				go p.cleanup(id)
			}
		}
	}()
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
	defer p.rwMutex.Unlock()

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
				p.logger.Debugf("")
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
	p.jobs[id].status = Waiting
}
