package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
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
		return nil, nil, fmt.Errorf("unrecognized keyUse. Must be 'ECDSA', or 'ECDH.'")
	}

	privKey_ecdsaPtr, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privKey.D = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.D.Bytes())
	privKey.X = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.X.Bytes())
	privKey.Y = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.Y.Bytes())

	//privKey_ecdsaPtr.ECDH()

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
		return nil, nil, fmt.Errorf("unrecognized keyUse. Must be 'ECDSA', or 'ECDH.'")
	}

	pubKey.X = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.X.Bytes())
	pubKey.Y = base64.URLEncoding.EncodeToString(privKey_ecdsaPtr.Y.Bytes())

	return &privKey, &pubKey, nil
}

func (jwk *JWK) ExportKeyAsGoType() (interface{}, error) {
	if jwk.Crv != "P-256" {
		return nil, fmt.Errorf("cannot convert to *ecdh.PublicKey. Incorrect 'Crv' property")
	}

	if jwk.Kty != "EC" {
		return nil, fmt.Errorf("cannot convert to *ecdh.PublicKey. Incorrect 'Kty' property")
	}

	// Step 1, create the 'common-to-all' base key. That is, an ecdsa.PublicKey
	pubKey := new(ecdsa.PublicKey)
	pubKey.Curve = elliptic.P256()

	bsX, err := base64.URLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret jwk.X coordinate as byte slice: %s", err)
	}
	pubKey.X = new(big.Int).SetBytes(bsX)

	bsY, err := base64.URLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret jwk.Y coordinate as byte slice: %s", err)
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
			return nil, fmt.Errorf("unable to interpret jwk.D coordinate as byte slice: %s", err)
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
	return nil, fmt.Errorf("unable to Export key. Unrecognized Key_ops")
}

// Currently only supports checking against pk of type *ecdsa.Private/Public & *ecdh.Private/Public
func (jwk *JWK) Equal(pk interface{}) bool {
	switch key := pk.(type) {
	case *ecdsa.PrivateKey:
		if !slices.Contains(jwk.Key_ops, "sign") ||
			jwk.Kty != "EC" ||
			jwk.Crv != key.Params().Name ||
			jwk.X != base64.URLEncoding.EncodeToString(key.PublicKey.X.Bytes()) ||
			jwk.Y != base64.URLEncoding.EncodeToString(key.PublicKey.Y.Bytes()) ||
			jwk.D != base64.URLEncoding.EncodeToString(key.D.Bytes()) {
			return false
		}
	case *ecdsa.PublicKey:
		if !slices.Contains(jwk.Key_ops, "verify") ||
			jwk.Kty != "EC" ||
			jwk.Crv != key.Params().Name ||
			jwk.X != base64.URLEncoding.EncodeToString(key.X.Bytes()) ||
			jwk.Y != base64.URLEncoding.EncodeToString(key.Y.Bytes()) ||
			jwk.D != "" {
			return false
		}

	case *ecdh.PrivateKey:
		if !slices.Contains(jwk.Key_ops, "deriveKey") ||
			jwk.Kty != "EC" ||
			jwk.Crv != fmt.Sprint(key.Curve()) ||
			jwk.X != base64.URLEncoding.EncodeToString(key.PublicKey().Bytes()[1:33]) ||
			jwk.Y != base64.URLEncoding.EncodeToString(key.PublicKey().Bytes()[33:]) ||
			jwk.D != base64.URLEncoding.EncodeToString(key.Bytes()) {
			return false
		}

	case *ecdh.PublicKey:
		if !slices.Contains(jwk.Key_ops, "deriveKey") ||
			jwk.Kty != "EC" ||
			jwk.Crv != fmt.Sprint(key.Curve()) ||
			jwk.X != base64.URLEncoding.EncodeToString(key.Bytes()[1:33]) ||
			jwk.Y != base64.URLEncoding.EncodeToString(key.Bytes()[33:]) ||
			jwk.D != "" {
			return false
		}

	case *JWK:
		for _, val := range jwk.Key_ops {
			if !slices.Contains(key.Key_ops, val) {
				return false
			}
		}

		if jwk.Kty != key.Kty ||
			jwk.Kid != key.Kid ||
			jwk.Crv != key.Crv ||
			jwk.X != key.X ||
			jwk.Y != key.Y ||
			jwk.D != key.D {
			return false
		}

	default:
		//fmt.Println("ERROR: At this time only ECDSA & ECDH keys are supported.")
		return false
	}

	return true
}

func (jwk1 *JWK) EqualToJWK(jwk2 *JWK) bool {
	for _, val := range jwk1.Key_ops {
		if !slices.Contains(jwk2.Key_ops, val) {
			return false
		}
	}

	if jwk1.Kty != jwk2.Kty {
		return false
	}

	if jwk1.Kid != jwk2.Kid {
		return false
	}

	if jwk1.Crv != jwk2.Crv {
		return false
	}

	if jwk1.X != jwk2.X {
		return false
	}

	if jwk1.Y != jwk2.Y {
		return false
	}

	if jwk1.D != jwk2.D {
		return false
	}

	return true
}

// func (*JWK) ExportECDSAKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error)
func (privateKey *JWK) GetECDHSharedSecret(publicKey *JWK) (*JWK, error) {
	// is public key public?
	if publicKey.D != "" {
		return nil, fmt.Errorf("function takes a public JWK as argument. Private key detected")
	}
	// is publis key for ECDH?
	if !slices.Contains(publicKey.Key_ops, "deriveKey") {
		return nil, fmt.Errorf("the public JWK passed in as argument must have 'deriveKey' as one of the key_ops")
	}
	// is private private?
	if privateKey.D == "" {
		return nil, fmt.Errorf("function receiver must be private JWK. No private key detected")
	}
	// is private for ECDH?
	if !slices.Contains(privateKey.Key_ops, "deriveKey") {
		return nil, fmt.Errorf("function receiver must have 'deriveKey' as one of the Key_ops")
	}
	// convert both to *ecdh.[private|public]Key
	privKey_unCasted, err := privateKey.ExportKeyAsGoType()
	if err != nil {
		return nil, fmt.Errorf("unable to export private key as Go Type: %s", err)
	}

	var privKey *ecdh.PrivateKey
	if pk, ok := privKey_unCasted.(*ecdh.PrivateKey); ok {
		privKey = pk
	} else {
		return nil, fmt.Errorf("returned private key could not be cast as *ecdh.PrivateKey")
	}

	pubKey_unCasted, err := publicKey.ExportKeyAsGoType()
	if err != nil {
		return nil, fmt.Errorf("unable to export public key as Go Type: %s", err)
	}

	var pubKey *ecdh.PublicKey
	if pk, ok := pubKey_unCasted.(*ecdh.PublicKey); ok {
		pubKey = pk
	} else {
		return nil, fmt.Errorf("returned public key could not be cast as *ecdh.PublicKey")
	}

	// Do ECDH
	ss, err := privKey.ECDH(pubKey)
	if err != nil {
		return nil, fmt.Errorf("ECDH failed: %s", err)
	}

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
		return nil, fmt.Errorf("receiver Key_ops must include 'encrypt' ")
	}

	ssBS, err := base64.URLEncoding.DecodeString(ss.X)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret ss.X coordinate as byte slice: %s", err)
	}
	blockCipher, err := aes.NewCipher(ssBS)
	if err != nil {
		return nil, fmt.Errorf("symmetric encryption failed @ 1 : %s", err)
	}
	aesgcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("symmetric encryption failed @ 2: %s", err)
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("symmetric encryption failed @ 3: %s", err)
	}

	cipherText := aesgcm.Seal(nonce, nonce, data, nil)

	return cipherText, nil
}

func (ss *JWK) SymmetricDecrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, fmt.Errorf("receiver Key_ops must include 'decrypt' ")
	}

	if !slices.Contains(ss.Key_ops, "decrypt") {
		return nil, fmt.Errorf("receiver Key_ops must include 'decrypt' ")
	}

	ssBS, err := base64.URLEncoding.DecodeString(ss.X)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret ss.X coordinate as byte slice: %s", err)
	}
	blockCipher, err := aes.NewCipher(ssBS)
	if err != nil {
		return nil, fmt.Errorf("symmetric encryption failed @ 1: %s", err)
	}
	aesgcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("symmetric encryption failed @ 2: %s", err)
	}
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("symmetric encryption failed @ 3: %s", err)
	}
	return plaintext, nil
}

func (privateKey *JWK) SignWithKey(data []byte) ([]byte, error) {
	ecdsaPrivateKey, err := privateKey.ExportKeyAsGoType()
	if err != nil {
		return nil, fmt.Errorf("unable to call ExportKeyAsGoType on function receiver. Error: %s", err)
	}
	var ecdsaKeyAsGoType *ecdsa.PrivateKey
	if privKey, ok := ecdsaPrivateKey.(*ecdsa.PrivateKey); ok {
		ecdsaKeyAsGoType = privKey
	} else {
		return nil, fmt.Errorf("receiver, privateKey, of type %T could not be cast as compatible Go type, *ecdsa.PrivateKey", privateKey)
	}

	hash32 := sha256.Sum256(data) //sha256 outputs a [32]byte that must be converted to []byte using the built in copy() method.
	var hash []byte
	copy(hash, hash32[:])
	signature, err := ecdsa.SignASN1(rand.Reader, ecdsaKeyAsGoType, hash)
	if err != nil {
		return nil, fmt.Errorf("unable to sign data. Internal error: %s", err)
	}
	return signature, nil
}

func (publicKey *JWK) CheckAgainstASN1Signature(signature, data []byte) (bool, error) {
	if !slices.Contains(publicKey.Key_ops, "verify") {
		return false, fmt.Errorf("check function receiver. *JWK Key_ops must include 'verify'")
	}

	ecdsaPublicKey, err := publicKey.ExportKeyAsGoType()
	if err != nil {
		return false, fmt.Errorf("unable to call ExportKeyAsGoType on function receiver. Error: %s", err)
	}
	var ecdsaKeyAsGoType *ecdsa.PublicKey
	if pubKey, ok := ecdsaPublicKey.(*ecdsa.PublicKey); ok {
		ecdsaKeyAsGoType = pubKey
	} else {
		return false, fmt.Errorf("receiver, publicKey, of type %T could not be cast as a compatible *ecdsa.PublicKey", publicKey)
	}

	hash32 := sha256.Sum256(data) //sha256 outputs a [32]byte that must be converted to []byte using the built in copy() method.
	var hash []byte
	copy(hash, hash32[:])
	verified := ecdsa.VerifyASN1(ecdsaKeyAsGoType, hash, signature)
	if !verified {
		return false, fmt.Errorf("signature validation failed for this public key")
	}

	return verified, nil
}

func VerifyASN1Signature(JWK *JWK, signature, data []byte) (bool, error) {
	result, err := JWK.CheckAgainstASN1Signature(signature, data)
	if err != nil {
		return false, fmt.Errorf("unable to verify signature: %s", err)
	}
	return result, nil
}

func SignData(JWK *JWK, data []byte) ([]byte, error) {
	ASN1Signature, err := JWK.SignWithKey(data)
	if err != nil {
		return nil, fmt.Errorf("unable to sign: %s", err)
	}

	return ASN1Signature, nil
}

func (JWK *JWK) ExportAsBase64() (string, error) {
	marshalled, err := json.Marshal(JWK)
	if err != nil {
		return "", fmt.Errorf("failure to export JWK as Base64 %s", err.Error())
	}

	return base64.URLEncoding.EncodeToString(marshalled), nil
}

func B64ToJWK(userPubJWK string) (*JWK, error) {
	userPubJWK_BS, err := base64.URLEncoding.DecodeString(userPubJWK)
	if err != nil {
		return nil, fmt.Errorf("failure to decode userPubJWK: %s", err.Error())
	}
	userPubJWKConverted := &JWK{}
	err = json.Unmarshal(userPubJWK_BS, userPubJWKConverted)
	if err != nil {
		return nil, fmt.Errorf("failure to unmarshal userPubJWK: %s", err.Error())
	}
	return userPubJWKConverted, nil
}

// PRIVATE FUNCTIONS FOR PKG UTILS
func bigIntEqual(a, b *big.Int) bool {
	return subtle.ConstantTimeCompare(a.Bytes(), b.Bytes()) == 1
}
