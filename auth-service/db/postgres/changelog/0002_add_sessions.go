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
			is_blocked boolean not null default false,
			refresh_token varchar not null,
			expires_at timestamptz not null,
			created_at timestamptz not null default (now())
		);
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
