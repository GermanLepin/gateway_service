package changelog

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddUsers, downAddUsers)
}

func upAddUsers(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table service.users (
			id uuid not null unique,
			first_name varchar(100) not null,
			last_name varchar(100) not null,
			password varchar(100) not null,
			email varchar(100) not null primary key,
			phone bigint not null,
			user_type varchar(100) not null
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func downAddUsers(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop table if exists service.users;
	`)
	if err != nil {
		return err
	}

	return nil
}
