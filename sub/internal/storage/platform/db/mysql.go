package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Must(db *sqlx.DB, err error) *sqlx.DB {
	if err != nil {
		panic(err)
	}
	return db
}

func NewConn(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", dsn)
	return db, err
}
