package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"slices"
)

type KeyUse int

var (
	ECDSA KeyUse = 1
	ECDH  KeyUse = 2
)

type JWK struct {
	Key_ops []string `json:"use,omitempty"` // ["sign", "verify", "encrypt", "decrypt", "wrapKey", "unwrapKey", "deriveKey", "deriveBits"]
	Kty     string   `json:"kty,omitempty"` // "EC", "RSA"
	Kid     string   `json:"kid,omitempty"` // Key ID
	Crv     string   `json:"crv,omitempty"` // "P-256"
	X       []byte   `json:"x,omitempty"`   // x coordinate
	Y       []byte   `json:"y,omitempty"`   // y coordinate
	D       []byte   `json:"d,omitempty"`   // private keys only
}

// This will give you a private and public key pair as *JWK structs.
// The Kid is shared but for the prefix "priv_" or "pub_"
// The Crv parameter is hardcoded for now as are the key opts.
// To get the *ecdh.PrivateKey/PublicKey or *ecdsa.PublicKey/*ecdsa.PrivateKey
// use *JWK.ExportKeyAsGoType() and then type assert the resulting key

func GenerateKeyPair(keyUse KeyUse) (*JWK, *JWK, error) {
	id := make([]byte, 16)
	rand.Read(id)
	id_str := base64.URLEncoding.EncodeToString(id)

	var privKey JWK
	privKey.Kty = "EC"
	privKey.Crv = "P-256"
	privKey.Kid = "priv_" + id_str
	privKey.Key_ops = []string{}
	if keyUse == ECDSA {
		privKey.Key_ops = append(privKey.Key_ops, "sign", "decrypt")
	} else if keyUse == ECDH {
		privKey.Key_ops = append(privKey.Key_ops, "deriveKey")
	} else {
		return nil, nil, fmt.Errorf("Unrecognized keyUse. Must be 'ECDSA', or 'ECDH.'")
	}

	privKey_ecdsaPtr, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privKey.D = privKey_ecdsaPtr.D.Bytes()
	privKey.X = privKey_ecdsaPtr.X.Bytes()
	privKey.Y = privKey_ecdsaPtr.Y.Bytes()

	var pubKey JWK
	pubKey.Kid = "pub_" + id_str
	pubKey.Kty = "EC"
	pubKey.Crv = "P-256"
	pubKey.Kid = "pub_" + id_str
	pubKey.Key_ops = []string{}
	if keyUse == ECDSA {
		pubKey.Key_ops = append(pubKey.Key_ops, "verify", "encrypt")
	} else if keyUse == ECDH {
		pubKey.Key_ops = append(pubKey.Key_ops, "deriveKey")
	} else {
		return nil, nil, fmt.Errorf("Unrecognized keyUse. Must be 'ECDSA', or 'ECDH.'")
	}

	pubKey.X = privKey_ecdsaPtr.X.Bytes()
	pubKey.Y = privKey_ecdsaPtr.Y.Bytes()

	return &privKey, &pubKey, nil
}

func (jwk *JWK) ExportKeyAsGoType() (interface{}, error) {
	if jwk.Crv != "P-256" {
		return nil, fmt.Errorf("Cannot convert to *ecdh.PublicKey. Incorrect 'Crv' property.")
	}

	if jwk.Kty != "EC" {
		return nil, fmt.Errorf("Cannot convert to *ecdh.PublicKey. Incorrect 'Kty' property.")
	}

	// Step 1, create the ecdsa.PublicKey
	pubKey := new(ecdsa.PublicKey)
	pubKey.Curve = elliptic.P256()
	pubKey.X = new(big.Int).SetBytes(jwk.X)
	pubKey.Y = new(big.Int).SetBytes(jwk.Y)

	// Step 2 decide if private
	privKey := new(ecdsa.PrivateKey)

	// now I have a pubKey and a privKey
	// to export, I can convert or not to ECDH()
	var keyUsage KeyUse
	var privateFlag bool
	if slices.Contains(jwk.Key_ops, "deriveKey") {
		keyUsage = ECDH
	} else {
		keyUsage = ECDSA
	}
	if jwk.D != nil {
		privateFlag = true
	}

	if privateFlag {
		privKey.PublicKey = *pubKey
		privKey.D = new(big.Int).SetBytes(jwk.D)
	}

	if keyUsage == ECDSA && !privateFlag {
		return pubKey, nil
	} else if keyUsage == ECDSA && privateFlag {
		return privKey, nil
	} else if keyUsage == ECDH && !privateFlag {
		key, err := pubKey.ECDH()
		if err != nil {
			return nil, err
		}
		return key, nil
	} else if keyUsage == ECDH && privateFlag {
		key, err := privKey.ECDH()
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	return nil, fmt.Errorf("Unable to Export key. Unrecognized Key_opts.")
}
