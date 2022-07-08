# enclaive-docker-openvpn-sgx

## About

This repository contains a fully working RA-TLS secret provisioning example using the [High-Level Secret Provisioning interface](https://gramine.readthedocs.io/en/latest/attestation.html#high-level-secret-provisioning-interface) of [gramine](https://github.com/gramineproject/gramine).

The application originally intended for demonstration is [OpenVPN](https://github.com/OpenVPN/openvpn), an open source VPN daemon. VPN applications like this one require the TUN/TAP kernel driver functionality which is currently not implemented in `gramine`, probably due to the underlying complexity.

Another possible candidate for enclaving the application would be [mystikos](https://github.com/deislabs/mystikos/), this was tried but also failed. For more information on this, see [mystikos.md](mystikos.md).

To see working secret provisioning after attestation, the `libos.entrypoint` in the manifest was changed to `/usr/bin/bash` to interactively explore the container. The required mount has been added and all files under `/usr/bin` are trusted so basic shell access is possible.

## Instructions

The first step is to build the container images using `docker-compose build`. There will be warnings about missing environment variables, which are only used during runtime and can be ignored for now.

The second command to be executed is the TLS-Certificate generation which will also generated necessary keys for file encryption:

```bash
./tls/cert-gen.sh
```

This will generate a self-signed CA and four certificates with private keys:

- PCCS, the Provisioning-Certificate-Caching-Service, created by [pccs/Dockerfile](pccs/Dockerfile)
- Provisioner, created by [provisioner/Dockerfile](provisioner/Dockerfile)
- OpenVPN server, which will be encrypted to `./files/`
- OpenVPN client

Both OpenVPN related certificates are using the appropriate extended key-usage attributes.

The script will also output the hex-encoded wrap-key used for the symmetric encryption of the files. Key generation and encryption is handled by `gramine-sgx-pf-crypt`.

The resulting image `enclaive/openvpn-sgx` will now be used to generate the signed manifest after mounting the signing key inside the container:

```bash
./gramine.sh
```

The required environment variables are later read from `.env`, a template can be found in `.env.template`.

- `PCCS_APIKEY` is the API-Key obtained from Intel after registration
- `KEY_DEFAULT` is the hex-encoded key generated in the first step
- `MRENCLAVE` and `MRSIGNER` are shown during the next step
- `ISV_PRODID` and `ISV_SVN` are configured in the manifest

`MRENCLAVE`, `MRSIGNER`, `ISV_PRODID` and `ISV_SVN` can be set to `false` to disable verification of these variables.

After setting all values in `.env`, the containers can be started.

To interactively explore the secret provisioning in the current stage without a working application, only the PCCS and Provisioner should be started detached:

```bash
docker-compose up -d pccs provisioner
docker-compose run ovpn
```

You will enter a bash shell inside of gramine.

The encrypted file can be read using `cat files/openvpn.key` mounted from `files/openvpn.key` and encrypted with the key provisioned to the application.

### TL;DR

```bash
# build the image
docker-compose build

# generate certificates, encryption key and encrypt files
./tls/cert-gen.sh

# sign certificate using host-mounted key
./gramine.sh

# configure using .env with KEY_DEFAULT from "Your wrap-key: REDACTED" during cert-gen.sh
# and setting MRENCLAVE, MRSIGNER, ISV_PRODID, ISV_SVN to false or to values from gramine.sh
nano .env

# launch pccs and provisioner
docker-compose up -d pccs provisioner

# start application container
docker-compose run ovpn

# inside the container
cat files/openvpn.key
```

## Explanation

The PCCS is using the reference implementation from Intel with a custom start-up script and certificate/key mounted from the host.

The provisioner is compiled against the `tools/ra-tls` libraries provided by `gramine`. It is basically a broken down version of the `ra-tls-secret-prov` [example](https://github.com/gramineproject/gramine/tree/master/CI-Examples/ra-tls-secret-prov) using Golang with CGO. A callback is passed to the library and as long as this callback returns a non-negative value the secret is provided to the application.

To configure the application in gramine to request and use this secret, the following settings are made inside the manifest:

```toml
loader.env.LD_PRELOAD = "/lib/libsecret_prov_attest.so"

loader.env.SECRET_PROVISION_CONSTRUCTOR = "1"
loader.env.SECRET_PROVISION_SET_KEY = "default"
loader.env.SECRET_PROVISION_CA_CHAIN_PATH = "/usr/local/share/ca-certificates/ca.crt"
loader.env.SECRET_PROVISION_SERVERS = "provisioner:4433"

sgx.remote_attestation = true
sgx.debug = false
```

The encrypted files are then mounted using the `encrytped`-mount type:

```toml
{ path = "/files", uri = "file:files", type = "encrypted", key_name = "default" },
```

`files` is a directory also mounted from the host containing the encrypted secrets generated during `cert-gen.sh`.

The secret provisioning server is started in [main.go](provisioner/src/main.go).

By using this approach a single provided secret can be used to decrypt as many files as necessary with minimal network traffic.

All secrets are only mounted during runtime and never copied into the image. All mounts can be seen in [docker-compose.yml](docker-compose.yml).

## About OpenVPN

Even if the TUN/TAP functionality would be implemented by `gramine` or `mystikos`, both would also need to implement [netlink-Sockets](https://en.wikipedia.org/wiki/Netlink). To avoid requiring this in addition the OpenVPN application is configured to use `iproute2` for network configuration with `--enable-iproute2`. As this would fork the `ip`-binary for the required calls, the file `openvpn_v2.5.6.diff` contains patches that disable these calls. As they are logged to the console, they can be set during the initialization in the entrypoint script `openvpn.sh`. The `ip`-binary would not even succeed as it requires the same `netlink`-sockets.

As an additional bug `if_nametoindex` is not implemented. As this value is constant, the function is simply overloaded to return this value.

During compilation and configuring OpenVPN support for the extended socket error capability is detected, this is also disable as it is not present when running with `gramine`. This is not a fatal error, but as we are already patching the application, we might as well do that too.

There is also another mode for OpenVPN where the device type is set to `null`, skipping all initialization. Disappointingly this mode is not correctly documented anywhere and the only useful reference I could find mentions that it is broken. When not using the TUN/TAP interface a `client-to-client` can be used but this is not the desired result.

## Other approaches to VPNs

One thing to look at could be IKEv2/IPSec based services using [strongSwan](https://github.com/strongswan/strongswan) but in the limited time for this project I was not able to explore this with a working and minimal configuration for internet routing.
