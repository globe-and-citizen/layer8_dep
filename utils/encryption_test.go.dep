package utilities

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	t.Run("generate random bytes", func(t *testing.T) {
		str, err := GenerateRandomString(32)
		assert.NoError(t, err)
		buf, err := hex.DecodeString(str)
		assert.NoError(t, err)
		assert.Equal(t, 32, len(buf))
	})
}

func TestGenerateAuthCode_DecodeAuthCode(t *testing.T) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	assert.NoError(t, err)

	tc := []struct {
		claims    *AuthCodeClaims
		expectErr bool
	}{
		{
			claims: &AuthCodeClaims{
				ClientID:    "client-id-1",
				UserID:      1,
				Scopes:      "scope1,scope2",
				RedirectURI: "https://example.com/callback/1",
				ExpiresAt:   time.Now().Add(time.Minute * 5).Unix(),
			},
			expectErr: false,
		},
		{
			claims: &AuthCodeClaims{
				ClientID:    "client-id-2",
				UserID:      2,
				Scopes:      "scope1,scope2",
				RedirectURI: "https://example.com/callback/2",
				ExpiresAt:   time.Now().Add(time.Minute * -5).Unix(),
			},
			expectErr: true,
		},
	}

	for _, c := range tc {
		code, err := GenerateAuthCode(string(secret), c.claims)
		assert.NoError(t, err)
		assert.NotEmpty(t, code)

		decoded, err := DecodeAuthCode(string(secret), code)
		if c.expectErr {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, c.claims.ClientID, decoded.ClientID)
		assert.Equal(t, c.claims.RedirectURI, decoded.RedirectURI)
		assert.Equal(t, c.claims.UserID, decoded.UserID)
		assert.Equal(t, c.claims.ExpiresAt, decoded.ExpiresAt)
	}
}

func TestGenerateUserToken_VerifyserToken(t *testing.T) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	assert.NoError(t, err)

	tc := []*AuthCodeClaims{
		{
			UserID:    1,
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
		{
			UserID:    2,
			ExpiresAt: time.Now().Add(time.Minute * -5).Unix(),
		},
	}

	for _, c := range tc {
		token, err := GenerateUserToken(string(secret), c.UserID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		id, err := VerifyUserToken(string(secret), token)
		assert.NoError(t, err)
		assert.Equal(t, c.UserID, id)
	}
}

func TestEncodePublicKey_DecodePublicKey(t *testing.T) {
	// key generation
	_, pub, err := GenerateKeyPair(ECDH_ALGO)
	if err != nil {
		t.Errorf("GenerateKeyPair() = %v", err)
	}

	encoded, err := EncodePublicKey(pub)
	assert.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := DecodePublicKey(encoded)
	assert.NoError(t, err)
	assert.Equal(t, pub, decoded)
}

func TestDeriveSharedSecret_SymmetricEncrypt_SymmetricDecrypt(t *testing.T) {
	// key generation
	pri1, pub1, err := GenerateKeyPair(ECDH_ALGO)
	if err != nil {
		t.Errorf("GenerateKeyPair() = %v", err)
	}
	pri2, pub2, err := GenerateKeyPair(ECDH_ALGO)
	if err != nil {
		t.Errorf("GenerateKeyPair() = %v", err)
	}

	// derive shared secret
	ss1 := DeriveSharedSecret(pri2, pub1)
	ss2 := DeriveSharedSecret(pri1, pub2)

	// encrypt/decrypt
	plaintext := []byte("Hello, World!")
	ciphertext, err := SymmetricEncrypt(plaintext, ss1)
	if err != nil {
		t.Errorf("SymmetricEncrypt() = %v", err)
	}
	decrypted, err := SymmetricDecrypt(ciphertext, ss2)
	if err != nil {
		t.Errorf("SymmetricDecrypt() = %v", err)
	}
	if string(decrypted) != string(plaintext) {
		t.Errorf("SymmetricDecrypt() = %s, want %s", string(decrypted), string(plaintext))
	}
}

func TestSignMessage_VerifySignature(t *testing.T) {
	// key pair
	pri, pub, err := GenerateKeyPair(ECDSA_ALGO)
	if err != nil {
		t.Errorf("GenerateKeyPair() = %v", err)
	}
	// message
	var message []byte
	_, err = rand.Read(message)
	if err != nil {
		t.Errorf("rand.Read() = %v", err)
	}
	// sign
	b64msg := base64.RawURLEncoding.EncodeToString(message)
	signature, err := SignMessage(b64msg, pri)
	if err != nil {
		t.Errorf("SignMessage() = %v", err)
	}
	// verify
	valid, err := VerifySignature(b64msg, signature, pub)
	if err != nil {
		t.Errorf("VerifySignature() = %v", err)
	}
	if !valid {
		t.Errorf("VerifySignature() = %v, want %v", valid, true)
	}
}
