package mailing

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type List struct {
	ID       int64  `db:"id"`
	LastSent string `db:"last_sent"`
}

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) UpdateLastSentDate(ctx context.Context, d time.Time) error {
	_, err := r.db.ExecContext(ctx, "REPLACE INTO mailing_list (id, last_sent) VALUES (1, $1)", d)
	return err
}
