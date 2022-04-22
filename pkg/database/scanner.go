package database

import (
	"database/sql"
	"time"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func scanRecipients(rows *sql.Rows) ([]*pb.MailRecipient, error) {
	var recipients []*pb.MailRecipient
	for rows.Next() {
		var recipient pb.MailRecipient
		var created time.Time
		var updated time.Time
		var usID *sql.NullInt64
		var confirmed *sql.NullBool
		err := rows.Scan(
			&recipient.Id,
			&recipient.Email,
			&usID,
			&created,
			&updated,
			&confirmed,
		)
		if err != nil {
			return nil, err
		}

		if usID.Valid {
			recipient.UsersServiceID = usID.Int64
		}

		if confirmed.Valid {
			recipient.Confirmed = confirmed.Bool
		}

		recipient.Created = timestamppb.New(created)
		recipient.Updated = timestamppb.New(updated)

		recipients = append(recipients, &recipient)

	}

	return recipients, nil
}
