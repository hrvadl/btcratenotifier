package subscriber

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Subscriber struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

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
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *Repo) Get(ctx context.Context) ([]Subscriber, error) {
	subscribers := []Subscriber{}
	if err := r.db.SelectContext(ctx, &subscribers, "SELECT * FROM subscribers"); err != nil {
		return nil, err
	}

	return subscribers, nil
}