package utils
import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"math/big"
)

func VerifyGroupSignature(hash []byte, rText, sText string, pubKey *ecdsa.PublicKey) (bool, error) {
	r := new(big.Int)
	s := new(big.Int)

	_, ok := r.SetString(rText, 16)
	if !ok {
		return false, errors.New("invalid R value")
	}
	_, ok = s.SetString(sText, 16)
	if !ok {
		return false, errors.New("invalid S value")
	}

	hashed := sha256.Sum256(hash)
	valid := ecdsa.Verify(pubKey, hashed[:], r, s)
	return valid, nil
}