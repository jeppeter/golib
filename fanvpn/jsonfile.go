package main

import (
	"aesext"
	"encoding/json"
	"encoding/pem"
	"crypto/rsa"
	"encoding/base64"
	"crypto/sha1"
	"crypto/x509"
	"crypto/rand"
	"dbgutil"
	"logutil"
	"reflect"
)


type VPNData struct {
	Key string 	`json:key`
	Iv string 	`json:iv`
	Data string `json:data`
}

type VPNParse struct {
	rsakey *rsa.PrivateKey
}

type VPNServer struct {
	Name string  `json:name`
	Flag string  `json:flag`
	Server string  `json:server`
	Port int `json:port`
}

type VPNConfig struct {
	Version int	  `json:version`
	Updated string    `json:updated`
	Nodes []VPNServer `json:nodes`
}

func NewVPNParse(privkey string) (retp *VPNParse,err error) {
	var block *pem.Block
	var keydata []byte
	var pkany any
	var leftbytes []byte
	retp = &VPNParse{}

	keydata = []byte(privkey)
	block, leftbytes = pem.Decode([]byte(privkey))
	if len(leftbytes) != len(keydata) {
		keydata = block.Bytes
	}

	pkany, err = x509.ParsePKCS8PrivateKey(keydata)
	if err != nil {
		err = dbgutil.FormatError("keydata not valid rsa %s", err.Error())
		return
	}
	switch pkany.(type) {
	case *rsa.PrivateKey:
		retp.rsakey = pkany.(*rsa.PrivateKey)
	default:
		err = dbgutil.FormatError("key is not rsakey type [%s]", reflect.TypeOf(pkany))
		return
	}

	err = nil
	return
}

func deocode_base64(inputs string) (retbuf []byte,err error) {
	var debuf []byte
	var n int
	debuf = make([]byte,base64.StdEncoding.DecodedLen(len(inputs)))
	n , err = base64.StdEncoding.Decode(debuf,[]byte(inputs))
	if err != nil {
		err = dbgutil.FormatError("decode [%s] error %s",inputs,err.Error())
		return
	}
	retbuf = make([]byte,n)
	for idx :=0 ;idx < n;idx += 1 {
		retbuf[idx] = debuf[idx]
	}
	err = nil
	return
}

func (retp *VPNParse)DecryptWithPrivateKey(ciphertext []byte) (plaintext []byte, err error) {
	hash := sha1.New()
	plaintext, err = rsa.DecryptOAEP(hash, rand.Reader, retp.rsakey, ciphertext, nil)
	if err != nil {
		err = dbgutil.FormatError("decrypt error %s", err.Error())
	}
	return
}


func (retp *VPNParse) GetConfig(jsons string) (retc *VPNConfig, err error) {
	var data *VPNData = &VPNData{}
	var keybuf,ivbuf,databuf []byte
	var keyplain []byte
	var cfgplain []byte
	retc = &VPNConfig{}
	err = json.Unmarshal([]byte(jsons),data)
	if err != nil {
		err = dbgutil.FormatError("parse error %s\n%s",err.Error(),jsons)
		return
	}

	/*now for the */
	keybuf, err = deocode_base64(data.Key)
	if err != nil {
		return
	}
	ivbuf ,err = deocode_base64(data.Iv)
	if err != nil {
		return
	}
	databuf, err = deocode_base64(data.Data)
	if err != nil {
		return
	}

	keyplain,err = retp.DecryptWithPrivateKey(keybuf)
	if err != nil {
		return
	}

	cfgplain,err = aesext.AesDecCbc(databuf,keyplain,ivbuf)
	if err != nil {
		return
	}
	logutil.Debug("cfgplain\n%s",string(cfgplain))
	err = json.Unmarshal(cfgplain,retc)
	if err != nil {
		err = dbgutil.FormatError("parse error %s\n%s",err.Error(),string(cfgplain))
		return
	}
	return

}