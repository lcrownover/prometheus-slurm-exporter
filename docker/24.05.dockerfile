FROM rockylinux:8
RUN dnf update -y && \
    dnf install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm && \
    dnf install -y --enablerepo=devel mariadb-devel python3-PyMySQL hwloc lz4-devel wget bzip2 perl munge-devel munge cmake jansson libjwt-devel libjwt json-c-devel json-c http-parser-devel http-parser && \
    dnf group install -y "Development Tools"

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

COPY slurm.conf /usr/local/etc/slurm.conf

WORKDIR /slurm
RUN wget https://download.schedmd.com/slurm/slurm-24.05-latest.tar.bz2 && tar -xvjf slurm-24.05-latest.tar.bz2 --strip-components=1

RUN /usr/sbin/create-munge-key && \
    chown munge:munge /etc/munge/munge.key && \
    chmod 400 /etc/munge/munge.key


RUN ./configure \
    --with-http-parser=/usr/local/ \
    --enable-slurmrestd \
    && make && make install

RUN /usr/sbin/munged

EXPOSE 6280

RUN env SLURM_CONF=/dev/null slurmrestd -d v0.0.41 -s slurmdbd,slurmctld --generate-openapi-spec > /slurm/v0.0.41.json
RUN env SLURM_CONF=/dev/null slurmrestd -d v0.0.40 -s slurmdbd,slurmctld --generate-openapi-spec > /slurm/v0.0.40.json
