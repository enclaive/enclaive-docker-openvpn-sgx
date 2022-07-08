#!/bin/bash

update-ca-certificates

/aesmd.sh

sleep 1

cd /manifest/

mkdir -p /dev/net
mknod /dev/net/tun c 10 200
chmod 600 /dev/net/tun

openvpn --mktun --dev ovpn --dev-type tun --user root --group root

################################################
# these are extracted after first start from log
/usr/sbin/ip link set dev ovpn up mtu 1500
/usr/sbin/ip link set dev ovpn up
/usr/sbin/ip addr add dev ovpn 10.8.0.1/24
################################################

#chown root:root files/*

#gramine-sgx ovpn server.conf
gramine-sgx ovpn
