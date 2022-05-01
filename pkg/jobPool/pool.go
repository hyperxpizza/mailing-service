package job_pool

import (
	"context"
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
	execJob chan Job
	done    chan string
	jobs    map[string]Job
	gron    gronx.Gronx
	rdc     redis.Client
}

func NewPool(ctx context.Context, logger logrus.FieldLogger) *Pool {
	return &Pool{
		ctx:     ctx,
		wg:      sync.WaitGroup{},
		rwMutex: sync.RWMutex{},
		logger:  logger,
		execJob: make(chan Job),
		done:    make(chan string),
		jobs:    make(map[string]Job),
		gron:    gronx.New(),
	}
}

func (p *Pool) Run() {
	p.logger.Info("Running mailing-service job pool")
	go func() {
		for {
			select {
			case job := <-p.execJob:
				p.wg.Add(1)
				go job.Exec(p.ctx, p.wg, p.logger, job.GetID(), p.done)
			case id := <-p.done:
				go p.cleanupAfterJobDone(id)
			case <-p.ctx.Done():
				return
			}
		}
	}()
	go func() {
		for {
			p.searchForJobs()
		}
	}()

	p.wg.Wait()
}

func (p *Pool) searchForJobs() {

	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	for _, j := range p.jobs {
		if j.GetStatus() == Waiting {
			due, err := p.gron.IsDue(j.GetCron(), time.Now())
			if err != nil {
				continue
			}

			if due {
				p.logger.Infof("job: %s ready to be executed", j.GetID())
				j.SetStatus(Running)
				p.execJob <- j
			}
		}
	}
}

func (p *Pool) AddJob(job Job) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.jobs[job.GetID()] = job
	p.logger.Infof("added job: %s to the pool", job.GetID())
}

func (p *Pool) cleanupAfterJobDone(id string) {
	p.logger.Infof("cleaning up after job: %s", id)
	p.rwMutex.Lock()
	for _, j := range p.jobs {
		if j.GetID() == id {
			j.SetStatus(Waiting)
			break
		}
	}
	p.rwMutex.Unlock()
}

func (p *Pool) RemoveJob(id string) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	delete(p.jobs, id)
}

func (p *Pool) PopulateFromDB(ids []string) {}
