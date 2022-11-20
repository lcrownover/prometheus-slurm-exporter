PROJECT_NAME = prometheus-slurm-exporter

all:
	mkdir -p bin/
	go build -o bin/prometheus-slurm-exporter

test:
	go test -v

run:
	go run cmd/prometheus-slurm-exporter/main.go
