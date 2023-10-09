package utilities

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strconv"
	"time"

	"github.com/aead/ecdh"
	"github.com/dgrijalva/jwt-go"
)

// AuthCodeClaims represents the claims of an auth code
type AuthCodeClaims struct {
	ClientID    string `json:"cid"`
	UserID      int64  `json:"uid"`
	RedirectURI string `json:"ruri"`
	Scopes      string `json:"scp"`
	ExpiresAt   int64  `json:"exp"`
	jwt.StandardClaims
}

// GenerateRandomString generates random bytes and encodes them to a hex string
func GenerateRandomString(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("could not generate random bytes: %s", err)
	}

	return hex.EncodeToString(buf), nil
}

// GenerateAuthCode generates a JWT self-encoded auth code
//
//	Args:
//		secret: secret key
//		claims: auth code claims
func GenerateAuthCode(secret string, claims *AuthCodeClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("could not generate auth code: %s", err)
	}
	return tokenString, nil
}

// DecodeAuthCode decodes a JWT self-encoded auth code
//
//	Args:
//		secret: secret key
//		code: auth code
func DecodeAuthCode(secret, code string) (*AuthCodeClaims, error) {
	token, err := jwt.ParseWithClaims(code, &AuthCodeClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not decode auth code: %s", err)
	}
	claims, ok := token.Claims.(*AuthCodeClaims)
	if !ok {
		return nil, fmt.Errorf("could not decode auth code: %s", err)
	}
	if claims.ExpiresAt < jwt.TimeFunc().Unix() {
		return nil, fmt.Errorf("could not decode auth code: token expired")
	}
	return claims, nil
}

// GenerateUserToken generates a JWT token for a user
func GenerateUserToken(secret string, userID int64) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.TimeFunc().Add(time.Hour * 24 * 7).Unix(), // expires in 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("could not generate user token: %s", err)
	}
	return tokenString, nil
}

// VerifyUserToken verifies a JWT token for a user and returns the user ID
func VerifyUserToken(secret, token string) (int64, error) {
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("could not verify user token: %s", err)
	}
	if claims.ExpiresAt < jwt.TimeFunc().Unix() {
		return 0, fmt.Errorf("could not verify user token: token expired")
	}
	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("could not verify user token: %s", err)
	}
	return id, nil
}

// Algorithm represents an encryption algorithm
type Algorithm int

var (
	ECDSA_ALGO Algorithm = 1
	ECDH_ALGO  Algorithm = 2
)

// GenerateKeyPair generates a public/private key pair using the P-256 curve
func GenerateKeyPair(algo Algorithm) (crypto.PrivateKey, crypto.PublicKey, error) {
	switch algo {
	case ECDSA_ALGO:
		pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("could not generate key pair: %s", err)
		}
		return pri, &pri.PublicKey, nil
	case ECDH_ALGO:
		pri, pub, err := ecdh.KeyExchange.GenerateKey(ecdh.X25519(), rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("could not generate key pair: %s", err)
		}
		return pri, pub, nil
	default:
		return nil, nil, fmt.Errorf("unknown algorithm")
	}
}

// EncodePublicKey encodes a public key to a hex string
func EncodePublicKey(pub crypto.PublicKey) (string, error) {
	switch k := pub.(type) {
	case *ecdsa.PublicKey:
		code := hex.EncodeToString(k.X.Bytes()) + hex.EncodeToString(k.Y.Bytes())
		return code + fmt.Sprintf("_ALG%d", ECDSA_ALGO), nil
	case [32]uint8:
		return hex.EncodeToString(k[:]) + fmt.Sprintf("_ALG%d", ECDH_ALGO), nil
	default:
		return "", fmt.Errorf("unknown public key type: %s", reflect.TypeOf(pub))
	}
}

// DecodePublicKey decodes a hex string to a public key
func DecodePublicKey(encoded string) (crypto.PublicKey, error) {
	var (
		key []byte
		err error
	)
	algo, err := strconv.Atoi(encoded[len(encoded)-1:])
	if err != nil {
		return nil, fmt.Errorf("could not decode public key: %s", err)
	}
	key, err = hex.DecodeString(encoded[:len(encoded)-5])
	if err != nil {
		return nil, fmt.Errorf("could not decode public key: %s", err)
	}
	switch Algorithm(algo) {
	case ECDSA_ALGO:
		x := new(big.Int).SetBytes(key[:len(key)/2])
		y := new(big.Int).SetBytes(key[len(key)/2:])
		return &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		}, nil
	case ECDH_ALGO:
		var pub [32]uint8
		copy(pub[:], key)
		return pub, nil
	default:
		return nil, fmt.Errorf("unknown algorithm")
	}
}

// DeriveSharedSecret derives a shared secret from a public key and a private key
// using the ECDH algorithm
func DeriveSharedSecret(pri crypto.PrivateKey, pub crypto.PublicKey) []byte {
	return ecdh.KeyExchange.ComputeSecret(ecdh.X25519(), pri, pub)
}

// SymmetricEncrypt encrypts a message using a shared secret
func SymmetricEncrypt(message, ss []byte) (ciphertext []byte, err error) {
	// AES-GCM encryption
	key := make([]byte, 16)
	copy(key, ss)
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	// generate nonce and encrypt
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	ciphertext = aesgcm.Seal(nonce, nonce, message, nil)
	return
}

// SymmetricDecrypt decrypts a message using a shared secret
func SymmetricDecrypt(ciphertext, secret []byte) (message []byte, err error) {
	// AES-GCM decryption
	key := make([]byte, 16)
	copy(key, secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	// extract nonce and decrypt
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	message, err = aesgcm.Open(nil, nonce, ciphertext, nil)
	return
}

// SignMessage signs a base64 encoded message using a private key
func SignMessage(message string, pri crypto.PrivateKey) (string, error) {
	// hash message
	hash := sha256.Sum256([]byte(message))
	// sign
	r, s, err := ecdsa.Sign(rand.Reader, pri.(*ecdsa.PrivateKey), hash[:])
	if err != nil {
		return "", fmt.Errorf("could not sign message: %s", err)
	}
	// encode signature
	return hex.EncodeToString(r.Bytes()) + hex.EncodeToString(s.Bytes()), nil
}

// VerifySignature verifies a signature for a base64 encoded message using a public key
func VerifySignature(message, signature string, pub crypto.PublicKey) (bool, error) {
	// hash message
	hash := sha256.Sum256([]byte(message))
	// decode signature
	r, err := hex.DecodeString(signature[:len(signature)/2])
	if err != nil {
		return false, fmt.Errorf("could not verify signature: %s", err)
	}
	s, err := hex.DecodeString(signature[len(signature)/2:])
	if err != nil {
		return false, fmt.Errorf("could not verify signature: %s", err)
	}
	// verify
	return ecdsa.Verify(
		pub.(*ecdsa.PublicKey),
		hash[:],
		new(big.Int).SetBytes(r),
		new(big.Int).SetBytes(s),
	), nil
}
