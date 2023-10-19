package utils

// TOMORROWS LABOUR CLEAN, ADD, REFACTOR TEST

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
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

// This will give you a private and public key pair both as *JWK structs.
// The Kid (i.e., "key-id") is equivalent shared but for the prefix "priv_" or "pub_"
// appended to a random base64 URL encoded string.The Crv parameter is hardcoded
// for now during key creation as are the key_ops.

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
		privKey.Key_ops = append(privKey.Key_ops, "sign")
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
	if len(jwk.D) != 0 {
		// if jwk.D != nil {
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

func (privateKey *JWK) GetECDHSharedSecret(publicKey *JWK) (*JWK, error) {
	// is public key public?
	if len(publicKey.D) != 0 {
		return nil, fmt.Errorf("Function takes a public JWK as argument. Private key detected.")
	}
	// is publis key for ECDH?
	if !slices.Contains(publicKey.Key_ops, "deriveKey") {
		return nil, fmt.Errorf("The public JWK passed in as argument must have 'deriveKey' as one of the key_opts")
	}
	// is private private?
	if len(privateKey.D) == 0 {
		return nil, fmt.Errorf("Function receiver must be private JWK. No private key detected.")
	}
	// is private for ECDH?
	if !slices.Contains(privateKey.Key_ops, "deriveKey") {
		return nil, fmt.Errorf("Function receiver must be for ECDH. Key_ops does not contain 'deriveKey' as an option.")
	}
	// convert both to *ecdh.[private|public]Key
	privKey_unCasted, err := privateKey.ExportKeyAsGoType()
	if err != nil {
		fmt.Println(err.Error())
	}

	var privKey *ecdh.PrivateKey
	if pk, ok := privKey_unCasted.(*ecdh.PrivateKey); ok {
		privKey = pk
	}

	pubKey_unCasted, err := publicKey.ExportKeyAsGoType()
	if err != nil {
		fmt.Println(err.Error())
	}

	var pubKey *ecdh.PublicKey
	if pk, ok := pubKey_unCasted.(*ecdh.PublicKey); ok {
		pubKey = pk
	}
	// to the ECDH
	ss, err := privKey.ECDH(pubKey)

	// Kid is derived from the private key's Kid

	symmetricJWK := &JWK{
		Kty:     "EC",
		Key_ops: []string{"encrypt", "decrypt"},
		Kid:     "shared_" + privateKey.Kid[4:],
		Crv:     privateKey.Crv,
		X:       ss,
	}

	// convert the result to a JWK.
	return symmetricJWK, err
}

func (ss *JWK) SymmetricEncrypt(data []byte) ([]byte, error) {
	if !slices.Contains(ss.Key_ops, "encrypt") {
		return nil, fmt.Errorf("Receiver Key_ops must include 'encrypt' ")
	}

	blockCipher, err := aes.NewCipher(ss.X)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %s, %d", err.Error())
	}
	aesgcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %s", err.Error())
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %s", err.Error())
	}

	cipherText := aesgcm.Seal(nonce, nonce, data, nil)

	return cipherText, nil
}

func (ss *JWK) SymmetricDecrypt(ciphertext []byte) ([]byte, error) {
	if !slices.Contains(ss.Key_ops, "decrypt") {
		return nil, fmt.Errorf("Receiver Key_ops must include 'decrypt' ")
	}

	blockCipher, err := aes.NewCipher(ss.X)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %s", err.Error())
	}
	aesgcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %s", err.Error())
	}
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %s", err.Error())
	}
	return plaintext, nil
}

func (privateKey *JWK) SignWithKey(data []byte) ([]byte, error) {
	ecdsaPrivateKey, err := privateKey.ExportKeyAsGoType()
	if err != nil {
		return nil, fmt.Errorf("Unable to call ExportKeyAsGoType on function receiver. Error: %s", err.Error())
	}
	var ecdsaKeyAsGoType *ecdsa.PrivateKey
	if privKey, ok := ecdsaPrivateKey.(*ecdsa.PrivateKey); ok {
		ecdsaKeyAsGoType = privKey
	} else {
		return nil, fmt.Errorf("Receiver, privateKey, of type %T could not be cast as compatible Go type, *ecdsa.PrivateKey.", privateKey)
	}

	hash32 := sha256.Sum256(data) //sha256 outputs a [32]byte that must be converted to []byte
	var hash []byte
	copy(hash, hash32[:])
	signature, err := ecdsa.SignASN1(rand.Reader, ecdsaKeyAsGoType, hash)
	if err != nil {
		return nil, fmt.Errorf("Unable to sign data. Internal error: %s", err.Error())
	}
	return signature, nil
}

func (publicKey *JWK) CheckAgainstASN1Signature(signature, data []byte) (bool, error) {
	if !slices.Contains(publicKey.Key_ops, "verify") {
		return false, fmt.Errorf("Check function receiver. Must be a *JWK with Key_ops including 'verify'")
	}

	ecdsPublicKey, err := publicKey.ExportKeyAsGoType()
	if err != nil {
		return false, fmt.Errorf("Unable to call ExportKeyAsGoType on function receiver. Error: %s", err.Error())
	}
	var ecdsaKeyAsGoType *ecdsa.PublicKey
	if pubKey, ok := ecdsPublicKey.(*ecdsa.PublicKey); ok {
		ecdsaKeyAsGoType = pubKey
	} else {
		return false, fmt.Errorf("Receiver, publicKey, of type %T could not be cast as a compatible *ecdsa.PublicKey.", publicKey)
	}

	hash32 := sha256.Sum256(data)
	var hash []byte
	copy(hash, hash32[:])
	verified := ecdsa.VerifyASN1(ecdsaKeyAsGoType, hash, signature)
	if verified != true {
		return false, fmt.Errorf("Signature validation failed for this public key")
	}

	return verified, nil
}

func VerifyASN1Signature(JWK *JWK, signature, data []byte) (bool, error) {
	result, err := JWK.CheckAgainstASN1Signature(signature, data)
	if err != nil {
		return false, fmt.Errorf("Unable to verify signature: %w", err)
	}
	return result, nil
}

func SignData(JWK *JWK, data []byte) ([]byte, error) {
	ASN1Signature, err := JWK.SignWithKey(data)
	if err != nil {
		return nil, fmt.Errorf("Unable to sign: %w", err)
	}

	return ASN1Signature, nil
}
