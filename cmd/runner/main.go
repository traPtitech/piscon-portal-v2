package main

import (
	"log"
	"time"

	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/runner"
	"github.com/traPtitech/piscon-portal-v2/runner/config"
	portalGrpc "github.com/traPtitech/piscon-portal-v2/runner/portal/grpc"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config, err := config.LoadFile()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return
	}

	conn, err := grpc.NewClient(config.Portal.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials())) //TODO: TLS を有効にする
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
		return
	}
	defer conn.Close()

	client := portalv1.NewBenchmarkServiceClient(conn)
	p := portalGrpc.NewPortal(client, time.Second)

	r, err := runner.Prepare(p, config.Problem)
	if err != nil {
		log.Fatalf("failed to prepare runner: %v", err)
		return
	}

	log.Printf("runner started: portal=%s, problem=%s\n", config.Portal.Address, config.Problem.Name)

	for {
		if err := r.Run(); err != nil {
			log.Printf("error: %v", err)
		}
	}
}
