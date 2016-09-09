#!/bin/bash

set -e 

apt-get update 
apt-get upgrade -y

apt-get install -y \
    htop \
    git \
    make \
    curl \
    supervisor \
    cgroup-lite \
    libapparmor-dev \
    libseccomp-dev \
    apparmor

export LOGFILE='/tmp/node-install.log'

echo 'Installing slave dependencies for ciru...'
(
    apt-get install -y \
        protobuf-c-compiler \
        libprotobuf-c0-dev \
        protobuf-compiler \
        libprotobuf-dev:amd64 \
        gcc \
        build-essential \
        bsdmainutils \
        python \
        git-core \
        asciidoc \
        xmlto \
        make 
) >> $LOGFILE

echo 'Building and installing criu...'
(
    cd /tmp
    git clone  https://github.com/xemul/criu.git
    cd /tmp/criu
    # checkout the current supported version at 1.5.2
    git checkout v1.5.2
    make clean
    make
    make install
    cd /tmp
    rm -rf /tmp/criu
) >> $LOGFILE

echo 'Criu built and installed successfully'
