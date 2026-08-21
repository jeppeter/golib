package aesext

import (
	"crypto/aes"
	"crypto/cipher"
	"dbgutil"
)

func AesEncEcb(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	var cblock cipher.Block
	var totaln int
	cblock, err = aes.NewCipher(key)
	if err != nil {
		err = dbgutil.FormatError("key cipher error %s", err.Error())
		return
	}
	if len(plaintext)%cblock.BlockSize() != 0 {
		err = dbgutil.FormatError("plaintext len %d %% %d != 0", len(plaintext), cblock.BlockSize())
		return
	}
	ciphertext = make([]byte, len(plaintext))
	for totaln = 0; totaln < len(plaintext); totaln += cblock.BlockSize() {
		cblock.Encrypt(ciphertext[totaln:], plaintext[totaln:])
	}
	err = nil
	return
}

func unpad_bytes(inbytes []byte,blksize int) (retbytes []byte,err error) {
	var leftsize int
	var idx int
	if (len(inbytes) % blksize ) != 0 {
		err = dbgutil.FormatError("input size [%d] %% %d != 0",len(inbytes),blksize)
		return
	}

	leftsize = int(inbytes[len(inbytes)-1])
	retbytes = make([]byte,len(inbytes) - leftsize)
	for idx = 0 ;idx<(len(inbytes) - leftsize) ; idx += 1 {
		retbytes[idx] = inbytes[idx]
	}

	for idx = len(retbytes);idx < len(inbytes) ; idx += 1 {
		if inbytes[idx] != byte(leftsize) {
			err = dbgutil.FormatError("[%d] [0x%02x] != 0x%02x", idx,inbytes[idx],leftsize)
			return
		}
	}
	err = nil
	return
}

func AesDecEcb(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	var cblock cipher.Block
	var totaln int
	cblock, err = aes.NewCipher(key)
	if err != nil {
		err = dbgutil.FormatError("key cipher error %s", err.Error())
		return
	}
	if len(ciphertext)%cblock.BlockSize() != 0 {
		err = dbgutil.FormatError("ciphertext len %d %% %d != 0", len(ciphertext), cblock.BlockSize())
		return
	}
	plaintext = make([]byte, len(ciphertext))
	for totaln = 0; totaln < len(ciphertext); totaln += cblock.BlockSize() {
		cblock.Decrypt(plaintext[totaln:], ciphertext[totaln:])
	}
	err = nil
	return
}

func AesDecCbc(ciphertext []byte, key []byte, iv []byte) (plaintext []byte, err error) {
	var obytes []byte
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	dec := cipher.NewCBCDecrypter(block, iv)
	obytes = make([]byte, len(ciphertext))
	dec.CryptBlocks(obytes, ciphertext)
	plaintext,err = unpad_bytes(obytes,aes.BlockSize)
	return
}

func pad_bytes(inbytes []byte,blksize int) (retbytes []byte) {
	var retlen int = len(inbytes)
	var i int
	retlen += (blksize - 1)
	retlen /= blksize
	retlen *= blksize
	if retlen == len(inbytes) {
		retlen += blksize
	}

	var padbyte byte
	padbyte = byte(retlen - len(inbytes))

	retbytes = make([]byte,retlen)
	for i =0 ; i < len(inbytes) ; i+= 1 {
		retbytes[i] = inbytes[i]
	}
	for i = len(inbytes);i < retlen;i+= 1 {
		retbytes[i] = padbyte
	}
	return 
}

func AesEncCbc(plaintext []byte, key []byte, iv []byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		err = dbgutil.FormatError("NewCipher for key error %s",err.Error())
		return
	}

	enc := cipher.NewCBCEncrypter(block, iv)
	nplain := pad_bytes(plaintext, aes.BlockSize)
	ciphertext = make([]byte, len(nplain))
	enc.CryptBlocks(ciphertext, nplain)
	err = nil
	return
}
