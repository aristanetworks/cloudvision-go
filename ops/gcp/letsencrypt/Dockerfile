# Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

FROM ubuntu:16.04
MAINTAINER Fabrice Rabaute fabrice@arista.com

# Install nginx and certbot
RUN export DEBIAN_FRONTEND=noninteractive \
  && apt-get update \
  && apt-get install -y --no-install-recommends \
    at \
    curl \
    nginx \
    openssh-server \
    software-properties-common \
    wget \
  && add-apt-repository ppa:certbot/certbot \
  && apt-get update \
  && apt-get install -y --no-install-recommends python-certbot-nginx \
  && rm -rf /var/lib/apt/lists/* \
  && rm /etc/nginx/sites-enabled/default \
  && mkdir -p /root/.ssh \
  && chmod 700 /root/.ssh \
  && mkdir /var/run/sshd \
  && chmod 0755 /var/run/sshd

COPY start.sh .
COPY authorized_keys /root/.ssh
COPY ssh/ssh_host_dsa_key /etc/ssh
COPY ssh/ssh_host_dsa_key.pub /etc/ssh
COPY ssh/ssh_host_ecdsa_key /etc/ssh
COPY ssh/ssh_host_ecdsa_key.pub /etc/ssh
COPY ssh/ssh_host_ed25519_key /etc/ssh
COPY ssh/ssh_host_ed25519_key.pub /etc/ssh
COPY ssh/ssh_host_rsa_key /etc/ssh
COPY ssh/ssh_host_rsa_key.pub /etc/ssh

# Check if private keys are correct
RUN echo "29f9fa8ddc22304cf5018761e959c606  /etc/ssh/ssh_host_dsa_key" | md5sum --strict --check - \
  && echo "953536fb6400532dbe4395d35fa22948  /etc/ssh/ssh_host_ecdsa_key" | md5sum --strict --check - \
  && echo "f0efa2a72204afde251da34ba18c15d9  /etc/ssh/ssh_host_ed25519_key" | md5sum --strict --check - \
  && echo "72a8ef882f90eeae3f46fe196058b8eb  /etc/ssh/ssh_host_rsa_key" | md5sum --strict --check -

CMD ["./start.sh"]
