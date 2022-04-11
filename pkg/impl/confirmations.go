package impl

import (
	"context"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MailingServiceServer) SendConfirmationEmail(ctx context.Context, req *pb.MailingServiceEmail) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *MailingServiceServer) ConfirmRecipient(ctx context.Context, req *pb.RecipientConfirmation) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
