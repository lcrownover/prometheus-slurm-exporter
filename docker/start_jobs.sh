#!/bin/bash

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
