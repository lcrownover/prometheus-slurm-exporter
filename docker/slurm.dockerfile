FROM nathanhess/slurm:full-root

RUN apt-get update && apt-get install -y \
    systemd \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/systemd/system/multi-user.target.wants

COPY slurm.conf /etc/slurm/slurm.conf
#COPY cgroup.conf /etc/slurm/cgroup.conf


ARG CPU=4
ARG MEMORY=8192

RUN echo "Container OS:" && cat /etc/os-release

RUN mkdir -p /container/jobs /container/output /container/err

WORKDIR /container/jobs

RUN echo '#!/bin/bash\n\
#SBATCH --job-name=hello_world\n\
#SBATCH --output=/container/output/hello_world.out\n\
#SBATCH --error=/container/err/hello_world.err\n\
#SBATCH --time=00:05:00\n\
#SBATCH --ntasks=1\n\n\
echo "Hello World"\n\
sleep 300' > /container/jobs/hello_world_job.sbatch

RUN echo '#!/bin/bash\n\
#SBATCH --job-name=lets_go\n\
#SBATCH --output=/container/output/lets_go.out\n\
#SBATCH --error=/container/err/lets_go.err\n\
#SBATCH --time=00:05:00\n\
#SBATCH --ntasks=1\n\n\
echo "Let'\''s Go"\n\
sleep 300' > /container/jobs/lets_go_job.sbatch

RUN chmod +x /container/jobs/hello_world_job.sbatch /container/jobs/lets_go_job.sbatch

CMD ["slurmctld", "slurmd", "-N"]
COPY start_slurm.sh /usr/local/bin/start_slurm.sh
#EXPOSE 6280
RUN chmod +x /usr/local/bin/start_slurm.sh
CMD ["/usr/local/bin/start_slurm.sh"]
CMD ["/bin/systemd"]
COPY slurm.conf /etc/slurm/slurm.conf
COPY cgroup.conf /etc/slurm/cgroup.conf
