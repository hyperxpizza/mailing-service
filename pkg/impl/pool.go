package impl

import (
	"context"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	m.logger.Infof("deleting job %s", req.Id)

	err := m.pool.RemoveJob(req.Id)
	if err != nil {
		return nil, status.Error(
			codes.NotFound,
			err.Error(),
		)
	}

	return &emptypb.Empty{}, nil
}
