package impl

import (
	"fmt"
	"net"

	"github.com/go-redis/redis/v8"
	"github.com/hyperxpizza/mailing-service/pkg/config"
	"github.com/hyperxpizza/mailing-service/pkg/database"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type MailingServiceServer struct {
	cfg    config.Config
	db     *database.Database
	logger logrus.FieldLogger
	rdc    redis.Client
	pb.UnimplementedMailingServiceServer
}

func NewMailingServiceServer() (*MailingServiceServer, error) {
	return &MailingServiceServer{}, nil
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
