package impl

import (
	"context"
	"errors"

	"github.com/hyperxpizza/mailing-service/pkg/customErrors"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"github.com/hyperxpizza/mailing-service/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MailingServiceServer) CreateGroup(ctx context.Context, req *pb.MailingServiceNewGroup) (*pb.MailingServiceID, error) {
	var id pb.MailingServiceID
	err := utils.ValidateNewGroup(req.Name)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	id.Id, err = m.db.InsertGroup(req.Name)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	return &id, nil
}

func (m *MailingServiceServer) GetGroups(ctx context.Context, req *emptypb.Empty) (*pb.MailGroups, error) {
	var groups pb.MailGroups
	var err error
	groups.Groups, err = m.db.GetGroups()
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	return &groups, nil
}

func (m *MailingServiceServer) DeleteGroup(ctx context.Context, req *pb.MailingServiceID) (*emptypb.Empty, error) {

	err := m.db.DeleteGroup(req.Id)
	if err != nil {
		var gErr *customErrors.NotFoundError
		if errors.As(err, &gErr) {
			return nil, status.Error(
				codes.NotFound,
				err.Error(),
			)
		}
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	return &emptypb.Empty{}, nil
}

func (m *MailingServiceServer) UpdateGroupName(ctx context.Context, req *pb.MailingServiceNewGroup) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
