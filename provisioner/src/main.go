package main

//#include <stdlib.h>
//#include "secret_prov.h"
//#include "go_export.h"
//
// #cgo CFLAGS: -O3 -Wall -Werror -std=c11 -I/gramine/Pal/src/host/Linux-SGX/tools/ra-tls -L/usr/local/lib/x86_64-linux-gnu
// #cgo LDFLAGS: -lmbedcrypto_gramine -lmbedtls_gramine -lmbedx509_gramine -Wl,--no-as-needed -lsgx_urts -lsecret_prov_verify_dcap
//
/*
static int communicate_with_client_callback(struct ra_tls_ctx* ctx) {
    secret_provision_close(ctx);
    return 0;
}

int start(uint8_t *secret, size_t secret_size) {
    return secret_provision_start_server(
        secret, secret_size,
        "4433", "server.crt", "server.key",
        goVerifyMeasurementsCallback,
        communicate_with_client_callback
    );
}
*/
import "C"

import (
	"log"
	"os"
	"unsafe"
)

const (
	EnvKey       = "KEY_DEFAULT"
	EnvMrEnclave = "MRENCLAVE"
	EnvMrSigner  = "MRSIGNER"
	EnvIsvProdId = "ISV_PRODID"
	EnvIsvSvn    = "ISV_SVN"
)

func main() {
	aesKey, ok := os.LookupEnv(EnvKey)

	if !ok {
		log.Fatal("No key was set using KEY_DEFAULT environment variable")
	}

	if len(aesKey) != 32 {
		log.Fatal("Wrong key size, provided KEY_DEFAULT is", len(aesKey), "and should be 32")
	}

	// it is very important to include the null-byte here, otherwise no secret will be provisioned
	_secret := aesKey + "\x00"
	secret := C.CString(_secret)
	defer C.free(unsafe.Pointer(secret))

	log.Println("starting provisioning server on :4433")
	log.Println(C.start((*C.uchar)(unsafe.Pointer(secret)), C.ulong(len(_secret))))
}
