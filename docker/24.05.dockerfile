FROM rockylinux:8
RUN dnf update -y && \
    dnf install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm && \
    dnf install -y --enablerepo=devel mariadb-devel python3-PyMySQL hwloc lz4-devel wget bzip2 perl munge-devel munge cmake jansson libjwt-devel libjwt json-c-devel json-c http-parser-devel http-parser libcgroup libcgroup-tools dbus-devel mariadb && \
    dnf group install -y "Development Tools"

RUN dnf install -y sudo

RUN dnf -y update && \
    dnf install -y systemd && \
    dnf clean all && \
    rm -rf /var/lib/apt/lists/*

RUN adduser slurm

# Install http_parser
RUN git clone --depth 1 --single-branch -b v2.9.4 https://github.com/nodejs/http-parser.git http_parser \
    && cd http_parser \
    && make \
    && make install

#RUN dnf install -y systemd

WORKDIR /slurm
RUN wget https://download.schedmd.com/slurm/slurm-24.05-latest.tar.bz2 && tar -xvjf slurm-24.05-latest.tar.bz2 --strip-components=1
RUN ./configure \
    --with-cgroup-v2 \
    --with-http-parser=/usr/local/ \
    --enable-slurmrestd \
    && make && make install

# Create the /var/log/slurm directory and set permissions
RUN mkdir -p /var/log/slurm && \
    chown slurm:slurm /var/log/slurm && \
    chmod 750 /var/log/slurm && \
    touch /var/log/slurm/slurmd.log /var/log/slurm/slurmctld.log && \
    chown slurm:slurm /var/log/slurm/slurmctld.log /var/log/slurm/slurmd.log

RUN getent group munge || groupadd -r munge && \
    getent passwd munge || useradd -r -g munge munge && \
    mkdir -p /var/log/munge && \
    chown munge:munge /var/log/munge && \
    chmod 750 /var/log/munge && \
    /usr/sbin/create-munge-key && \
    chown munge:munge /etc/munge/munge.key && \
    chmod 400 /etc/munge/munge.key

RUN touch /var/log/munge/munged.log && \
    chown munge:munge /var/log/munge/munged.log

COPY slurm.conf /usr/local/etc/slurm.conf

USER root
COPY cgroup.conf /usr/local/etc/cgroup.conf
COPY slurm.conf /usr/local/etc/slurm.conf
COPY slurmdbd.conf /usr/local/etc/slurmdbd.conf
RUN chown slurm:slurm /usr/local/etc/slurmdbd.conf
RUN chmod 600 /usr/local/etc/slurmdbd.conf
COPY start_slurm.sh /start_slurm.sh

ENV SLURM_CONF=/usr/local/etc/slurm.conf
RUN chmod 755 /start_slurm.sh

RUN mkdir -p /var/spool/slurm /var/spool/slurmd && \
    chown slurm:slurm /var/spool/slurm /var/spool/slurmd && \
    chmod 750 /var/spool/slurmd && \
    touch /var/spool/slurmd/cred_state && \
    chown slurm:slurm /var/spool/slurmd/cred_state

RUN chown -R slurm:slurm /slurm/src/ 

RUN mkdir -p /jobs /jobs/output /jobs/err && \
    chown root:slurm /jobs /jobs/output /jobs/err

# Create sample SLURM job scripts

COPY hello_world_job.sbatch /jobs/hello_world_job.sbatch
COPY lets_go_job.sbatch /jobs/lets_go_job.sbatch

RUN chmod +x /jobs/hello_world_job.sbatch /jobs/lets_go_job.sbatch

# Ask Lucas about what other ports need to be exposed or if I need to build slurm with this port exposed from the getgo 
EXPOSE 6280

RUN ln -s /slurm/src/slurmd/slurmd/slurmd /bin/slurmd # I only added this to make it easier to run the slurmd executable during daemon start troubleshooting 
RUN ln -s /slurm/src/slurmdbd/slurmdbd /bin/slurmdbd # I only added this to make it easier to run the slurmd executable during daemon start troubleshooting 
RUN ln -s /slurm/src/slurmrestd/slurmrestd /bin/slurmrestd # I only added this to make it easier to run the slurmd executable during daemon start troubleshooting 

RUN env SLURM_CONF=/dev/null slurmrestd -d v0.0.41 -s slurmdbd,slurmctld --generate-openapi-spec > /slurm/v0.0.41.json
RUN env SLURM_CONF=/dev/null slurmrestd -d v0.0.40 -s slurmdbd,slurmctld --generate-openapi-spec > /slurm/v0.0.40.json

ENTRYPOINT ["/start_slurm.sh"]
