#!/bin/bash

set -eu

[[ -z "$APIKEY" ]] && echo "APIKEY was not set" && exit 1

export HASH_USER=$( head /dev/random | sha512sum | awk '{print $1;}')
export HASH_ADMIN=$(head /dev/random | sha512sum | awk '{print $1;}')

cat <<< $(
    jq '.hosts = "0.0.0.0" | .ApiKey = $ENV.APIKEY | .UserTokenHash = $ENV.HASH_USER | .AdminTokenHash = $ENV.HASH_ADMIN' config/default.json
) > config/default.json

/usr/bin/node -r esm /opt/intel/sgx-dcap-pccs/pccs_server.js
