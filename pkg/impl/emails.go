package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"github.com/hyperxpizza/mailing-service/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	confirmationTokenKey   = "confirmation-token-%s"
	tokenNotFoundError     = "token was not found in the database"
	tokenMismatchError     = "tokens are not matching"
	recipientNotFoundError = "recipient was not found in the database"
)

func (m *MailingServiceServer) SendConfirmationEmail(ctx context.Context, req *pb.MailingServiceEmail) (*emptypb.Empty, error) {
	m.logger.Infof("sending confirmation to")

	key := fmt.Sprintf(confirmationTokenKey, req.Email)

	//check if token already exists
	token, err := m.rdc.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return &emptypb.Empty{}, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	if token != "" {
		_, err := m.rdc.Del(ctx, key).Result()
		if err != nil {
			return &emptypb.Empty{}, status.Error(
				codes.Internal,
				err.Error(),
			)
		}
	}

	confToken := utils.GenerateOneTimeToken()
	exp := time.Duration(m.cfg.MailingService.ConfirmationTokenExpirationMinutes) * time.Minute
	go m.rdc.Set(ctx, key, confToken, exp)

	//send confirmation email

	return &emptypb.Empty{}, nil
}

func (m *MailingServiceServer) ConfirmRecipient(ctx context.Context, req *pb.RecipientConfirmation) (*emptypb.Empty, error) {
	key := fmt.Sprintf(confirmationTokenKey, req.Email)

	confToken, err := m.rdc.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &emptypb.Empty{}, status.Error(
				codes.NotFound,
				tokenNotFoundError,
			)
		}

		return &emptypb.Empty{}, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	if confToken != req.Token {
		return &emptypb.Empty{}, status.Error(
			codes.InvalidArgument,
			tokenMismatchError,
		)
	}

	//update database record
	err = m.db.ConfirmRecipient(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &emptypb.Empty{}, status.Error(
				codes.NotFound,
				err.Error(),
			)
		}

		return &emptypb.Empty{}, status.Error(
			codes.Internal,
			err.Error(),
		)

	}

	return &emptypb.Empty{}, nil
}

func (m *MailingServiceServer) CheckIfRecipientIsConfirmed(ctx context.Context, req *pb.CheckIfConfirmedRequest) (*pb.Cofirmed, error) {
	var confirmed pb.Cofirmed

	conf, err := m.db.CheckIfRecipientConfirmed(req.UsersServiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(
				codes.NotFound,
				recipientNotFoundError,
			)
		}

		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	confirmed.Confirmed = *conf

	return &confirmed, nil
}
