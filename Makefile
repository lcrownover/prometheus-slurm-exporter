PROJECT_NAME = prometheus-slurm-exporter

all:
	mkdir -p bin/
	go build -o bin/prometheus-slurm-exporter cmd/prometheus-slurm-exporter/main.go

test:
	go test -v ./...

run:
	go run cmd/prometheus-slurm-exporter/main.go

openapi:
	oapi-codegen -config=oapi-codegen-config.yml openapi-specs/v0.0.41.json
