//go:build integration

package sub

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	mysqlc "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/validator"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
)

var mySQLContainer *mysqlc.MySQLContainer

const (
	dbName   = "test_converter"
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = "password"
)

func mustGetDSN(c *mysqlc.MySQLContainer) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	dsn, err := c.ConnectionString(ctx, "parseTime=true", "tls=skip-verify")
	if err != nil {
		panic(err)
	}
	return dsn
}

func toFileScheme(path string) string {
	return "file://" + path
}

func toMySQLScheme(dsn string) string {
	return "mysql://" + dsn
}

func setup(ctx context.Context) error {
	var err error
	mySQLContainer, err = mysqlc.RunContainer(ctx,
		testcontainers.WithImage("mysql:latest"),
		mysqlc.WithDatabase(dbName),
		mysqlc.WithUsername(user),
		mysqlc.WithPassword(password),
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	migrations, err := filepath.Abs(filepath.Join("../../../", "migrations"))
	if err != nil {
		return fmt.Errorf("failed to find migrations dir %s: %w", migrations, err)
	}

	m, err := migrate.New(toFileScheme(migrations), toMySQLScheme(mustGetDSN(mySQLContainer)))
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return m.Up()
}

func teardown(ctx context.Context) error {
	return mySQLContainer.Terminate(ctx)
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	if err := setup(ctx); err != nil {
		panic(err)
	}

	exitCode := m.Run()

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	if err := teardown(ctx); err != nil {
		panic(err)
	}

	if exitCode != 0 {
		panic("test run returned negative exit code")
	}
}

func TestServiceSend(t *testing.T) {
	type args struct {
		ctx  context.Context
		mail string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should subscribe correctly",
			args: args{
				ctx:  context.TODO(),
				mail: "test@mail.com",
			},
			wantErr: false,
		},
		{
			name: "Should not subscribe when email is incorrect",
			args: args{
				ctx:  context.TODO(),
				mail: "tetmail.com",
			},
			wantErr: true,
		},
		{
			name: "Should not subscribe when email already exists",
			args: args{
				ctx:  context.TODO(),
				mail: "test@mail.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := db.NewConn(mustGetDSN(mySQLContainer))
			require.NoError(t, err, "Failed to connect to db")

			rs := subscriber.NewRepo(db)
			v := validator.NewStdlib()
			s := NewService(rs, v)
			id, err := s.Subscribe(tt.args.ctx, tt.args.mail)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NotZero(t, id)
		})
	}
}
