# Prometheus Slurm Exporter

Prometheus collector and exporter for metrics extracted from the [Slurm](https://slurm.schedmd.com/overview.html) resource scheduling system.

This project was forked from [https://github.com/vpenso/prometheus-slurm-exporter](https://github.com/vpenso/prometheus-slurm-exporter) and, for now, aims to be backwards-compatible from SLURM 23.11 forward.
This means the existing Grafana Dashboard should plug directly into this exporter and work roughly the same.

Unlike previous slurm exporters, this project leverages the SLURM REST API (`slurmrestd`) for data retreival.
Due to that difference, you are no longer required to run this exporter on a cluster node, as the exporter does not depend on having SLURM installed or connected to the head node!
I will be releasing containerized versions of this exporter soon.

## Installation

This repository contains precompiled binaries for the three most recent major versions of SLURM _(Note: currently only two versions, but will be three when 24.11 releases)_.
In the [releases](https://github.com/lcrownover/prometheus-slurm-exporter/releases) page, download the newest version of the exporter that matches your SLURM version.
The included systemd file assumes you've saved this binary to `/usr/local/sbin/prometheus-slurm-exporter`, so drop it there or take note to change the systemd file if you choose to use it.

## Configuration

The expoter requires several environment variables to be set:

* `SLURM_EXPORTER_LISTEN_ADDRESS`

  This should be the full address for the exporter to listen on.

  _Default: `0.0.0.0:8080`_

* `SLURM_EXPORTER_API_URL`

  This is the URL to your slurmrestd server.

  _Example: `http://head1.domain.edu:6820`_
  _Example: `unix://path/to/unix/socket`_

* `SLURM_EXPORTER_API_USER`

  The user specified in the token command.

* `SLURM_EXPORTER_API_TOKEN`

  This is the [SLURM token to authenticate against slurmrestd](https://slurm.schedmd.com/jwt.html).

  The easiest way to generate this is by running the following line on your head node:

  ```bash
  scontrol token username=myuser lifespan=someseconds
  ```

  `myuser` should probably be the `slurm` user, or some other privileged account.

  `lifespan` is specified in seconds. I set mine for 1 year (`lifespan=31536000`).

* `SLURM_EXPORTER_ENABLE_TLS`

  Set to `true` to enable TLS support. You must also provide paths to your certificate and key.

* `SLURM_EXPORTER_TLS_CERT_PATH`

  Path to your TLS certificate.

* `SLURM_EXPORTER_TLS_KEY_PATH`

  Path to your TLS key, it should be `0600`.

## Systemd

A systemd unit file is [included](https://github.com/lcrownover/prometheus-slurm-exporter/blob/develop/extras/systemd/prometheus-slurm-exporter.service) for ease of deployment.

This unit file assumes you've written your environment variables to `/etc/prometheus-slurm-exporter/env.conf` in the format:

```
SLURM_EXPORTER_API_URL="http://head.domain.edu:6820"
SLURM_EXPORTER_API_USER="root"
SLURM_EXPORTER_API_TOKEN="mytoken"
```

_Don't forget to `chmod 600 /etc/prometheus-slurm-exporter/env.conf`!_

## Prometheus Server Scrape Config

This is an example scrape config for your prometheus server:

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
I will be releasing a new version of the dashboard soon that will receive new features.

![Status of the Nodes](images/Node_Status.png)

![Status of the Jobs](images/Job_Status.png)

![SLURM Scheduler Information](images/Scheduler_Info.png)

## Contributing

Check out the [CONTRIBUTING.md](CONTRIBUTING.md) document.
