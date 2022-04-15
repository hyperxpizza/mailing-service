package main

import (
	"context"
	"flag"
	"net"
	"testing"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"github.com/hyperxpizza/mailing-service/pkg/impl"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	buffer                = 1024 * 1024
	target                = "bufnet"
	configPathNotSetError = "config path is not set"
)

var lis *bufconn.Listener
var ctx = context.Background()
var configPathOpt = flag.String("config", "", "path to config file")
var loglevelOpt = flag.String("loglevel", "", "loglevel")

func mockGrpcServer(configPath string, secure bool) error {
	lis = bufconn.Listen(buffer)
	server := grpc.NewServer()

	logger := logrus.New()
	if level, err := logrus.ParseLevel(*loglevelOpt); err == nil {
		logger.Level = level
	}

	mailingServiceServer, err := impl.NewMailingServiceServer()
	if err != nil {
		return err
	}

	pb.RegisterMailingServiceServer(server, mailingServiceServer)

	if err := server.Serve(lis); err != nil {
		logger.Fatal(err)
		return err
	}

	return nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func sampleRecipientRequest() *pb.NewMailRecipient {
	return &pb.NewMailRecipient{
		Email:          "hyperxpizza@kernelpanic.live",
		UsersServiceID: 1,
		GroupName:      "Customers",
		Confirmed:      false,
	}
}

func TestMailingServer(t *testing.T) {
	flag.Parse()

	go mockGrpcServer(*configPathOpt, false)

	connection, err := grpc.DialContext(ctx, target, grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)

	defer connection.Close()

	client := pb.NewMailingServiceClient(connection)

	t.Run("Add recipients test", func(t *testing.T) {
		id, err := client.AddRecipient(ctx, sampleRecipientRequest())
		assert.NoError(t, err)
	})
}
