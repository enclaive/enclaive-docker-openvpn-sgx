#!/bin/bash

DIRECTORY=$(dirname "${BASH_SOURCE[0]}")
cd "${DIRECTORY}"

[[ ! -f /root/.config/gramine/enclave-key.pem ]] && gramine-sgx-gen-private-key

gramine-manifest -Darch_libdir=/lib/x86_64-linux-gnu ovpn.manifest.template ovpn.manifest
gramine-sgx-sign --manifest ovpn.manifest --output ovpn.manifest.sgx
gramine-sgx-get-token -s ovpn.sig -o ovpn.token
