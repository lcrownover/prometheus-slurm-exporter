#!/bin/bash
#SBATCH --job-name=hello_world
#SBATCH --output=/container/output/hello_world.out
#SBATCH --error=/container/err/hello_world.err
#SBATCH --time=00:05:00
#SBATCH --ntasks=1
echo "Hello World"
sleep 300
