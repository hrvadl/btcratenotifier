package date

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type LatestSent struct {
	ID   int64     `db:"id"`
	Date time.Time `db:"last_sent"`
}

type Repo struct {
	db *sqlx.DB
}

func NewLastSentRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) UpdateLatestSent(ctx context.Context, d time.Time) error {
	_, err := r.db.ExecContext(ctx, "REPLACE INTO mailing_list (id, last_sent) VALUES (1, $1)", d)
	return err
}

func (r *Repo) GetLatestSent(ctx context.Context) (*LatestSent, error) {
	l := new(LatestSent)
	if err := r.db.SelectContext(ctx, l, "SELECT * FROM mailing_list where id = 1"); err != nil {
		return nil, err
	}
	return l, nil
}
