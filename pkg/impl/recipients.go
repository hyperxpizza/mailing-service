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

	id.Id = dbID
	return &id, nil
}

func (m *MailingServiceServer) RemoveRecipient(ctx context.Context, req *pb.MailingServiceID) (*emptypb.Empty, error) {

	err := m.db.DeleteRecipientByID(req.Id)
	if err != nil {
		var idNotFoundErr *customErrors.NotFoundError
		if errors.As(err, &idNotFoundErr) {
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

func (m *MailingServiceServer) GetRecipient(ctx context.Context, req *pb.MailingServiceID) (*pb.MailRecipient, error) {
	recipient, err := m.db.GetRecipientByID(req.Id)
	if err != nil {
		var idNotFoundErr *customErrors.NotFoundError
		if errors.As(err, &idNotFoundErr) {
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

	return recipient, nil
}

func (m *MailingServiceServer) GetRecipients(ctx context.Context, req *pb.GetRecipientsRequest) (*pb.MailRecipients, error) {
	var recipients pb.MailRecipients

	rec, err := m.db.GetRecipients(req)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}
	recipients.MailRecipients = rec

	return &recipients, nil
}

func (m *MailingServiceServer) GetRecipientsByGroup(ctx context.Context, req *pb.GetRecipientsByGroupRequest) (*pb.MailRecipients, error) {
	var recipients pb.MailRecipients
	rec, err := m.db.GetRecipientsByGroup(req)
	if err != nil {
		var gNotFoundErr *customErrors.NotFoundError
		if errors.As(err, &gNotFoundErr) {
			return nil, err
		}

		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}
	recipients.MailRecipients = rec
	return &recipients, nil
}

func (m *MailingServiceServer) SearchRecipients(ctx context.Context, req *pb.SearchRequest) (*pb.MailRecipients, error) {
	var recipients pb.MailRecipients
	return &recipients, nil
}

func (m *MailingServiceServer) CountRecipients(ctx context.Context, req *pb.MailingServiceGroup) (*pb.Count, error) {
	var count pb.Count

	c, err := m.db.CountRecipients(req.Group)
	if err != nil {
		var gNotFoundErr *customErrors.NotFoundError
		if errors.As(err, &gNotFoundErr) {
			return nil, err
		}

		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	count.Num = c

	return &count, nil
}
