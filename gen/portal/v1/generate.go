package portalv1

//go:generate go tool mockgen -package mock -destination mock/benchmark.go -typed . BenchmarkServiceClient

//go:generate go tool mockgen -package mock -destination mock/progress_stream.go -typed google.golang.org/grpc ClientStreamingClient
