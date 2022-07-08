#!/bin/bash

function generate() {
    NAME="$1"

    if [[ -f tls/"${NAME}".key ]]; then
        echo "Using existing certificate for ${NAME}"
        return
    else
        echo "Generating certificate for ${NAME}"
    fi

    openssl genrsa -out tls/"${NAME}".key 2048
    openssl req -new -key tls/"${NAME}".key -out tls/"${NAME}".csr -config tls/ca.conf -subj "/CN=${NAME}" -addext "subjectAltName = DNS:${NAME}" -extensions "${2}"
    openssl x509 -req -days 365 -in tls/"${NAME}".csr -CA tls/ca.crt -CAkey tls/ca.key -CAcreateserial -out tls/"${NAME}".crt -extfile tls/ca.conf -extensions "${2}"
}

if [[ ! -f tls/ca.crt || "$1" == "force" ]]; then
    echo "No TLS configuration found, generating self-signed certificates"
    openssl req -x509 -new -nodes -newkey rsa:2048 -keyout tls/ca.key -sha512 -days 1024 -out tls/ca.crt -config tls/ca.conf -extensions v3_ca
else
    echo "Using provided TLS configuration"
fi

generate pccs common
generate provisioner common
generate openvpn v3_vpn_server
generate client v3_vpn_client

echo ""
echo "Fixing permissions for pccs key"
chmod 0644 tls/pccs.key

echo ""
echo "Creating TA-Key for elliptic curves if not exists"
[[ ! -f tls/ta.key ]] && openvpn --genkey tls-auth tls/ta.key

echo ""
echo "TLS generation done, creating wrap-key for provisioning"
echo ""

if [[ ! -f tls/wrap-key ]]; then
    gramine-sgx-pf-crypt gen-key -w tls/wrap-key
else
    echo "Using existing wrap-key"
fi

echo "Your wrap-key: $(xxd -p tls/wrap-key)"
echo ""

[[ ! -d files/ ]] && mkdir -p files/

gramine-sgx-pf-crypt encrypt -w tls/wrap-key -i tls/openvpn.key -o files/openvpn.key
gramine-sgx-pf-crypt encrypt -w tls/wrap-key -i tls/ta.key -o files/ta.key

#chmod 0600 files/openvpn.key
#chmod 0600 files/ta.key
