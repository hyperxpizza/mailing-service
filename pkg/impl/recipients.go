package impl

import (
	"context"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MailingServiceServer) AddRecipient(ctx context.Context, req *pb.NewMailRecipient) (*pb.MailingServiceID, error) {
	var id pb.MailingServiceID
	return &id, nil
}

func (m *MailingServiceServer) RemoveRecipiet(ctx context.Context, req *pb.MailingServiceID) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *MailingServiceServer) GetRecipients(ctx context.Context, req *pb.GetRecipientsRequest) (*pb.MailRecipients, error) {
	var recipients pb.MailRecipients
	return &recipients, nil
}

func (m *MailingServiceServer) GetRecipientsByGroup(ctx context.Context, req *pb.GetRecipientsByGroupRequest) (*pb.MailRecipients, error) {
	var recipients pb.MailRecipients
	return &recipients, nil
}

func (m *MailingServiceServer) SearchRecipients(ctx context.Context, req *pb.SearchRequest) (*pb.MailRecipients, error) {
	var recipients pb.MailRecipients
	return &recipients, nil
}

func (m *MailingServiceServer) CountRecipients(ctx context.Context, req *pb.MailingServiceGroup) (*pb.Count, error) {
	var count pb.Count
	return &count, nil
}
