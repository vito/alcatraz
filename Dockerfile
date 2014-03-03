FROM ubuntu

RUN apt-get update

RUN apt-get -y install build-essential
RUN apt-get -y install rsync
RUN apt-get -y install openssh-server

RUN useradd -m -s /bin/bash vcap

RUN mkdir /var/run/sshd

RUN locale-gen en_US en_US.UTF-8

CMD /usr/sbin/sshd -D -e
