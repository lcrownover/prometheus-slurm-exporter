#!/bin/bash
set -e

# Start slurmctld
slurmctld &

# Start slurmd
slurmd & 

# Keep the container running
wait

