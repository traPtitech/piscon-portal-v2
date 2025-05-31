package db_test

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"
	"time"

	"github.com/stephenafamo/bob"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	dbrepo "github.com/traPtitech/piscon-portal-v2/server/repository/db"
)

var mysqlContainer *mysql.MySQLContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := mysql.Run(ctx,
		"mysql:8",
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
	)
	if err != nil {
		panic(err)
	}
	defer container.Terminate(ctx) //nolint errcheck

	connection := container.MustConnectionString(ctx)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}
	retry(30, func() error { return db.Ping() }) // ensure db is ready
	db.Close()

	mysqlContainer = container

	m.Run()
}

func setupRepository(t *testing.T) (*dbrepo.Repository, bob.Executor) {
	t.Helper()

	ctx := context.Background()

	dbName := randomDBName()
	if err := createDatabase(dbName); err != nil {
		t.Fatal(err)
	}
	connection := mysqlContainer.MustConnectionString(ctx, "parseTime=true", "loc=Local")
	db, err := sql.Open("mysql", connection)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("USE " + dbName); err != nil {
		t.Fatal(err)
	}

	return dbrepo.NewRepository(db), bob.NewDB(db)
}

func createDatabase(name string) error {
	connection := mysqlContainer.MustConnectionString(context.Background(), "multiStatements=true")
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return fmt.Errorf("create database %s: %w", name, err)
	}
	schema, err := os.ReadFile("testdata/schema.sql")
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	if _, err := db.Exec("USE " + name); err != nil {
		return fmt.Errorf("use database %s: %w", name, err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}
	return nil
}

// retry retries f until it returns nil or n retries are reached. panic if n retries are reached.
func retry(n int, f func() error) {
	var err error
	for i := 0; i < n; i++ {
		err = f()
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
	panic(err)
}

func randomDBName() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var buf [10]byte
	for i := range buf {
		buf[i] = charset[rand.IntN(len(charset))]
	}
	return "test_" + string(buf[:])
}
