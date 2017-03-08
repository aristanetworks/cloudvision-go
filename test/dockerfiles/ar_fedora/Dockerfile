# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

# Custom Fedora image with sshd installed, used for tests in 'ardc-config' package, for
# Ansible server.

# We need Fedora 18, but Docker library only has images for Fedora 20 and above.
FROM fedora:23
MAINTAINER Ren Lee <ren@arista.com>

# The following are installed:
# sudo - apparently this is missing from base image
# openssh - to host ssh connection to server
# python - dependency for ansible
# hostname - we want to be able to query the server's name
# sshpass - Ansible requires this
RUN yum install -y \
   sudo \
   openssh-server \
   python \
   hostname \
   sshpass

RUN yum install -y http://dist.aristanetworks.com/ToolsV2/repo/x86_64/RPMS/telegraf-1.0.0-1.x86_64.rpm
RUN yum install -y http://dist.aristanetworks.com/ToolsV2/repo/x86_64/RPMS/telegraf-Linux-1.6-3.noarch.rpm
RUN yum install -y ftp://195.220.108.108/linux/fedora/linux/updates/23/x86_64/a/ansible-2.2.0.0-3.fc23.noarch.rpm

EXPOSE 22

COPY ssh /root/.ssh
RUN chown root /root/.ssh
RUN chmod 700 /root/.ssh/
RUN chmod 600 /root/.ssh/*

COPY ssh/id_rsa /etc/ssh/ssh_host_rsa_key
COPY ssh/id_rsa.pub /etc/ssh/ssh_host_rsa_key.pub
RUN chmod 600 /etc/ssh/ssh_host_*

CMD ["/usr/sbin/sshd", "-D"]
