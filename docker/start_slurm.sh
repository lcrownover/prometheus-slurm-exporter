#!/bin/bash

# Ensure logfile and /var/log/munge have the correct ownership
# Start munge daemon as munge user
echo "Starting the munge daemon"
sudo -u munge /usr/sbin/munged

# Check if munge daemon started successfully
if ps aux | grep -q '[m]unged'; then
  echo "Munge daemon started successfully"
else
  echo "Failed to start Munge daemon"
  exit 1
fi

# Output the Slurm configuration
#cat /usr/local/etc/slurm.conf
# Start the slurmctld daemon
echo "Starting the slurmctld daemon"
if /slurm/src/slurmctld/slurmctld; then
  echo "slurmctld daemon started successfully"
else
  echo "Failed to start slurmctld daemon"
  exit 1
fi

# Start the slurmd daemon
echo "Starting the slurmd daemon"
if /slurm/src/slurmd/slurmd/slurmd; then
  echo "slurmd daemon started successfully"
else
  echo "Failed to start slurmd daemon"
  exit 1
fi

echo "Submitting jobs"
sbatch /jobs/hello_world_job.sbatch
if [ $? -eq 0 ]; then
  echo "hello_world_job.sbatch submitted successfully"
else
  echo "Failed to submit hello_world_job.sbatch"
fi

sbatch /jobs/lets_go_job.sbatch
if [ $? -eq 0 ]; then
  echo "lets_go_job.sbatch submitted successfully"
else
  echo "Failed to submit lets_go_job.sbatch"
fi

# Keep the container running
#tail -f /dev/null
