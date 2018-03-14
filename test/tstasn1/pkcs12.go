package main

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
)

var (
	oidDataContentType          = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 7, 1})
	oidEncryptedDataContentType = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 7, 6})

	oidFriendlyName     = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 9, 20})
	oidLocalKeyID       = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 9, 21})
	oidMicrosoftCSPName = asn1.ObjectIdentifier([]int{1, 3, 6, 1, 4, 1, 311, 17, 1})
)

type encryptedContentInfo struct {
	ContentType                asn1.ObjectIdentifier
	ContentEncryptionAlgorithm pkix.AlgorithmIdentifier
	EncryptedContent           []byte `asn1:"tag:0,optional"`
}

type encryptedData struct {
	Version              int
	EncryptedContentInfo encryptedContentInfo
}

type contentInfo struct {
	ContentType asn1.ObjectIdentifier
	Content     asn1.RawValue `asn1:"tag:0,explicit,optional"`
}

// from PKCS#7:
type digestInfo struct {
	Algorithm pkix.AlgorithmIdentifier
	Digest    []byte
}

type macData struct {
	Mac        digestInfo
	MacSalt    []byte
	Iterations int `asn1:"optional,default:1"`
}

type pfxPdu struct {
	Version  int
	AuthSafe contentInfo
	MacData  macData `asn1:"optional"`
}

func newPfxPdu() *pfxPdu {
	self := &pfxPdu{}
	return self
}

func (self *pfxPdu) Verify(passwd []byte) error {
	return nil
}

type Pkcs12 struct {
}

func newPkcs12() *Pkcs12 {
	self := &Pkcs12{}
	return self
}

func (self *Pkcs12) Decode(data []byte, passwd []byte) error {
	var rest []byte
	var err error
	var p *pfxPdu

	p = newPfxPdu()

	rest, err = asn1.Unmarshal(data, p)
	if err != nil {
		return err
	}

	if len(rest) > 0 {
		return fmt.Errorf("left %v", rest)
	}

	if p.Version != 3 {
		return fmt.Errorf("%d != 3 version", p.Version)
	}

	if !p.AuthSafe.ContentType.Equal(oidDataContentType) {
		return fmt.Errorf("not content type")
	}

	rest, err = asn1.Unmarshal(p.AuthSafe.Content.Bytes, &(p.AuthSafe.Content))
	if err != nil {
		return err
	}

	if len(rest) != 0 {
		return fmt.Errorf("left %v", rest)
	}

	if len(p.MacData.Mac.Algorithm.Algorithm) == 0 {
		return fmt.Errorf("no algorithm found")
	}

	err = p.Verify(passwd)
	if err != nil {
		return err
	}

	return nil
}
