#!/bin/bash
#SBATCH --job-name=lets_go
#SBATCH --output=/jobs/output/lets_go.out
#SBATCH --error=/jobs/err/lets_go.err
#SBATCH --time=00:05:00
#SBATCH --ntasks=1
echo "Lets Go"
sleep 180
