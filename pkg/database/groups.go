package database

import (
	"strings"
	"time"

	pb "github.com/hyperxpizza/mailing-service/pkg/grpc"
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
