//go:build integration

package sub

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/validator"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestServiceSend(t *testing.T) {
	type args struct {
		ctx  context.Context
		mail string
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should subscribe correctly",
			args: args{
				ctx:  context.Background(),
				mail: "test@mail.com",
			},
			wantErr: false,
		},
		{
			name: "Should not subscribe correctly when it takes too long",
			args: args{
				ctx:  newImmediateCtx(),
				mail: "test1111@mail.com",
			},
			wantErr: true,
		},
		{
			name: "Should not subscribe when email is empty",
			args: args{
				ctx:  context.Background(),
				mail: "",
			},
			wantErr: true,
		},
		{
			name: "Should not subscribe when email is incorrect",
			args: args{
				ctx:  context.Background(),
				mail: "tetmail.com",
			},
			wantErr: true,
		},
		{
			name: "Should not subscribe when email already exists",
			args: args{
				ctx:  context.Background(),
				mail: "test@mail.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			dsn := os.Getenv("SUB_DSN")
			require.NotZero(t, dsn, "test DSN can not be empty")
			db, err := db.NewConn(dsn)
			require.NoError(t, err, "Failed to connect to db")

			rs := subscriber.NewRepo(db)
			v := validator.NewStdlib()
			s := NewService(rs, v)
			id, err := s.Subscribe(tt.args.ctx, tt.args.mail)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, id)
		})
	}
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}
