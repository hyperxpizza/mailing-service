package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
)

const (
	GroupNotFoundError = "group with name: %s does not exist"
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
			return 0, fmt.Errorf(GroupNotFoundError, req.GroupName)
		} else {
			return 0, err
		}
	}

	//insert mailRecipient
	stmt, err := tx.Prepare(`insert into mailRecipients (id, email, usersServiceID, created, updated) values (default, $1, $2, $3, $4) returning id`)
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

	err = stmt.QueryRow(req.Email, usersServiceID, time.Now(), time.Now()).Scan(id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, nil
}
