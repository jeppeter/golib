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
	plaintext = []byte{}
	err = nil
	return
}

func AesEncCbc(plaintext []byte, key []byte, iv []byte, blksize int) (ciphertext []byte, err error) {
	ciphertext = []byte{}
	err = nil
	return
}
