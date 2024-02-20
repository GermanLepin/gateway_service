package connection

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "gateway-service/db/postgres/changelog"
	"gateway-service/internal/application/helper/logging"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"

	"github.com/pressly/goose"
)

var (
	driver = "pgx"
)

func StartDB() *sql.DB {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s  database=%s sslmode=disable timezone=UTC connect_timeout=5",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DATABASE"),
	)

	conn := connectToDB(ctx, dsn)
	if conn == nil {
		logger.Panic("cannot connect to Postgres")
	}

	if err := goose.Up(conn, "/var"); err != nil {
		logger.Panic("cannot run the migrations")
	}

	// if smth goes wrong we always can run down Migrations goose.Down()
	// if err := goose.Down(conn, "/var"); err != nil {
	// 	log.Panic("cannot run the migrations")
	// }

	return conn
}

func connectToDB(ctx context.Context, dsn string) *sql.DB {
	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	var counts int64
	for {
		connection, err := openDB(ctx, dsn)
		if err != nil {
			log.Println("postgres is not ready yet")
			counts++
		} else {
			log.Println("connected to Postgres")
			return connection
		}

		if counts > 10 {
			logger.Error("cannot connect to database", zap.Error(err))
			return nil
		}

		log.Println("backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(ctx context.Context, dsn string) (*sql.DB, error) {
	logger := logging.LoggerFromContext(ctx)

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		logger.Error("connection open is failed", zap.Error(err))
		return nil, err
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err = conn.Ping(); err != nil {
		logger.Error("ddatabase ping is failed", zap.Error(err))
		return nil, err
	}

	return conn, nil
}
