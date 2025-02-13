package portalv1

//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -package mock -destination mock/benchmark.go -typed . BenchmarkServiceClient

//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -package mock -destination mock/progress_stream.go -typed google.golang.org/grpc ClientStreamingClient
