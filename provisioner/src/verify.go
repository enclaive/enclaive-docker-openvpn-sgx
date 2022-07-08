package main

import (
	"log"
	"os"
	"strconv"
)

var (
	expectMrEnclave = os.Getenv(EnvMrEnclave)
	expectMrSigner  = os.Getenv(EnvMrSigner)
	expectISVProdId = os.Getenv(EnvIsvProdId)
	expectISVSVN    = os.Getenv(EnvIsvSvn)

	expectISVProdIdValue uint16
	expectISVSVNValue    uint16

	expectedMrEnclave = expectMrEnclave != "false" && expectMrEnclave != ""
	expectedMrSigner  = expectMrSigner != "false" && expectMrSigner != ""
	expectedISVProdId = expectISVProdId != "false" && expectISVProdId != ""
	expectedISVSVN    = expectISVSVN != "false" && expectISVSVN != ""
)

func init() {
	var err error
	var value uint64

	if expectedISVProdId {
		value, err = strconv.ParseUint(expectISVProdId, 10, 16)

		if err == nil {
			expectISVProdIdValue = uint16(value)
		}
	}

	if expectedISVSVN {
		value, err = strconv.ParseUint(expectISVSVN, 10, 16)

		if err == nil {
			expectISVSVNValue = uint16(value)
		}
	}
}

func checkProvidedValues(mrenclave, mrsigner string, isvProdId, isvSvn uint16) int {
	if expectedMrEnclave {
		log.Println("Checking MRENCLAVE", expectMrEnclave, "...")

		if mrenclave != expectMrEnclave {
			return McEnclave
		}
	}

	if expectedMrSigner {
		log.Println("Checking MRSIGNER against", expectMrSigner, "...")

		if mrsigner != expectMrSigner {
			return McSigner
		}
	}

	if expectedISVProdId {
		log.Println("Checking ISV Product ID against", expectISVProdIdValue, "...")

		if isvProdId != expectISVProdIdValue {
			return McIsvProdId
		}
	}

	if expectedISVSVN {
		log.Println("Checking ISV SVN against ", expectISVSVNValue, "...")

		if isvSvn != expectISVSVNValue {
			return McIsvSvn
		}
	}

	return McOk
}
