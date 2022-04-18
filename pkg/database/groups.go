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

	return mg, nil
}
