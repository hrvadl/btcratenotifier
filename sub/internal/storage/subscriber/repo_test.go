//go:build !integration

package subscriber

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestNewRepo(t *testing.T) {
	t.Parallel()
	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should create repo with correct db conn",
			args: args{
				db: &sqlx.DB{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewRepo(tt.args.db)
			require.NotNil(t, got)
		})
	}
}
