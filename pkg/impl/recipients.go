package impl

import (
	"context"
	"errors"

	"github.com/hyperxpizza/mailing-service/pkg/database"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"github.com/hyperxpizza/mailing-service/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MailingServiceServer) AddRecipient(ctx context.Context, req *pb.NewMailRecipient) (*pb.MailingServiceID, error) {
	var id pb.MailingServiceID

	m.logger.Debugf("adding new recipient with email: %s", req.Email)
	err := utils.ValidateEmail(req.Email)
	if err != nil {
		m.logger.Debugf("email: %s not valid", req.Email)
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	dbID, err := m.db.InsertMailRecipient(req)
	if err != nil {
		var gErr *database.NotFoundError
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

	id.Id = dbID
	return &id, nil
}

func (m *MailingServiceServer) RemoveRecipiet(ctx context.Context, req *pb.MailingServiceID) (*emptypb.Empty, error) {

	return &emptypb.Empty{}, nil
}

func (m *MailingServiceServer) GetRecipient(ctx context.Context, req *pb.MailingServiceID) (*pb.MailRecipient, error) {
	recipient, err := m.db.GetRecipientByID(req.Id)
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
