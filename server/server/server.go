package server

import (
	"fmt"
	"log"
	"net"

	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Start(port int, benchmarkSvc portalv1.BenchmarkServiceServer) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	s := grpc.NewServer()

	portalv1.RegisterBenchmarkServiceServer(s, benchmarkSvc)

	log.Println("gRPC server started on port", port)
	err = s.Serve(listener)
	if err != nil {
		return fmt.Errorf("grpc serve: %w", err)
	}
	defer s.Stop()
	return nil
}

func handleError(msg string, err error) error {
	log.Println(msg, err)
	return status.Error(codes.Internal, "internal error")
}
