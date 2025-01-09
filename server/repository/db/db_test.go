package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/mysql"
	dbrepo "github.com/traPtitech/piscon-portal-v2/server/repository/db"
)

func setupRepository(t *testing.T) *dbrepo.Repository {
	t.Helper()

	ctx := context.Background()

	container, err := mysql.Run(ctx,
		"mysql:8",
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
		mysql.WithDatabase("test"),
		mysql.WithScripts("testdata/schema.sql"),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	connection := container.MustConnectionString(ctx, "parseTime=true")
	db, err := sql.Open("mysql", connection)
	if err != nil {
		t.Fatal(err)
	}

	for range 30 {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return dbrepo.NewRepository(db)
}
