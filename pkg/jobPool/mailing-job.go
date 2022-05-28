package job_pool

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/adhocore/gronx"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	cronNotValidError = "provided cron expression is not valid"
)

var gron = gronx.New()

type ExecutionFunction func(context.Context, logrus.FieldLogger, chan string)

type Job interface {
	Exec(context.Context, sync.WaitGroup, logrus.FieldLogger, chan string)
	GetStatus() JobStatus
	SetStatus(JobStatus)
	GetID() string
	GetCron() string
	GetLastExecuted() time.Time
	SetLastExecuted(time.Time)
}

type MailingJob struct {
	id           string
	cron         string
	jobType      JobType
	status       JobStatus
	lastExecuted time.Time
	ExecFn       ExecutionFunction
}

func NewMailingJob(cron string, jType JobType) (*MailingJob, error) {
	if !gron.IsValid(cron) {
		return nil, errors.New(cronNotValidError)
	}

	uuID := uuid.New()

	return &MailingJob{
		id:      uuID.String(),
		cron:    cron,
		jobType: jType,
		status:  Waiting,
	}, nil
}

func (m *MailingJob) SetExecFN(fn ExecutionFunction) {
	m.ExecFn = fn
}

func (m *MailingJob) GetLastExecuted() time.Time {
	return m.lastExecuted
}

func (m *MailingJob) SetLastExecuted(t time.Time) {
	m.lastExecuted = t
}

func (m *MailingJob) Exec(ctx context.Context, logger logrus.FieldLogger, done chan string) {
	m.ExecFn(ctx, logger, done)
}

func (m *MailingJob) GetStatus() JobStatus {
	return m.status
}

func (m *MailingJob) SetStatus(s JobStatus) {
	m.status = s
}

func (m *MailingJob) GetID() string {
	return m.id
}

func (m *MailingJob) GetCron() string {
	return m.cron
}
