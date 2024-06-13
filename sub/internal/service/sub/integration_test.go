//go:build integration

package sub

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/pkg/tests"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/validator"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
)

var mySQLContainer *mysql.MySQLContainer

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	var err error
	mySQLContainer, err = tests.NewDB(ctx)
	if err != nil {
		panic(err)
	}

	defer func() {
		tctx, tcancel := context.WithTimeout(context.Background(), time.Second*10)
		defer tcancel()
		if err = mySQLContainer.Terminate(tctx); err != nil {
			slog.Info("Failed to terminate test container", slog.Any("err", err))
		}
	}()

	migrations, err := filepath.Abs(filepath.Join("../../../", "migrations"))
	if err != nil {
		panic(fmt.Sprintf("failed to find migrations dir %s: %v", migrations, err))
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	conn, err := mySQLContainer.ConnectionString(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to get conn string from db container: %v", err))
	}

	if err := tests.MigrateDB(migrations, conn); err != nil {
		panic(err)
	}

	exitCode := m.Run()

	if exitCode != 0 {
		panic("test run returned negative exit code")
	}
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
			db, err := db.NewConn(tests.MustGetDSN(mySQLContainer))
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
