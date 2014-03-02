FROM ubuntu

RUN apt-get update

RUN apt-get -y install build-essential
RUN apt-get -y install rsync
RUN apt-get -y install zlib1g-dev

ADD https://matt.ucc.asn.au/dropbear/releases/dropbear-2014.63.tar.bz2 /tmp/dropbear.tar.bz2

RUN tar jxf /tmp/dropbear.tar.bz2 -C /tmp
RUN mv /tmp/dropbear-2014.63 /tmp/dropbear

RUN cd /tmp/dropbear && ./configure && make
RUN cp /tmp/dropbear/dropbear /sbin

RUN mkdir -p /etc/dropbear
RUN /tmp/dropbear/dropbearkey -t rsa -f /etc/dropbear/dropbear_rsa_host_key

RUN rm -rf /tmp/dropbear

RUN useradd -m -s /bin/bash vcap

CMD /sbin/dropbear -F -E
