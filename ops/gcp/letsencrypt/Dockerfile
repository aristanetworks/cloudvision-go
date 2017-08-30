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

RUN chmod 600 /root/.ssh/authorized_keys

CMD ["./start.sh"]
