package impl

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hyperxpizza/mailing-service/pkg/config"
	"github.com/hyperxpizza/mailing-service/pkg/database"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	redisPONG            = "PONG"
	redisConnectionError = "cannot connect to redis"
)

type MailingServiceServer struct {
	cfg        *config.Config
	db         *database.Database
	logger     logrus.FieldLogger
	rdc        *redis.Client
	smtpClient *smtp.Client
	pb.UnimplementedMailingServiceServer
}

func NewMailingServiceServer(lgr logrus.FieldLogger, c *config.Config) (*MailingServiceServer, error) {
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

	return &MailingServiceServer{
		cfg:        c,
		db:         db,
		logger:     lgr,
		rdc:        rdc,
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

func (m *MailingServiceServer) Run() {
	grpcServer := grpc.NewServer()
	pb.RegisterMailingServiceServer(grpcServer, m)

	addr := fmt.Sprintf(":%d", m.cfg.MailingService.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		m.logger.Fatalf("net.Listen failed: %s", err.Error())
	}

	m.logger.Infof("auth service server running on %s:%d", m.cfg.MailingService.Host, m.cfg.MailingService.Port)

	if err := grpcServer.Serve(lis); err != nil {
		m.logger.Fatalf("failed to serve: %s", err.Error())
	}
}
