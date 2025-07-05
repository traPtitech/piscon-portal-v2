package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/pflag"
	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/runner"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	benchImpl "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl"
	privateisu "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl/private_isu"
	portalGrpc "github.com/traPtitech/piscon-portal-v2/runner/portal/grpc"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	problemExample    string = "example"
	problemPrivateIsu string = "private_isu"
)

var (
	problems = []string{problemExample, problemPrivateIsu}
	target   = pflag.StringP("target", "t", "", "portal server address (host:port)")
	problem  = pflag.StringP("problem", "p", "", fmt.Sprintf("problem name: one of %v", problems))
	help     = pflag.BoolP("help", "h", false, "show help (this message)")

	problemBenchmarks = map[string]func(config map[string]any) benchmarker.Benchmarker{
		problemExample: func(_ map[string]any) benchmarker.Benchmarker {
			return benchImpl.NewExample()
		},
		problemPrivateIsu: func(_ map[string]any) benchmarker.Benchmarker {
			return privateisu.New()
		},
	}
)

func validateFlags() error {
	if target == nil || *target == "" {
		return errors.New("target is required")
	}
	if problem == nil || *problem == "" {
		return errors.New("problem is required")
	}
	return nil
}

func main() {
	pflag.Parse()

	if help != nil && *help {
		pflag.Usage()
		return
	}

	if err := validateFlags(); err != nil {
		pflag.Usage()
		log.Fatalf("validation error: %v", err)
	}

	conn, err := grpc.NewClient(*target,
		grpc.WithTransportCredentials(insecure.NewCredentials())) //TODO: TLS を有効にする
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := portalv1.NewBenchmarkServiceClient(conn)
	p := portalGrpc.NewPortal(client, time.Second)

	benchConfig := map[string]any{}

	bench, ok := problemBenchmarks[*problem]
	if !ok {
		log.Fatalf("problem %q is not found", *problem)
	}

	r := runner.Prepare(p, bench(benchConfig))

	log.Printf("runner started: target=%s, problem=%s\n", *target, *problem)

	for {
		if err := r.Run(); err != nil {
			log.Printf("error: %v", err)
		}
	}
}
