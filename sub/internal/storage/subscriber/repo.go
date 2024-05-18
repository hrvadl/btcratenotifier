package subscriber

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/hrvadl/btcratenotifier/sub/internal/storage/platform/db"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Save(ctx context.Context, s Subscriber) (int64, error) {
	res, err := r.db.ExecContext(ctx, "INSERT INTO subscribers (email) VALUES (?)", s.Email)
	if err == nil {
		return res.LastInsertId()
	}

	var mySQLErr *mysql.MySQLError
	if errors.As(err, &mySQLErr) && mySQLErr.Number == db.AlreadyExistsErrCode {
		return 0, ErrAlreadyExists
	}

	return 0, err
}

func (r *Repo) Get(ctx context.Context) ([]Subscriber, error) {
	subscribers := []Subscriber{}
	if err := r.db.SelectContext(ctx, &subscribers, "SELECT * FROM subscribers"); err != nil {
		return nil, err
	}

	return subscribers, nil
}
