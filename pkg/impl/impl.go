package impl

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hyperxpizza/mailing-service/pkg/config"
	"github.com/hyperxpizza/mailing-service/pkg/database"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	job_pool "github.com/hyperxpizza/mailing-service/pkg/jobPool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	redisPONG            = "PONG"
	redisConnectionError = "cannot connect to redis"
)

type MailingServiceServer struct {
	ctx        context.Context
	cfg        *config.Config
	db         *database.Database
	logger     logrus.FieldLogger
	rdc        *redis.Client
	smtpClient *smtp.Client
	pool       *job_pool.Pool
	poolCtx    context.Context
	pb.UnimplementedMailingServiceServer
}

func NewMailingServiceServer(ctx context.Context, lgr logrus.FieldLogger, c *config.Config) (*MailingServiceServer, error) {
	db, err := database.Connect(c)
	if err != nil {
		return nil, err
	}

	rdc := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		DB:   int(c.Redis.DB),
	})

	if err := checkRedisConnection(rdc); err != nil {
		return nil, err
	}

	smtpAddr := fmt.Sprintf("%s:%d", c.SMTP.Host, c.SMTP.Port)
	to := time.Duration(2) * time.Second
	smtpConn, err := net.DialTimeout("tcp", smtpAddr, to)
	if err != nil {
		return nil, err
	}

	smtpClient, err := smtp.NewClient(smtpConn, c.SMTP.Host)
	if err != nil {
		return nil, err
	}

	pool := job_pool.NewPool(ctx, lgr, rdc)

	return &MailingServiceServer{
		cfg:        c,
		db:         db,
		logger:     lgr,
		rdc:        rdc,
		pool:       pool,
		smtpClient: smtpClient,
	}, nil
}

func checkRedisConnection(rdc *redis.Client) error {
	stat, err := rdc.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	if stat != redisPONG {
		return errors.New(redisConnectionError)
	}

	return nil
}

func (m *MailingServiceServer) WithPoolCtx(ctx context.Context) *MailingServiceServer {
	m.poolCtx = ctx
	return m
}

func (m *MailingServiceServer) Run() error {

	cert, err := tls.LoadX509KeyPair(m.cfg.TLS.CertPath, m.cfg.TLS.KeyPath)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMailingServiceServer(grpcServer, m)

	addr := fmt.Sprintf(":%d", m.cfg.MailingService.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		m.logger.Fatalf("net.Listen failed: %s", err.Error())
		return err
	}

	m.logger.Infof("auth service server running on %s:%d", m.cfg.MailingService.Host, m.cfg.MailingService.Port)

	if err := grpcServer.Serve(lis); err != nil {
		m.logger.Fatalf("failed to serve: %s", err.Error())
		return err
	}

	return nil
}
