package impl

import (
	"context"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MailingServiceServer) AddJob(ctx context.Context, req *pb.JobRequest) (*pb.JobID, error) {
	var id pb.JobID
	return &id, nil
}

func (m *MailingServiceServer) JobStream(req *emptypb.Empty, stream pb.MailingService_JobStreamServer) error {

	return nil
}

func (m *MailingServiceServer) DeleteJob(ctx context.Context, req *pb.JobID) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
