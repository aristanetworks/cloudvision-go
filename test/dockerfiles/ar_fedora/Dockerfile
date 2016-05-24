# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

# Custom Fedora image with sshd installed, used for tests in 'ardc-config' package, for
# Ansible server.

# We need Fedora 18, but Docker library only has images for Fedora 20 and above.
FROM fedora:20
MAINTAINER Ren Lee <ren@arista.com>

# The following are installed:
# openssh - to host ssh connection to server
# python - dependency for ansible
# hostname - we want to be able to query the server's name
# sshpass - Ansible requires this
# ansible - we want to run test playbooks against this server
RUN yum install -y \
   openssh-server \
   python \
   hostname \
   sshpass \
   ansible

EXPOSE 22
ENV ANSIBLE_INVENTORY=/test_ansible_hosts

COPY ssh /root/.ssh
RUN chown root /root/.ssh
RUN chmod 700 /root/.ssh/
RUN chmod 600 /root/.ssh/*

COPY ssh/id_rsa /etc/ssh/ssh_host_rsa_key
COPY ssh/id_rsa.pub /etc/ssh/ssh_host_rsa_key.pub
RUN chmod 600 /etc/ssh/ssh_host_*

CMD ["/usr/sbin/sshd", "-D"]
