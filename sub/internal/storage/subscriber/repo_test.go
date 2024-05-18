package subscriber

import (
	"testing"

	"github.com/jmoiron/sqlx"
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
			if got := NewRepo(tt.args.db); got == nil {
				t.Errorf("NewRepo() = %v, want not nil", got)
			}
		})
	}
}
