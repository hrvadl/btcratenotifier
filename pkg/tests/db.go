package tests

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	mysqlc "github.com/testcontainers/testcontainers-go/modules/mysql"
)

const (
	dbName   = "test_converter"
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = "password"
)

const dsnTimeout = 5 * time.Second

func MustGetDSN(c *mysqlc.MySQLContainer) string {
	ctx, cancel := context.WithTimeout(context.Background(), dsnTimeout)
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

func NewDB(ctx context.Context) (*mysqlc.MySQLContainer, error) {
	container, err := mysqlc.RunContainer(ctx,
		testcontainers.WithImage("mysql:latest"),
		mysqlc.WithDatabase(dbName),
		mysqlc.WithUsername(user),
		mysqlc.WithPassword(password),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	return container, nil
}

func MigrateDB(path, dsn string) error {
	migrations, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to find migrations dir %s: %w", migrations, err)
	}

	m, err := migrate.New(toFileScheme(migrations), toMySQLScheme(dsn))
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return m.Up()
}
