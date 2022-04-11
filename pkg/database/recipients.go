package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
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

func (db *Database) DeleteRecipientByEmail(email string) error {

	tx, err := db.BeginTx(bgContext, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`delete from mailRecipients where email=$1`)
	if err != nil {
		return err
	}

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
