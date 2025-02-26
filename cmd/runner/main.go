package main

import (
	"log"
	"time"

	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/runner"
	benchImpl "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl"
	portalGrpc "github.com/traPtitech/piscon-portal-v2/runner/portal/grpc"
)

func main() {
	client := portalv1.NewBenchmarkServiceClient(nil) //TODO: Implement
	p := portalGrpc.NewPortal(client, time.Second)

	benchmarker := benchImpl.Example() //TODO: 設定を読み込んで動的に変えるようにする

	r := runner.Prepare(p, benchmarker)

	for {
		if err := r.Run(); err != nil {
			log.Printf("error: %v", err)
		}
	}
}
