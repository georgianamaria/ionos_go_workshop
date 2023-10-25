
CODEGEN ?= "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"

internal/api/quotav1/api.gen.go: openapi/quota.yaml openapi/config.yaml
	go run $(CODEGEN) -config openapi/config.yaml $< > $@
