package impl

import (
	"context"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MailingServiceServer) CreateGroup(ctx context.Context, req *pb.MailingServiceNewGroup) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *MailingServiceServer) GetGroups(ctx context.Context, req *emptypb.Empty) (*pb.MailGroups, error) {
	var groups pb.MailGroups
	return &groups, nil
}

func (m *MailingServiceServer) DeleteGroup(ctx context.Context, req *pb.MailingServiceID) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *MailingServiceServer) UpdateGroupName(ctx context.Context, req *pb.MailingServiceNewGroup) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
