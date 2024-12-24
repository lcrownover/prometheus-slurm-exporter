#!/bin/bash
#SBATCH --job-name=lets_go
#SBATCH --output=/jobs/output/lets_go.out
#SBATCH --error=/jobs/err/lets_go.err
#SBATCH --ntasks=1
printf "Lets go\n"
