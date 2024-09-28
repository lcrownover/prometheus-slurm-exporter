# Development

You must have access to a slurm head node running `slurmrestd` and a valid token
for that service. Take note of your slurm version, such as `24.05`, as you'll
use this version when building.

`develop` is the default branch for this repository, and `main` is used for
releases.

## Install Go from source

```bash
export VERSION=1.22.5 OS=linux ARCH=amd64
wget https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz
tar -xzvf go$VERSION.$OS-$ARCH.tar.gz
export PATH=$PWD/go/bin:$PATH
```

_Alternatively install Go using the packaging system of your Linux
distribution._

## Adding Support for New Openapi Versions

### Install openapi-generator-cli and openjdk

Install `openapi-generator-cli` globally with NPM:

```bash
npm install -g @openapitools/openapi-generator-cli`
```

This package depends on having the `java` executable in `PATH`, so install java.

For mac, `brew install java`, then following the brew message, symlink the JDK,
`sudo ln -sfn /usr/local/opt/openjdk/libexec/openjdk.jdk /Library/Java/JavaVirtualMachines/openjdk.jdk`

For ubuntu, `sudo snap install openjdk`.

## Building

### Clone this repository and build

Use Git to clone the source code:

```bash
git clone https://github.com/lcrownover/prometheus-slurm-exporter.git
cd prometheus-slurm-exporter
```

Build the binary for your SLURM version, for example 24.05:

```bash
SLURM_VERSION=24.05 make
```

Run tests for a specific SLURM version:

```bash
SLURM_VERSION=24.05 make test
```

Run the tests for all SLURM versions:

```bash
SLURM_VERSION=all make test
```

Start the exporter:

```bash
./bin/prometheus-slurm-exporter
```

If you wish to run the exporter on a different port, or the default port (8080)
is already in use, run with the following argument:

```bash
./bin/prometheus-slurm-exporter --listen-address="0.0.0.0:<port>"
```

Query all metrics:

```bash
curl http://localhost:8080/metrics
```

### Generating and Saving Openapi specs from SLURM using Docker

Navigate to the `docker` directory and use the python script to automatically
grab and store an openapi yaml spec from a target slurm version.

### Generating the Openapi code for new SLURM versions

I do this for every new SLURM version, so it should already be done.

Assuming 23.11:

```bash
openapi-generator-cli generate \
    -g go \
    -i openapi-specs/23.11.json \
    -o ../openapi-slurm-23-11 \
    --package-name openapi_slurm_23_11 \
    --git-user-id lcrownover \
    --git-repo-id openapi-slurm-23-11
```

This will generate an entire git repository that you can toss up in GitHub.

Navigate to it.

Now we just strip out all instances of the API version. In this case, 23.11 is API
version 0.0.40:

```bash
find . -type f -not -path "./.git/*" -exec gsed -i 's/V0040//g' {} +
```

### Cutting releases

1. Update the `version` var in `main.go`

2. Merge `develop` into `main`

3. Tag the release

   `git tag -a v1.0.1 -m 'release note'`

4. Push the tag

   `git push origin v1.0.1`

5. Make sure you have `GITHUB_TOKEN` exported, then use `goreleaser` to create
   releases:

   `goreleaser release`

## Ideas

docker container

end-to-end testing with docker slurm cluster
