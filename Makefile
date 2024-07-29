PROJECT_NAME = prometheus-slurm-exporter

all:
	mkdir -p bin/
	go build -o bin/prometheus-slurm-exporter cmd/prometheus-slurm-exporter/main.go

test:
	go test -v ./...

run:
	go run cmd/prometheus-slurm-exporter/main.go

# You can get openapi.json from your slurmrestd at: localhost:6820/openapi.json
openapi:
	oapi-codegen --package=types --generate types openapi.json > internal/types/openapi.gen.go
