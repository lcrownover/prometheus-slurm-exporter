# Prometheus Slurm Exporter

Prometheus collector and exporter for metrics extracted from the [Slurm](https://slurm.schedmd.com/overview.html) resource scheduling system.

This project was forked from [https://github.com/vpenso/prometheus-slurm-exporter](https://github.com/vpenso/prometheus-slurm-exporter) and, for now, aims to be backwards-compatible from SLURM 23.11 forward. This means the existing Grafana Dashboard should plug directly into this exporter and work roughly the same.

## Installation

This repository contains precompiled binaries for the three most recent major versions of SLURM. In the [releases](https://github.com/lcrownover/prometheus-slurm-exporter/releases) page, select the newest version of the exporter that matches your SLURM version. The included systemd file assumes you've saved this binary to `/usr/local/sbin/prometheus-slurm-exporter`.

## Configuration

The expoter requires several environment variables to be set:

* `SLURM_EXPORTER_LISTEN_ADDRESS`

This should be the full address for the exporter to listen on.

_Default: `0.0.0.0:8080`_

* `SLURM_EXPORTER_API_URL`

This is the URL to your slurmrestd server.

_Example: `http://head1.domain.edu:6820`_

* `SLURM_EXPORTER_API_TOKEN`

This is the [SLURM token to authenticate against slurmrestd](https://slurm.schedmd.com/jwt.html).

The easiest way to generate this is by running the following line on your head node:

```bash
scontrol token username=myuser lifespan=someseconds
```

`myuser` should probably be the `slurm` user, or some other privileged account.

`lifespan` is specified in seconds. I set mine for 1 year (`lifespan=31536000`).

* `SLURM_EXPORTER_API_USER`

The user specified in the token command.

## Systemd

A systemd unit file is [included](https://github.com/lcrownover/prometheus-slurm-exporter/blob/develop/extras/systemd/prometheus-slurm-exporter.service) for ease of deployment.

This unit file assumes you've written your environment variables to `/etc/prometheus-slurm-exporter/env.conf` in the format:

```
SLURM_EXPORTER_API_URL="http://head.domain.edu:6820"
SLURM_EXPORTER_API_USER="root"
SLURM_EXPORTER_API_TOKEN="mytoken"
```

## Prometheus Server Scrape Config

```
scrape_configs:
  - job_name: 'slurm_exporter'
    scrape_interval:  30s
    scrape_timeout:   30s
    static_configs:
      - targets: ['exporter_host.domain.edu:8080']
```

## Grafana Dashboard

The [dashboard](https://grafana.com/dashboards/4323) published by the previous author should work the same with this exporter.

![Status of the Nodes](images/Node_Status.png)

![Status of the Jobs](images/Job_Status.png)

![SLURM Scheduler Information](images/Scheduler_Info.png)

## Contributing

Check out the [CONTRIBUTING.md](CONTRIBUTING.md) document.
