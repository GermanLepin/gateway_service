package changelog

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upInit, downInit)
}

func upInit(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create schema service;
	`)
	if err != nil {
		return err
	}

	return nil
}

func downInit(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop schema service;
	`)
	if err != nil {
		return err
	}

	return nil
}
