//go:build integration

package subscriber

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

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
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

func TestSave(t *testing.T) {
	type args struct {
		ctx context.Context
		sub Subscriber
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should save subscriber correctly",
			args: args{
				ctx: context.Background(),
				sub: Subscriber{Email: "test@mail.com"},
			},
			wantErr: false,
		},
		{
			name: "Should not save subscriber twice",
			args: args{
				ctx: context.Background(),
				sub: Subscriber{Email: "test@mail.com"},
			},
			wantErr: true,
		},
		{
			name: "Should not get subscribers correctly when it takes too long",
			args: args{
				ctx: newImmediateCtx(),
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, err := db.NewConn(tests.MustGetDSN(mySQLContainer))
			require.NoError(t, err, "Failed to connect to test DB")

			r := NewRepo(db)
			id, err := r.Save(tt.args.ctx, tt.args.sub)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, id)
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should get subscribers correctly",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Should not get subscribers correctly when it takes too long",
			args: args{
				ctx: newImmediateCtx(),
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, err := db.NewConn(tests.MustGetDSN(mySQLContainer))
			require.NoError(t, err, "Failed to connect to test DB")

			r := NewRepo(db)
			want := seed(t, r, 30)

			got, err := r.Get(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Subset(t, mapSubsToMails(got), mapSubsToMails(want))
		})
	}
}

func seed(t *testing.T, repo *Repo, amount int) []Subscriber {
	t.Helper()

	subs := make([]Subscriber, 0, amount)
	for range amount {
		mail := fmt.Sprintf("mail%v@mail.com", time.Now().Nanosecond())
		sub := Subscriber{Email: mail}
		subs = append(subs, sub)
		id, err := repo.Save(context.Background(), sub)
		require.NoError(t, err)
		require.NotZero(t, id)
	}

	return subs
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}

func mapSubsToMails(s []Subscriber) []string {
	mails := make([]string, 0, len(s))
	for i := range s {
		mails = append(mails, s[i].Email)
	}
	return mails
}
