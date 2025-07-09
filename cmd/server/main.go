package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/piscon-portal-v2/server/handler"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	dbrepo "github.com/traPtitech/piscon-portal-v2/server/repository/db"
	"github.com/traPtitech/piscon-portal-v2/server/server"
	"github.com/traPtitech/piscon-portal-v2/server/services/instance"
	"github.com/traPtitech/piscon-portal-v2/server/services/instance/aws"
	"github.com/traPtitech/piscon-portal-v2/server/services/instance/fake"
	"github.com/traPtitech/piscon-portal-v2/server/services/oauth2"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	handlerConfig := handler.Config{
		RootURL:       os.Getenv("ROOT_URL"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		Oauth2: oauth2.Config{
			Issuer:       "https://q.trap.jp",
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			AuthURL:      "https://q.trap.jp/api/v3/oauth2/authorize",
			TokenURL:     "https://q.trap.jp/api/v3/oauth2/token",
		},
	}

	useCaseConfig, err := provideUseCaseConfig()
	if err != nil {
		log.Fatal("failed to create use case config:", err)
	}
	manager, err := provideInstanceManager()
	if err != nil {
		log.Fatal("failed to create instance manager:", err)
	}
	repo, err := provideRepository()
	if err != nil {
		log.Fatal("failed to create repository:", err)
	}

	useCase := usecase.New(useCaseConfig, repo, manager)
	handler, err := handler.New(useCase, repo, handlerConfig)
	if err != nil {
		panic(err)
	}
	handler.SetupRoutes(e)

	benchService := server.NewBenchmarkService(useCase)
	go func() {
		log.Fatal("error in bench service.", server.Start(50051, benchService))
	}()

	e.Logger.Fatal(e.Start(":8080"))
}

func provideRepository() (repository.Repository, error) {
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
		return nil, fmt.Errorf("open database: %w", err)
	}
	return dbrepo.NewRepository(db), nil
}

func provideUseCaseConfig() (usecase.Config, error) {
	instanceLimit := 3 // default instance limit
	if str := os.Getenv("INSTANCE_LIMIT"); str != "" {
		var err error
		instanceLimit, err = strconv.Atoi(str)
		if err != nil {
			return usecase.Config{}, fmt.Errorf("invalid INSTANCE_LIMIT: %w", err)
		}
	}
	return usecase.Config{InstanceLimit: instanceLimit}, nil
}

func provideInstanceManager() (instance.Manager, error) {
	manager := os.Getenv("INSTANCE_MANAGER")
	if manager == "aws" {
		return provideAWSInstanceManager()
	}
	if manager == "" || manager == "mock" {
		return provideMockInstanceManager()
	}
	return nil, fmt.Errorf("unknown INSTANCE_MANAGER: %s", manager)
}

func provideMockInstanceManager() (instance.Manager, error) {
	root, err := os.OpenRoot("/app/.dev/instance")
	if err != nil {
		return nil, fmt.Errorf("open root directory: %w", err)
	}
	return fake.NewManager(root)
}

func provideAWSInstanceManager() (instance.Manager, error) {
	awsConfig := aws.Config{
		ImageID:         os.Getenv("AWS_IMAGE_ID"),
		InstanceType:    os.Getenv("AWS_INSTANCE_TYPE"),
		Region:          os.Getenv("AWS_REGION"),
		AccessKey:       os.Getenv("AWS_ACCESS_KEY"),
		SecretKey:       os.Getenv("AWS_SECRET_KEY"),
		SubnetID:        os.Getenv("AWS_SUBNET_ID"),
		SecurityGroupID: os.Getenv("AWS_SECURITY_GROUP_ID"),
		KeyPairName:     os.Getenv("AWS_KEY_PAIR_NAME"),
	}
	return aws.NewClient(awsConfig)
}
