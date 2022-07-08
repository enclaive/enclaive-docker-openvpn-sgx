package main

//#include "go_export.h"
import "C"

import (
	"encoding/hex"
	"log"
	"unsafe"
)

const (
	McOk = -iota
	McError
	McEnclave
	McSigner
	McMrCombo
	McIsvProdId
	McIsvSvn
	McIsvCombo
)

var (
	mcErrors = map[int]string{
		McOk:        "Everything OK",
		McError:     "General error occured",
		McEnclave:   "Wrong MRENCLAVE value",
		McSigner:    "Wrong MRSIGNER value",
		McMrCombo:   "Invalid combination of MRSIGNER and MRENCLAVE value",
		McIsvProdId: "Wrong ISV product id value",
		McIsvSvn:    "Wrong ISV security version value",
		McIsvCombo:  "Invalid ISV combination",
	}
)

//export goVerifyMeasurementsCallback
func goVerifyMeasurementsCallback(_mrenclave, _mrsigner, _isv_prod_id, _isv_svn *C.cchar_t) C.int {
	if _mrenclave == nil || _mrsigner == nil || _isv_prod_id == nil || _isv_svn == nil {
		return C.int(McError)
	}

	mrenclave := hex.EncodeToString(C.GoBytes(unsafe.Pointer(_mrenclave), C.int(32)))
	mrsigner := hex.EncodeToString(C.GoBytes(unsafe.Pointer(_mrsigner), C.int(32)))
	isvProdId := (*uint16)(unsafe.Pointer(_isv_prod_id))
	isvSvn := (*uint16)(unsafe.Pointer(_isv_svn))

	if isvProdId == nil || isvSvn == nil {
		return C.int(McError)
	}

	ret := checkProvidedValues(mrenclave, mrsigner, *isvProdId, *isvSvn)

	log.Println("Received client measurements:")
	log.Println("\tMRENCLAVE:", mrenclave)
	log.Println("\tMRSIGNER:", mrsigner)
	log.Println("\tISV Product ID:", *isvProdId)
	log.Println("\tISV Security Version:", *isvSvn)
	log.Println("Check returned:", mcErrors[ret])

	return C.int(ret)
}
