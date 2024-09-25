#!/usr/bin/env python3

# This script is just a quick tool to generate a slurmrestd container
# and dump the latest openapi spec to ./openapi-specs

import subprocess
import sys

if len(sys.argv) != 2:
    print(f"usage: {sys.argv[0]} <slurm xx.xx version>")
    exit(1)

slurm_version = sys.argv[1]

versions = {
    "24.05": {
        "api_version": "0.0.41",
        "container_version": "24.05",
    },
    "23.11": {
        "api_version": "0.0.40",
        "container_version": "24.05",
    },
}

if slurm_version not in versions:
    print(
        "supported slurm versions: {}".format(", ".join([v for v in versions.keys()]))
    )
    exit(1)

oapi_version = versions[slurm_version]["api_version"]
container_version = versions[slurm_version]["container_version"]


def cleanup_container(container_version: str):
    container_delete_command = f"docker rm -f slurm-{container_version}"
    s = subprocess.run(
        container_delete_command.split(),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        universal_newlines=True,
    )

    if s.returncode != 0:
        raise Exception(f"Failed to clean up container: {s.stderr}")


def build_container(container_version: str):
    build_command = f"docker build -t slurm_{container_version} --file {container_version}.dockerfile ."
    s = subprocess.run(
        build_command.split(),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        universal_newlines=True,
    )

    if s.returncode != 0:
        raise Exception(f"Failed to build SLURM: {s.stderr}")


def create_container(container_version: str) -> str:
    create_command = (
        f"docker create --name slurm-{container_version} slurm_{container_version}"
    )
    s = subprocess.run(
        create_command.split(),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        universal_newlines=True,
    )

    if s.returncode != 0:
        raise Exception(f"Failed to create SLURM container: {s.stderr}")

    return s.stdout.strip()


def copy_container_file(container_id: str, oapi_version: str):
    copy_command = f"docker cp {container_id}:/slurm/v{oapi_version}.json ../openapi-specs/{slurm_version}.json"
    s = subprocess.run(
        copy_command.split(),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        universal_newlines=True,
    )

    if s.returncode != 0:
        raise Exception(f"Failed to copy Openapi specs from container: {s.stderr}")


print(
    f"Building SLURM {container_version} to get Openapi manifest version {oapi_version}"
)

try:
    build_container(container_version)
    container_id = create_container(container_version)
    copy_container_file(container_id, oapi_version)
    cleanup_container(container_version)
    print(f"Copied openapi spec {oapi_version} to ../openapi-specs/{slurm_version}.json")

except Exception as e:
    print(f"Failed to copy openapi spec: {e}")
    cleanup_container(container_version)
