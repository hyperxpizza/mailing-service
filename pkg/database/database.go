package database

import (
	"database/sql"
	"fmt"

	"github.com/hyperxpizza/mailing-service/pkg/config"
	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func Connect(cfg *config.Config) (*Database, error) {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{database}, nil
}

func (db *Database) InsertRecipient()                   {}
func (db *Database) GetRecipients()                     {}
func (db *Database) GetRecipientByUserServiceID(id int) {}
