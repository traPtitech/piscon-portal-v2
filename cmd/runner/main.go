package main

import (
	"log"
	"time"

	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/runner"
	portalGrpc "github.com/traPtitech/piscon-portal-v2/runner/portal/grpc"
)

func main() {
	client := portalv1.NewBenchmarkServiceClient(nil) //TODO: Implement
	p := portalGrpc.NewPortal(client, time.Second)

	r := runner.Prepare(p, nil) //TODO: Implement

	for {
		if err := r.Run(); err != nil {
			log.Printf("error: %v", err)
		}
	}
}
