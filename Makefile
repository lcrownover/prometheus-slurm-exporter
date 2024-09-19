PROJECT_NAME = prometheus-slurm-exporter

24.05: openapi_24.05
	mkdir -p bin/
	go build -tags=2405 -o bin/prometheus-slurm-exporter cmd/prometheus-slurm-exporter/main.go

23.11:
	mkdir -p bin/
	go build -tags=2405 -o bin/prometheus-slurm-exporter cmd/prometheus-slurm-exporter/main.go

test: test_24.05 test_23.11
	go test -v ./...

test_24.05:
	go test -v ./... --tags=2405

test_23.11:
	go test -v ./... --tags=2311

run:
	go run cmd/prometheus-slurm-exporter/main.go

openapi_24.05:
	oapi-codegen --generate types --package types openapi-specs/24.05.json > internal/types/2405_openapi.gen.go

openapi_23.11:
	oapi-codegen --generate types --package types openapi-specs/23.11.json > internal/types/2311_openapi.gen.go
