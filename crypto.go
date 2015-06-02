package gitcha

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func genSecret(name string) (string, error) {
	if n := len(name); n > keyLen {
		return "", ErrNameTooLong
	}

	kb := make([]byte, keyLen)
	copy(kb, []byte(name))
	if _, err := rand.Read(kb); err != nil {
		return "", err
	}
	sec := base64.URLEncoding.EncodeToString(kb)[:secLen]

	return sec, nil
}

func genID(sec string) (string, error) {
	if n := len(sec); n != secLen {
		return "", fmt.Errorf("expected secret length %d, got %d", keyLen, n)
	}

	cip, err := aes.NewCipher([]byte(sec[:32]))
	if err != nil {
		return "", err
	}

	idb := make([]byte, cip.BlockSize())
	kb := make([]byte, cip.BlockSize())
	if _, err = rand.Read(kb); err != nil {
		return "", err
	}
	cip.Encrypt(idb, kb)
	id := base64.URLEncoding.EncodeToString(idb)

	return id, nil
}
