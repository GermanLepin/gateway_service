package changelog

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upPaymentInformation, downPaymentInformation)
}

func upPaymentInformation(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table service.payment_information (
			operation_id uuid not null primary key,
			user_id uuid not null,
			amount numeric(10,2) not null,
			card_number numeric(10,2) not null,
			card_holder_name varchar(100) not null,
			cvv numeric(10,2) not null,
			payment_status varchar(100) not null
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func downPaymentInformation(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop table if exists service.payment_information;
	`)
	if err != nil {
		return err
	}

	return nil
}
