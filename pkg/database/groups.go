package database

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/hyperxpizza/mailing-service/pkg/customErrors"
	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (db *Database) InsertGroup(name string) (int64, error) {
	var id int64

	stmt, err := db.Prepare(`insert into mailGroups(id, groupName, created, updated) values(default, $1, $2, $3) returning id`)
	if err != nil {
		return 0, err
	}

	err = stmt.QueryRow(strings.ToUpper(name), time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) GetGroups() ([]*pb.MailGroup, error) {
	var mg []*pb.MailGroup

	rows, err := db.Query(`select * from mailGroups`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var group pb.MailGroup
		var created time.Time
		var updated time.Time
		err := rows.Scan(
			&group.Id,
			&group.Name,
			&created,
			&updated,
		)
		if err != nil {
			return nil, err
		}

		mg = append(mg, &group)
	}

	return mg, nil
}

func (db *Database) GetGroup(id int64) (*pb.MailGroup, error) {
	var group pb.MailGroup
	var created time.Time
	var updated time.Time
	err := db.QueryRow(`select * from mailGroups where id=$1`).Scan(
		&group.Id,
		&group.Name,
		&created,
		&updated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewIDNotFoundError(id)
		}

		return nil, err
	}

	group.Created = timestamppb.New(created)
	group.Updated = timestamppb.New(updated)

	return &group, nil
}

func (db *Database) DeleteGroup(id int64) error {
	_, err := db.Exec(`delete from mailGroups where id=$1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customErrors.NewIDNotFoundError(id)
		}

		return err
	}

	return nil
}

func (db *Database) UpdateGroupName(newName string, id int64) error {
	_, err := db.Exec(`update mailGroups set groupName=$1, updated=$2 where id=$3`, newName, time.Now(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customErrors.NewIDNotFoundError(id)
		}

		return err
	}
	return nil
}
