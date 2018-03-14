package main

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
)

var (
	oidDataContentType          = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 7, 1})
	oidEncryptedDataContentType = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 7, 6})
	oidFriendlyName             = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 9, 20})
	oidLocalKeyID               = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 9, 21})
	oidMicrosoftCSPName         = asn1.ObjectIdentifier([]int{1, 3, 6, 1, 4, 1, 311, 17, 1})
	oidSHA256                   = asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 1})
	oidSHA1                     = asn1.ObjectIdentifier([]int{1, 3, 14, 3, 2, 26})
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

func (self *pfxPdu) verifySha1(passwd []byte) error {
	return nil
}

func (self *pfxPdu) verifySha256(passwd []byte) error {
	return nil
}

func (self *pfxPdu) Verify(passwd []byte) error {
	switch {
	case self.MacData.Mac.Algorithm.Algorithm.Equal(oidSHA1):
		return self.verifySha1(passwd)
	case self.MacData.Mac.Algorithm.Algorithm.Equal(oidSHA256):
		return self.verifySha256(passwd)
	}
	return fmt.Errorf("unknown type %v", self.MacData.Mac.Algorithm.Algorithm.String())
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
