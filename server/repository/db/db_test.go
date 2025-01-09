package db_test

import (
	"context"
	"database/sql"
	_ "embed"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/mysql"
	dbrepo "github.com/traPtitech/piscon-portal-v2/server/repository/db"
	"github.com/traPtitech/piscon-portal-v2/server/utils/random"
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

func setupRepository(t *testing.T) *dbrepo.Repository {
	t.Helper()

	ctx := context.Background()

	dbName := "test_" + random.String(10)
	if err := createDatabase(dbName); err != nil {
		t.Fatal(err)
	}
	connection := mysqlContainer.MustConnectionString(ctx, "parseTime=true")
	db, err := sql.Open("mysql", connection)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("USE " + dbName); err != nil {
		t.Fatal(err)
	}

	return dbrepo.NewRepository(db)
}

func createDatabase(name string) error {
	connection := mysqlContainer.MustConnectionString(context.Background(), "multiStatements=true")
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	schema, err := os.ReadFile("testdata/schema.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec("USE " + name); err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	return err
}

// retry retries f until it returns nil or max retries are reached. panic if max retries are reached.
func retry(max int, f func() error) {
	var err error
	for i := 0; i < max; i++ {
		err = f()
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
	panic(err)
}
