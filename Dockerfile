FROM enclaive/gramine-os:impish-1.3-425e8450

RUN apt-get update && apt-get install -y build-essential
RUN apt-get install -y autoconf iproute2 git
RUN apt-get install -y libtool-bin libssl-dev liblz4-dev liblzo2-dev libpam0g-dev python3-docutils

RUN git clone --depth=1 --branch=v2.5.6 https://github.com/OpenVPN/openvpn

WORKDIR ./openvpn/

RUN sed -i '/^CONFIGURE_DEFINES=/s/set/env/g' configure.ac
RUN autoreconf --force --install

RUN ./configure --prefix=/usr --sbindir=/usr/sbin --enable-iproute2
RUN make -j4

COPY ./openvpn_v2.5.6.diff /openvpn/
RUN patch -p1 < ./openvpn_v2.5.6.diff \
    && make -j4 \
    && make -j4 install \
    && libtool --finish /usr/lib/openvpn/plugins

WORKDIR /

COPY ./openvpn.sh .

EXPOSE 1194/udp

ENTRYPOINT [ "/openvpn.sh" ]
