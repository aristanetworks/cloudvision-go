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

# Allow SSH login as root, with this generic password. 
RUN echo 'root:arista' | chpasswd
RUN sed -i 's/PermitRootLogin without-password/PermitRootLogin yes/' /etc/ssh/sshd_config

EXPOSE 22
ENV ANSIBLE_INVENTORY=/test_ansible_hosts

RUN /usr/bin/ssh-keygen -A
RUN mkdir /.ssh
COPY test_authorized_keys /.ssh/authorized_keys
COPY test_id_rsa /.ssh/id_rsa
COPY test_id_rsa.pub /.ssh/id_rsa.pub

COPY test_id_rsa /etc/ssh/ssh_host_rsa_key
COPY test_id_rsa.pub /etc/ssh/ssh_host_rsa_key.pub

CMD ["/usr/sbin/sshd", "-D"]
