# Development

You must have access to a slurm head node running `slurmrestd` and a valid token
for that service.

## Install Go from source

```bash
export VERSION=1.22.5 OS=linux ARCH=amd64
wget https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz
tar -xzvf go$VERSION.$OS-$ARCH.tar.gz
export PATH=$PWD/go/bin:$PATH
```

_Alternatively install Go using the packaging system of your Linux distribution._

## Clone this repository and build

Use Git to clone the source code of the exporter, run all the tests and build the binary:

```bash
git clone https://github.com/lcrownover/prometheus-slurm-exporter.git
cd prometheus-slurm-exporter
make
```

To just run the tests:

```bash
make test
```

Start the exporter (foreground), and query all metrics:

```bash
./bin/prometheus-slurm-exporter
```

If you wish to run the exporter on a different port, or the default port (8080) is already in use, run with the following argument:

```bash
./bin/prometheus-slurm-exporter --listen-address="0.0.0.0:<port>"
...

# query all metrics (default port)
curl http://localhost:8080/metrics
```
