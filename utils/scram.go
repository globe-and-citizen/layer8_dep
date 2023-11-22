package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// GenerateNonce generates a random 16-byte nonce
func GenerateNonce() (string, error) {
	nonce := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(nonce), nil
}

// HmacSha256 function calculates the HMAC of a message
func HmacSha256(key, message []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// XOR function performs a bitwise XOR of two byte slices
func XOR(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, errors.New("unequal length")
	}

	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result, nil
}

// HI function computes SaltedPassword using PBKDF2
func HI(password, salt string, iterationCount int) []byte {
	normalizedPassword := []byte(password)
	saltedPassword := pbkdf2.Key(normalizedPassword, []byte(salt), iterationCount, sha256.Size, sha256.New)
	return saltedPassword
}

// H function computes StoredKey using SHA-256
func H(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// ConcatenateScramAttributes concatenates the given attributes as a message sequence
func ConcatenateScramAttributes(attributes map[string]string) string {
	if len(attributes) == 0 {
		return ""
	}

	attributeSequence := ""
	for k, v := range attributes {
		attributeSequence += "," + k + "=" + v
	}
	return attributeSequence[1:]
}

// ParseScramAttributes parses the given message sequence into a map of attributes
func ParseScramAttributes(attributeSequence string) map[string]string {
	attributes := make(map[string]string)
	if attributeSequence == "" {
		return attributes
	}

	attributePairs := strings.Split(attributeSequence, ",")
	for _, attributePair := range attributePairs {
		// string.Index is used here instead of strings.SplitN because the value may contain "=" for
		// base64 encoded data
		key := attributePair[:strings.Index(attributePair, "=")]
		value := attributePair[strings.Index(attributePair, "=")+1:]
		attributes[key] = value
	}
	return attributes
}

// ServerVerifyClientProof function verifies the client's proof by
// calculating the ClientSignature based on the stored key and
// exclusive-ORing the client proof and the calculated client signature
// to recover the ClientKey and then verifying the correctness of the
// client key by applying the hash function and comparing the result
// to the stored key
func ServerVerifyClientProof(username, combinedNonce string, storedKey, clientProof []byte) bool {
	// calculating the ClientSignature based on the stored key
	authMessage := "n=" + combinedNonce
	clientSignature := HmacSha256(storedKey, []byte(authMessage))

	// exclusive-ORing to recover the ClientKey
	clientKey, err := XOR(clientProof, clientSignature)
	if err != nil {
		return false
	}

	// verifing the correctness of the client key
	calcStoredKey := sha256.Sum256(clientKey)
	return hmac.Equal(calcStoredKey[:], storedKey)
}

// ServerGenerateServerSignature function generates the server's signature
func ServerGenerateServerSignature(username, combinedNonce string, serverKey []byte) string {
	authMessage := "n=" + combinedNonce
	serverSignature := HmacSha256(serverKey, []byte(authMessage))
	return base64.StdEncoding.EncodeToString(serverSignature)
}
