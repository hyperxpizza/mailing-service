package main

import (
	"context"
	"flag"
	"net"
	"testing"

	"github.com/hyperxpizza/mailing-service/pkg/config"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"github.com/hyperxpizza/mailing-service/pkg/impl"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	buffer                = 1024 * 1024
	target                = "bufnet"
	configPathNotSetError = "config path is not set"
	sampleGroupName       = "CUSTOMERS2"
)

var lis *bufconn.Listener
var ctx = context.Background()
var configPathOpt = flag.String("config", "", "path to config file")
var loglevelOpt = flag.String("loglevel", "info", "loglevel")
var deleteOpt = flag.Bool("delete", true, "delete records after test?")

func mockGrpcServer(configPath string, secure bool) error {
	lis = bufconn.Listen(buffer)
	server := grpc.NewServer()

	logger := logrus.New()
	if level, err := logrus.ParseLevel(*loglevelOpt); err == nil {
		logger.Level = level
	}

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	mailingServiceServer, err := impl.NewMailingServiceServer(logger, cfg)
	if err != nil {
		return err
	}

	pb.RegisterMailingServiceServer(server, mailingServiceServer)

	if err := server.Serve(lis); err != nil {
		panic(err)
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
		GroupName:      sampleGroupName,
		Confirmed:      false,
	}
}

func sampleNewGroup() *pb.MailingServiceNewGroup {
	return &pb.MailingServiceNewGroup{
		Name: sampleGroupName,
	}
}

//go test -v ./tests/ --run TestMailingServer --config=/home/hyperxpizza/dev/golang/reusable-microservices/mailing-service/config.dev.json
func TestMailingServer(t *testing.T) {
	flag.Parse()

	go mockGrpcServer(*configPathOpt, false)

	connection, err := grpc.DialContext(ctx, target, grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)

	defer connection.Close()

	client := pb.NewMailingServiceClient(connection)

	//insertAndGetRecipient := func

	t.Run("Add recipient test", func(t *testing.T) {

		gId, err := client.CreateGroup(ctx, sampleNewGroup())
		assert.NoError(t, err)

		sampleRecipient := sampleRecipientRequest()

		id, err := client.AddRecipient(ctx, sampleRecipient)
		assert.NoError(t, err)

		recipient, err := client.GetRecipient(ctx, id)
		assert.NoError(t, err)

		assert.Equal(t, recipient.Email, sampleRecipient.Email)
		assert.Equal(t, recipient.UsersServiceID, sampleRecipient.UsersServiceID)
		assert.Equal(t, recipient.Confirmed, sampleRecipient.Confirmed)

		checkGroup := func(id int64, name string, group []*pb.MailGroup) {
			assert.Equal(t, 1, len(group))
			assert.Equal(t, id, group[0].Id)
			assert.Equal(t, name, group[0].Name)
		}

		checkGroup(gId.Id, sampleGroupName, recipient.Groups)

		if *deleteOpt {
			_, err = client.RemoveRecipient(ctx, id)
			assert.NoError(t, err)

			_, err = client.DeleteGroup(ctx, gId)
			assert.NoError(t, err)
		}
	})

	t.Run("Add group test", func(t *testing.T) {
		sg := sampleNewGroup()
		gId, err := client.CreateGroup(ctx, sg)
		assert.NoError(t, err)

		groups, err := client.GetGroups(ctx, &emptypb.Empty{})
		assert.NoError(t, err)
		assert.Greater(t, len(groups.Groups), 0)

		var group *pb.MailGroup
		group = nil
		for _, g := range groups.Groups {
			if g.Id == gId.Id {
				group = g
			}
		}

		assert.Equal(t, group.Name, sg.Name)

		if *deleteOpt {
			_, err := client.DeleteGroup(ctx, gId)
			assert.NoError(t, err)
		}
	})

	t.Run("Update group test", func(t *testing.T) {
		sg := sampleNewGroup()
		gId, err := client.CreateGroup(ctx, sg)
		assert.NoError(t, err)

		updatedGroupName := "updated group name"
		req := pb.UpdateGroupRequest{
			Id:      gId.Id,
			NewName: updatedGroupName,
		}

		_, err = client.UpdateGroupName(ctx, &req)
		assert.NoError(t, err)

		updatedGroup, err := client.GetGroup(ctx, gId)
		assert.NoError(t, err)
		assert.Equal(t, updatedGroup.Name, updatedGroupName)

		if *deleteOpt {
			_, err := client.DeleteGroup(ctx, gId)
			assert.NoError(t, err)
		}
	})
}
