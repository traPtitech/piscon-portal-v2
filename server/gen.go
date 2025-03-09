package portal

// OpenAPI generator
//go:generate go tool ogen -config ogen.yaml -package openapi -target handler/openapi -clean ../openapi/openapi.yml

// ORM generator
//go:generate go tool bobgen-mysql -c bobgen.yaml
