package impl

import (
	"github.com/hyperxpizza/mailing-service/pkg/database"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"github.com/sirupsen/logrus"
)

type MailingServiceServer struct {
	db     *database.Database
	logger logrus.FieldLogger
	pb.UnimplementedMailingServiceServer
}

func NewMailingServiceServer() {}

func (m *MailingServiceServer) Run() {}
