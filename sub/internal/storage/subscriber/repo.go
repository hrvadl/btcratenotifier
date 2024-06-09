package subscriber

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
)

// Repo is a thin abstraction to not do sqlx queries
// directly in the services. Therefore specific underlying DB could
// be more easily changed in the future.
type Repo struct {
	db *sqlx.DB
}

// NewRepo constructs repo with provided sqlx DB connection.
// NOTE: it expectes db connection to be connection MySQL.
func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// Save method saves subscriber to the repo and then returns
// newly created ID. Could return an error if email is not valid, or such email
// already exists.
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

// Get method gets all subscribers from the DB.
func (r *Repo) Get(ctx context.Context) ([]Subscriber, error) {
	subscribers := []Subscriber{}
	if err := r.db.SelectContext(ctx, &subscribers, "SELECT * FROM subscribers"); err != nil {
		return nil, err
	}

	return subscribers, nil
}
