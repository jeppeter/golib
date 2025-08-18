package main

import (
	"crypto/x509"
)

func x509req_CheckSignature(c *x509.CertificateRequest) error {
	return checkSignature(c.SignatureAlgorithm, c.RawTBSCertificateRequest, c.Signature, c.PublicKey, true)
}
