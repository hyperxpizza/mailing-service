package database

func (db *Database) InsertGroup(name string) (int64, error) {
	var id int64

	stmt, err := db.Prepare(`insert into mailGroups(id, groupName, created, updated) values(default, $1, $2, $3) returning id`)
	if err != nil {
		return 0, err
	}

	err = stmt.QueryRow(name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
