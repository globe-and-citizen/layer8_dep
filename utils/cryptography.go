package utils

// TOMORROWS LABOUR CLEAN, ADD, REFACTOR TEST

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
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
	X       string   `json:"x,omitempty"`   // x coordinate as base64 URL encoded string.
	Y       string   `json:"y,omitempty"`   // y coordinate as base64 URL encoded string.
	D       string   `json:"d,omitempty"`   // d coordinate as base64 URL encoded string. Private keys only.
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
	privKey.D = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.D.Bytes())
	privKey.X = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.X.Bytes())
	privKey.Y = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.Y.Bytes())

	privKey_ecdsaPtr.ECDH()

	var pubKey JWK
	pubKey.Kid = "pub_" + id_str
	pubKey.Kty = "EC"
	pubKey.Crv = "P-256"
	pubKey.Kid = "pub_" + id_str
	pubKey.Key_ops = []string{}
	if keyUse == ECDSA {
		pubKey.Key_ops = append(pubKey.Key_ops, "verify")
	} else if keyUse == ECDH {
		pubKey.Key_ops = append(pubKey.Key_ops, "deriveKey")
	} else {
		return nil, nil, fmt.Errorf("Unrecognized keyUse. Must be 'ECDSA', or 'ECDH.'")
	}

	pubKey.X = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.X.Bytes())
	pubKey.Y = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.Y.Bytes())

	return &privKey, &pubKey, nil
}

func (jwk *JWK) ExportKeyAsGoType() (interface{}, error) {
	if jwk.Crv != "P-256" {
		return nil, fmt.Errorf("Cannot convert to *ecdh.PublicKey. Incorrect 'Crv' property.")
	}

	if jwk.Kty != "EC" {
		return nil, fmt.Errorf("Cannot convert to *ecdh.PublicKey. Incorrect 'Kty' property.")
	}

	// Step 1, create the 'common-to-all' base key. That is, an ecdsa.PublicKey
	pubKey := new(ecdsa.PublicKey)
	pubKey.Curve = elliptic.P256()

	bsX, err := base64.URLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, fmt.Errorf("Unable to interpret jwk.X coordinate as byte slice: %w", err)
	}
	pubKey.X = new(big.Int).SetBytes(bsX)

	bsY, err := base64.URLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, fmt.Errorf("Unable to interpret jwk.Y coordinate as byte slice: %w", err)
	}
	pubKey.Y = new(big.Int).SetBytes(bsY)

	// Step 2 decide if key is to be private

	var keyUsage KeyUse
	var privateFlag bool
	if slices.Contains(jwk.Key_ops, "deriveKey") {
		keyUsage = ECDH
	} else {
		keyUsage = ECDSA
	}
	if jwk.D != "" {
		privateFlag = true
	}

	privKey := new(ecdsa.PrivateKey)
	if privateFlag {
		privKey.PublicKey = *pubKey
		bsD, err := base64.URLEncoding.DecodeString(jwk.D)
		if err != nil {
			return nil, fmt.Errorf("Unable to interpret jwk.D coordinate as byte slice: %w", err)
		}
		privKey.D = new(big.Int).SetBytes(bsD)
	}

	if keyUsage == ECDSA && !privateFlag { // return an *ecdsa.PublicKey
		return pubKey, nil
	} else if keyUsage == ECDSA && privateFlag { // return an *ecdsa.PrivateKey
		return privKey, nil
	} else if keyUsage == ECDH && !privateFlag { // return an *ecdh.PublicKey
		publicKey, err := pubKey.ECDH()
		if err != nil {
			return nil, err
		}
		return publicKey, nil
	} else if keyUsage == ECDH && privateFlag { // return an *ecdsa.PrivateKey
		privateKey, err := privKey.ECDH()
		if err != nil {
			return nil, err
		}
		return privateKey, nil
	}
	return nil, fmt.Errorf("Unable to Export key. Unrecognized Key_ops.")
}

func (jwk *JWK) Equal(x crypto.PublicKey) bool {

	xx, ok := x.(*ecdsa.PublicKey)
	if !ok {
		return false
	}
	XBS, err := base64.URLEncoding.DecodeString(jwk.X)
	if err != nil {
		return false
	}

	XBI := new(big.Int).SetBytes(XBS)

	YBS, err := base64.URLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return false
	}
	YBI := new(big.Int).SetBytes(YBS)

	return bigIntEqual(XBI, xx.X) && bigIntEqual(YBI, xx.Y) &&
		jwk.Crv == xx.Curve.Params().Name

}

// func (*JWK) ExportECDSAKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error)
func (privateKey *JWK) GetECDHSharedSecret(publicKey *JWK) (*JWK, error) {
	// is public key public?
	if publicKey.D != "" {
		return nil, fmt.Errorf("Function takes a public JWK as argument. Private key detected.")
	}
	// is publis key for ECDH?
	if !slices.Contains(publicKey.Key_ops, "deriveKey") {
		return nil, fmt.Errorf("The public JWK passed in as argument must have 'deriveKey' as one of the key_ops")
	}
	// is private private?
	if privateKey.D == "" {
		return nil, fmt.Errorf("Function receiver must be private JWK. No private key detected.")
	}
	// is private for ECDH?
	if !slices.Contains(privateKey.Key_ops, "deriveKey") {
		return nil, fmt.Errorf("Function receiver must have 'deriveKey' as one of the Key_ops.")
	}
	// convert both to *ecdh.[private|public]Key
	privKey_unCasted, err := privateKey.ExportKeyAsGoType()
	if err != nil {
		return nil, fmt.Errorf("Unable to export private key as Go Type: %w", err)
	}

	var privKey *ecdh.PrivateKey
	if pk, ok := privKey_unCasted.(*ecdh.PrivateKey); ok {
		privKey = pk
	} else {
		return nil, fmt.Errorf("Returned private key could not be cast as *ecdh.PrivateKey")
	}

	pubKey_unCasted, err := publicKey.ExportKeyAsGoType()
	if err != nil {
		return nil, fmt.Errorf("Unable to export public key as Go Type: %w", err)
	}

	var pubKey *ecdh.PublicKey
	if pk, ok := pubKey_unCasted.(*ecdh.PublicKey); ok {
		pubKey = pk
	} else {
		return nil, fmt.Errorf("Returned public key could not be cast as *ecdh.PublicKey")
	}

	// Do ECDH
	ss, err := privKey.ECDH(pubKey)

	// Kid is derived from the private key's Kid

	symmetricJWK := &JWK{
		Kty:     "EC",
		Key_ops: []string{"encrypt", "decrypt"},
		Kid:     "shared_" + privateKey.Kid[4:],
		Crv:     privateKey.Crv,
		X:       base64.URLEncoding.EncodeToString(ss),
	}

	// convert the result to a JWK.
	return symmetricJWK, nil
}

func (ss *JWK) SymmetricEncrypt(data []byte) ([]byte, error) {
	if !slices.Contains(ss.Key_ops, "encrypt") {
		return nil, fmt.Errorf("Receiver Key_ops must include 'encrypt' ")
	}

	ssBS, err := base64.URLEncoding.DecodeString(ss.X)
	if err != nil {
		return nil, fmt.Errorf("Unable to interpret ss.X coordinate as byte slice: %w", err)
	}
	blockCipher, err := aes.NewCipher(ssBS)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %w", err)
	}
	aesgcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %w", err)
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %w", err)
	}

	cipherText := aesgcm.Seal(nonce, nonce, data, nil)

	return cipherText, nil
}

func (ss *JWK) SymmetricDecrypt(ciphertext []byte) ([]byte, error) {
	if !slices.Contains(ss.Key_ops, "decrypt") {
		return nil, fmt.Errorf("Receiver Key_ops must include 'decrypt' ")
	}

	ssBS, err := base64.URLEncoding.DecodeString(ss.X)
	if err != nil {
		return nil, fmt.Errorf("Unable to interpret ss.X coordinate as byte slice: %w", err)
	}
	blockCipher, err := aes.NewCipher(ssBS)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %w", err)
	}
	aesgcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %w", err)
	}
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("Symmetric encryption failed: %w", err)
	}
	return plaintext, nil
}

func (privateKey *JWK) SignWithKey(data []byte) ([]byte, error) {
	ecdsaPrivateKey, err := privateKey.ExportKeyAsGoType()
	if err != nil {
		return nil, fmt.Errorf("Unable to call ExportKeyAsGoType on function receiver. Error: %w", err)
	}
	var ecdsaKeyAsGoType *ecdsa.PrivateKey
	if privKey, ok := ecdsaPrivateKey.(*ecdsa.PrivateKey); ok {
		ecdsaKeyAsGoType = privKey
	} else {
		return nil, fmt.Errorf("Receiver, privateKey, of type %T could not be cast as compatible Go type, *ecdsa.PrivateKey.", privateKey)
	}

	hash32 := sha256.Sum256(data) //sha256 outputs a [32]byte that must be converted to []byte using the built in copy() method.
	var hash []byte
	copy(hash, hash32[:])
	signature, err := ecdsa.SignASN1(rand.Reader, ecdsaKeyAsGoType, hash)
	if err != nil {
		return nil, fmt.Errorf("Unable to sign data. Internal error: %w", err)
	}
	return signature, nil
}

func (publicKey *JWK) CheckAgainstASN1Signature(signature, data []byte) (bool, error) {
	if !slices.Contains(publicKey.Key_ops, "verify") {
		return false, fmt.Errorf("Check function receiver. *JWK Key_ops must include 'verify'")
	}

	ecdsaPublicKey, err := publicKey.ExportKeyAsGoType()
	if err != nil {
		return false, fmt.Errorf("Unable to call ExportKeyAsGoType on function receiver. Error: %w", err)
	}
	var ecdsaKeyAsGoType *ecdsa.PublicKey
	if pubKey, ok := ecdsaPublicKey.(*ecdsa.PublicKey); ok {
		ecdsaKeyAsGoType = pubKey
	} else {
		return false, fmt.Errorf("Receiver, publicKey, of type %T could not be cast as a compatible *ecdsa.PublicKey.", publicKey)
	}

	hash32 := sha256.Sum256(data) //sha256 outputs a [32]byte that must be converted to []byte using the built in copy() method.
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

// Util Utils?
func bigIntEqual(a, b *big.Int) bool {
	return subtle.ConstantTimeCompare(a.Bytes(), b.Bytes()) == 1
}
