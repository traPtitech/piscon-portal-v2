package main

import (
	"log"
	"time"

	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/runner"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	benchImpl "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl"
	privateisu "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl/private_isu"
	"github.com/traPtitech/piscon-portal-v2/runner/config"
	portalGrpc "github.com/traPtitech/piscon-portal-v2/runner/portal/grpc"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	problemExample    string = "example"
	problemPrivateIsu string = "private_isu"
)

var (
	problemBenchmarks = map[string]func(config map[string]any) benchmarker.Benchmarker{
		problemExample: func(_ map[string]any) benchmarker.Benchmarker {
			return benchImpl.NewExample()
		},
		problemPrivateIsu: func(_ map[string]any) benchmarker.Benchmarker {
			return privateisu.New()
		},
	}
)

func main() {
	config, err := config.LoadFile()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	conn, err := grpc.NewClient(config.Portal.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials())) //TODO: TLS を有効にする
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := portalv1.NewBenchmarkServiceClient(conn)
	p := portalGrpc.NewPortal(client, time.Second)

	bench, ok := problemBenchmarks[config.Problem.Name]
	if !ok {
		log.Fatalf("problem %q is not found", config.Problem.Name)
	}

	r := runner.Prepare(p, bench(config.Problem.Options))

	log.Printf("runner started: portal=%s, problem=%s\n", config.Portal.Address, config.Problem.Name)

	for {
		if err := r.Run(); err != nil {
			log.Printf("error: %v", err)
		}
	}
}
