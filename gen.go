package portal

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@v1.8.1 -config ogen.yaml -package openapi -target server/openapi -clean openapi/openapi.yml
//go:generate go run github.com/stephenafamo/bob/gen/bobgen-mysql@v0.29.0 -c bobgen.yaml
