# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

# Custom Fedora image with sshd installed, used for tests in 'ardc-config' package, for
# Ansible server.
FROM fedora:25
MAINTAINER Ren Lee <ren@arista.com>

RUN yum install -y \
   yum \
   sudo \
   openssh-server \
   python \
   hostname \
   sshpass \
   ftp://rpmfind.net/linux/fedora/linux/updates/25/x86_64/a/ansible-2.4.0.0-1.fc25.noarch.rpm \
   http://dist/ToolsV2/repo/x86_64/RPMS/telegraf-1.2.0-1.x86_64.rpm \
   http://dist/ToolsV2/repo/x86_64/RPMS/telegraf-Linux-1.8-15.noarch.rpm

EXPOSE 22

COPY ssh /root/.ssh
RUN chown root /root/.ssh
RUN chmod 700 /root/.ssh/
RUN chmod 600 /root/.ssh/*

COPY ssh/id_rsa /etc/ssh/ssh_host_rsa_key
COPY ssh/id_rsa.pub /etc/ssh/ssh_host_rsa_key.pub
RUN chmod 600 /etc/ssh/ssh_host_*

CMD ["/usr/sbin/sshd", "-D"]
