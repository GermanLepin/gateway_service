package changelog

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upSessions, downSessions)
}

func upSessions(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table service.sessions (
			id uuid primary key, 
			user_id uuid not null,
			first_name varchar(100) not null,
			last_name varchar(100) not null,
			user_type varchar(100) not null,
			user_agent varchar(100) not null,
			client_ip varchar(100) not null,
			is_blocked boolean not null default false,
			expires_at timestamptz not null,
			created_at timestamptz not null default (now())
		);

		alter table service.sessions add foreign key ("user_id") references service.users ("id");		
	`)
	if err != nil {
		return err
	}

	return nil
}

func downSessions(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop table if exists service.sessions;
	`)
	if err != nil {
		return err
	}

	return nil
}
