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

		create table service.payment_status (
			operation_id uuid not null primary key,
			status varchar(100) not null
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func downInit(tx *sql.Tx) error {
	_, err := tx.Exec(`
		delete table service.paymen_status;

		drop schema service;
	`)
	if err != nil {
		return err
	}

	return nil
}
