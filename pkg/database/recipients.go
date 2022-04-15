package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var bgContext = context.Background()

func (db *Database) InsertMailRecipient(req *pb.NewMailRecipient) (int64, error) {
	var id int64

	tx, err := db.BeginTx(bgContext, nil)
	if err != nil {
		return 0, err
	}

	//check if group exists
	var groupID int
	err = tx.QueryRow(`select id from mailGroups where groupName=$1`, req.GroupName).Scan(&groupID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return 0, NewGroupNotFoundError(req.GroupName)
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
		*usersServiceID = req.UsersServiceID
	}

	err = stmt.QueryRow(req.Email, usersServiceID, time.Now(), time.Now(), req.Confirmed).Scan(id)
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

	_, err = stmt.Exec(groupID, id)
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

	_, err = stmt.Exec(email)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteRecipientByEmail(email string) error {

	tx, err := db.BeginTx(bgContext, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`delete from mailRecipients where email=$1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return NewRecipientNotFoundError(email)
		}

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
		return nil, err
	}

	rows, err := tx.Query(`select m.groupID, g.groupName, g.created, g.updated from recipientGroupMap as m join mailGroups as g on g.id=m.groupID where recipientID=$1`, recipient.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for rows.Next() {
		var mailGroup pb.MailGroup
		var gCreated time.Time
		var gUpdated time.Time
		err := rows.Scan(
			&mailGroup.Id,
			&mailGroup.Name,
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

	recipient.Created = timestamppb.New(created)
	recipient.Updated = timestamppb.New(updated)
	if confirmed.Valid {
		recipient.Confirmed = confirmed.Bool
	}
	recipient.Groups = groups

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

	} else {
		query = base
	}

	return count, nil
}
