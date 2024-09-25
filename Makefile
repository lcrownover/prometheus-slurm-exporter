PROJECT_NAME = prometheus-slurm-exporter

ifndef SLURM_VERSION
$(error SLURM_VERSION environment variable is not set)
endif

slurm_version := ${SLURM_VERSION}

# If SLURM_VERSION is "all", print an error message for the default build target
ifeq ($(slurm_version),all)
build:
	$(error You must unset SLURM_VERSION to build)
else
build:
	mkdir -p bin/
	go build -tags=$(subst .,,$(slurm_version)) -o bin/prometheus-slurm-exporter cmd/prometheus-slurm-exporter/main.go
endif

test:
ifeq ($(slurm_version),all)
	# Generate and test for version 24.05
	go test -tags=2405 -v ./...
	# Generate and test for version 23.11
	go test -tags=2311 -v ./...
else
	go test -tags=$(subst .,,$(slurm_version)) -v ./...
endif

install:
	cp bin/prometheus-slurm-exporter /usr/local/sbin/prometheus-slurm-exporter
