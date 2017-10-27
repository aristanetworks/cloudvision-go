FROM fedora:20
MAINTAINER Adam Setters <asetty@arista.com>

RUN yum install -y \
   python \
   openssh-server \
   MySQL-python.x86_64

RUN useradd -m arastra
RUN useradd -m asb

CMD [ "/bin/bash" ]
