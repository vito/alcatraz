FROM ubuntu

RUN apt-get update

RUN apt-get -y install build-essential

ADD wsh /tmp/wsh
RUN cd /tmp/wsh && make
RUN mv /tmp/wsh/wshd /sbin/wshd

RUN apt-get -y install rsync

RUN useradd -m -s /bin/bash vcap

CMD /sbin/wshd --run /share/run
