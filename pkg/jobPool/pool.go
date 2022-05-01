package job_pool

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Descriptior struct {
	Type    string
	Due     time.Time
	Retries int64
	Args    map[string]interface{}
}

type Job interface {
	Exec(context.Context, sync.WaitGroup)
	GetDescriptor() *Descriptior
	GetStatus() JobStatus
	SetStatus(JobStatus)
	GetID() string
}

type Pool struct {
	ctx     context.Context
	wg      sync.WaitGroup
	rwMutex sync.RWMutex
	logger  logrus.FieldLogger
	execJob chan Job
	done    chan string
	jobs    map[string]Job
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
	}
}

func (p *Pool) Run() {
	var wg sync.WaitGroup
	wg.Add(2)
	p.logger.Info("Running mailing-service job pool")
	go func() {
		for {
			select {
			case <-p.execJob:
				job := <-p.execJob
				p.wg.Add(1)
				go job.Exec(p.ctx, p.wg)

			case <-p.ctx.Done():
				wg.Done()
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				wg.Done()
				return
			default:
				jobs, err := p.searchForJobs()
				if err == nil {
					for _, j := range jobs {
						p.execJob <- j
					}
				}
			}
		}
	}()
	wg.Wait()
	p.wg.Wait()
}

func (p *Pool) searchForJobs() ([]Job, error) {

	if len(p.jobs) == 0 {
		return nil, errors.New("no jobs")
	}

	p.rwMutex.RLock()
	defer p.rwMutex.Unlock()

	var jobs []Job
	for _, j := range p.jobs {
		now := time.Now()
		desc := j.GetDescriptor()
		if desc.Due.Equal(now) {

			jobs = append(jobs, j)
		}
	}

	return jobs, nil
}

func (p *Pool) AddRecurrentJob(job Job) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.jobs[job.GetID()] = job
}

func (p *Pool) RemoveJob(id string) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	delete(p.jobs, id)
}

func (p *Pool) cleanup() {

}
