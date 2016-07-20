# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

# Custom Fedora image with sshd installed, used for tests in 'ardc-config' package, for
# Ansible server.

# We need Fedora 18, but Docker library only has images for Fedora 20 and above.
FROM fedora:20
MAINTAINER Ren Lee <ren@arista.com>

# To get around the Rpmdb checksum issue (https://github.com/docker/docker/issues/10180)
yum clean all

# The following are installed:
# openssh - to host ssh connection to server
# python - dependency for ansible
# hostname - we want to be able to query the server's name
# sshpass - Ansible requires this
RUN yum install -y \
   openssh-server \
   python \
   hostname \
   sshpass 

# ansible - we want to run test playbooks against this server
# This is roundabout way to avoid needing to keep up with fedora archive updates
COPY ansible-2.1.0.0-1.fc23.noarch.rpm /tmp
RUN yum install -y /tmp/ansible-2.1.0.0-1.fc23.noarch.rpm

EXPOSE 22

COPY ssh /root/.ssh
RUN chown root /root/.ssh
RUN chmod 700 /root/.ssh/
RUN chmod 600 /root/.ssh/*

COPY ssh/id_rsa /etc/ssh/ssh_host_rsa_key
COPY ssh/id_rsa.pub /etc/ssh/ssh_host_rsa_key.pub
RUN chmod 600 /etc/ssh/ssh_host_*

CMD ["/usr/sbin/sshd", "-D"]
