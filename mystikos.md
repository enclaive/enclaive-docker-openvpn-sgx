# enclaive-docker-openvpn-sgx

This file contains basic instructions on installing mystikos, compiling mininmal OpenVPN and packaging it to run with mystikos. It is currently broken, as the required `/dev/net/tun` can neither be mounted inside the process nor created by the process as [mknod(2)](https://man7.org/linux/man-pages/man2/mknod.2.html) is not implemented.

## Mystikos

```bash
apt-get update
apt-get install -y wget libsgx-enclave-common libsgx-dcap-ql libsgx-dcap-ql-dev libsgx-quote-ex libmbedtls-dev

wget -q https://github.com/deislabs/mystikos/releases/download/v0.9.0/Ubuntu-2004_mystikos-0.9.0-x86_64.deb
dpkg -i Ubuntu-2004_mystikos-0.9.0-x86_64.deb

mkdir -p /app/
export PATH="$PATH:/opt/mystikos/bin"

apt-get install -y nano file strace
apt-get install -y build-essential autoconf iproute2 git libtool-bin libssl-dev liblz4-dev liblzo2-dev libpam0g-dev python3-docutils pkg-config musl-tools

git clone --depth=1 --branch=v2.5.6 https://github.com/OpenVPN/openvpn

cd /openvpn/
sed -i '/^CONFIGURE_DEFINES=/s/set/env/g' configure.ac
autoreconf --force --install
./configure --prefix=/usr --sbindir=/usr/sbin --disable-lzo --disable-lz4 \
    --disable-plugins --disable-management --disable-fragment --disable-multihome \
    --disable-port-share --disable-debug --disable-def-auth --disable-pf \
    --disable-ofb-cfb --disable-plugin-auth-pam --disable-plugin-down-root \
    --with-crypto-library=mbedtls
make -j4
make -j4 install
cd -

cp /usr/sbin/openvpn /app/
while read lib;do
    DIRECTORY=$(dirname -- "$lib")
    mkdir -p /app/"$DIRECTORY"
    cp "$lib" /app/"$DIRECTORY"
done < <(ldd app/openvpn | grep -o '/lib/[^ ]*')

myst mkcpio /app/ rootfs
myst exec-sgx --app-config-path config.json --mount /dev/net/=/dev/net/ rootfs /openvpn
```

The `config.json` is a file containing the application configuration including host mounts. None of them work. This does not even try to run the application and will only try to create the necessary device.

```json
{
    "Debug": 1,
    "ApplicationPath": "/openvpn",
    "ApplicationParameters": ["--mktun", "--dev", "tun0"],
    "Mount": [
        {
        "Target": "/dev/net/",
        "Type": "hostfs",
        "Flags": []
        }
    ],
}
```

Verification that this is actually broken can be performed using a simpler program that only tries to call `mknod(2)` and `ioctl(TUNSETIFF)`.

```c
#include <string.h>
#include <stdio.h>
#include <fcntl.h>

#include <sys/sysmacros.h>
#include <sys/stat.h>
#include <sys/ioctl.h>

#include <net/if.h>
#include <linux/if_tun.h>

#define CLEAR(x) memset(&(x), 0, sizeof(x))

static void
strncpynt(char *dest, const char *src, size_t maxlen)
{
    if (maxlen > 0)
    {
        strncpy(dest, src, maxlen-1);
        dest[maxlen - 1] = 0;
    }
}

int main() {
        int fd;
        int ret;
        struct ifreq ifr;
        const char *node = "/dev/net/tun";

        ret = mknod(node, S_IFCHR|0666, makedev(10, 200));

        printf("mknod returned: %d\n", ret);

        if ((fd = open(node, O_RDWR)) < 0) {
                printf("ERROR: Cannot open TUN/TAP dev %s\n", node);
                return 1;
        }

        CLEAR(ifr);
        ifr.ifr_flags = IFF_NO_PI;
        ifr.ifr_flags |= IFF_TUN;
        strncpynt(ifr.ifr_name, "ovpn", IFNAMSIZ);

        if ((ret = ioctl(fd, TUNSETIFF, (void *) &ifr)) < 0) {
                printf("ERROR ioctl(TUNSETIFF): %d\n", ret);
                return 1;
        }
        return 0;
}
```

A more basic run can be achieved with this, which will also show the broken calls:

```bash
gcc test.c -o app/test
myst mkcpio /app/ rootfs
myst exec-linux --strace rootfs /test
```
