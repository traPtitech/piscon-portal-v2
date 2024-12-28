package portal

// OpenAPI generator
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@v1.8.1 -config ogen.yaml -package openapi -target handler/openapi -clean ../openapi/openapi.yml

// ORM generator
//go:generate go run github.com/stephenafamo/bob/gen/bobgen-mysql@v0.29.0 -c bobgen.yaml
