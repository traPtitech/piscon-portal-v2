package main

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/piscon-portal-v2/server/handler"
)

func main() {
	dbConfig := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_NAME"),
		ParseTime: true,
	}
	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())

	config := handler.Config{
		RootURL:       os.Getenv("ROOT_URL"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		Oauth2: handler.Oauth2Config{
			Issuer:       "https://q.trap.jp",
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			AuthURL:      "https://q.trap.jp/api/v3/oauth2/authorize",
			TokenURL:     "https://q.trap.jp/api/v3/oauth2/token",
		},
	}
	handler, err := handler.New(db, config)
	if err != nil {
		panic(err)
	}
	handler.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
