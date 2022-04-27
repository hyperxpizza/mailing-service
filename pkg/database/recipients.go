package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hyperxpizza/mailing-service/pkg/customErrors"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var bgContext = context.Background()

const (
	getRecipientsBase          = "select id, email, usersServiceID, created, updated, confirmed from mailRecipients"
	getRecipientsBaseGroupName = "select r.id, r.email, r.usersServiceID, r.created, r.updated, r.confirmed from mailRecipients as r join recipientGroupMap as m on r.id = m.recipientID join mailGroups as g on m.groupID = g.id"
)

func (db *Database) InsertMailRecipient(req *pb.NewMailRecipient) (int64, error) {
	var id int64

	tx, err := db.BeginTx(bgContext, nil)
	if err != nil {
		return 0, err
	}

	//check if group exists
	var groupID int
	err = tx.QueryRow(`select id from mailGroups where groupName=$1`, strings.ToUpper(req.GroupName)).Scan(&groupID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return 0, customErrors.NewGroupNotFoundError(req.GroupName)
		} else {
			return 0, err
		}
	}

	//insert mailRecipient
	stmt, err := tx.Prepare(`insert into mailRecipients (id, email, usersServiceID, created, updated, confirmed) values (default, $1, $2, $3, $4, $5) returning id`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	var usersServiceID *int64
	if req.UsersServiceID == 0 {
		usersServiceID = nil
	} else {
		usersServiceID = &req.UsersServiceID
	}

	err = stmt.QueryRow(req.Email, usersServiceID, time.Now(), time.Now(), req.Confirmed).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	//insert into recipientGroupMap
	stmt2, err := tx.Prepare(`insert into recipientGroupMap(id, groupID, recipientID) values(default, $1, $2)`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	defer stmt2.Close()

	_, err = stmt2.Exec(groupID, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, nil
}

func (db *Database) ConfirmRecipient(email string) error {
	stmt, err := db.Prepare(`update mailRecipients set confirmed=true where email=$1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteRecipientByID(id int64) error {

	tx, err := db.BeginTx(bgContext, nil)
	if err != nil {
		return err
	}

	stmt2, err := tx.Prepare(`delete from recipientGroupMap where recipientID=$1`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt2.Close()

	_, err = stmt2.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare(`delete from mailRecipients where id=$1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return customErrors.NewIDNotFoundError(id)
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (db *Database) GetRecipientByID(id int64) (*pb.MailRecipient, error) {
	var recipient pb.MailRecipient
	var groups []*pb.MailGroup

	tx, err := db.BeginTx(bgContext, nil)
	if err != nil {
		return nil, err
	}

	var usID sql.NullInt64
	var confirmed sql.NullBool
	var created time.Time
	var updated time.Time

	err = tx.QueryRow(`select id, email, usersServiceID, created, updated, confirmed from mailRecipients where id=$1`, id).Scan(
		&recipient.Id,
		&recipient.Email,
		&usID,
		&created,
		&updated,
		&confirmed,
	)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewIDNotFoundError(id)
		}

		return nil, err
	}

	defer func() {
		recipient.Created = timestamppb.New(created)
		recipient.Updated = timestamppb.New(updated)
		if confirmed.Valid {
			recipient.Confirmed = confirmed.Bool
		}
		recipient.Groups = groups
		if usID.Valid {
			recipient.UsersServiceID = usID.Int64
		}
	}()

	rows, err := tx.Query(`select m.groupID, g.groupName, g.created, g.updated from recipientGroupMap as m join mailGroups as g on g.id=m.groupID where recipientID=$1`, recipient.Id)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewIDNotFoundError(id)
		}

		return nil, err
	}

	for rows.Next() {
		var mailGroup pb.MailGroup
		var gCreated time.Time
		var gUpdated time.Time
		err := rows.Scan(
			&mailGroup.Id,
			&mailGroup.Name,
			&gCreated,
			&gUpdated,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		mailGroup.Created = timestamppb.New(gCreated)
		mailGroup.Updated = timestamppb.New(gUpdated)
		groups = append(groups, &mailGroup)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &recipient, nil
}

func (db *Database) CheckIfRecipientConfirmed(usID int64) (*bool, error) {
	var confirmed bool

	var c sql.NullBool
	err := db.QueryRow(`select confirmed from mailRecipients where usersServiceID=$1`).Scan(&c)
	if err != nil {
		return nil, err
	}

	if c.Valid {
		confirmed = c.Bool
	}

	return &confirmed, nil
}

func (db *Database) CountRecipients(groupName string) (int64, error) {
	var count int64

	base := "select COUNT(*) from mailRecipients"
	extra := "as r join recipientGroupMap as m on r.id = m.recipientID join mailGroups as g on m.groupID = g.id where g.groupName = $1"
	var query string
	if groupName != "" {
		query = fmt.Sprintf("%s %s", base, extra)
		err := db.QueryRow(query, strings.ToUpper(groupName)).Scan(&count)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, customErrors.NewGroupNotFoundError(groupName)
			}

			return 0, err
		}

	} else {
		query = base
		err := db.QueryRow(query).Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (db *Database) GetRecipients(req *pb.GetRecipientsRequest) ([]*pb.MailRecipient, error) {
	query, allowedVars := buildGetRecipientsQuery(req.Order.String(), req.Pagination.Limit, req.Pagination.Offset)
	fmt.Println(query)
	rows, err := db.Query(query, allowedVars...)
	if err != nil {
		return nil, err
	}

	recipients, err := scanRecipients(rows)
	if err != nil {
		return nil, err
	}

	return recipients, nil
}

func (db *Database) GetRecipientsByGroup(req *pb.GetRecipientsByGroupRequest) ([]*pb.MailRecipient, error) {
	query, allowedVars := buildGetRecipientsWhereGroupQuery(req.Order.String(), req.Group, req.Pagination.Limit, req.Pagination.Offset)
	rows, err := db.Query(query, allowedVars...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewGroupNotFoundError(req.Group)
		}
		return nil, err
	}

	recipients, err := scanRecipients(rows)
	if err != nil {
		return nil, err
	}

	return recipients, nil
}

func (db *Database) SearchRecipients(req *pb.SearchRequest) ([]*pb.MailRecipient, error) {
	query, allowedVars := buildSearchQuery(req.Query, req.Order.String(), req.Pagination.Limit, req.Pagination.Offset)
	rows, err := db.Query(query, allowedVars...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewNoResultsError(req.Query)
		}

		return nil, err
	}

	recipients, err := scanRecipients(rows)
	if err != nil {
		return nil, err
	}

	return recipients, nil
}
