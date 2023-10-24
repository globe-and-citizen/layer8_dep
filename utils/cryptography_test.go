package utils

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"slices"
	"testing"
	"time"
)

// SHOULD YOU IMPLEMENT A CONVERSION FUNCTION?
func Test_One(t *testing.T) {
	//Everything it Seems I need to do the conversion?

	privKey_ECDH, pubKey_ECDH, err := GenerateKeyPair(ECDH)
	if err != nil {
		t.Errorf(err.Error())
	}
	pKeyInt, _ := privKey_ECDH.ExportKeyAsGoType()
	privKeyCasted, _ := pKeyInt.(*ecdh.PrivateKey)
	D_bytes := privKeyCasted.Bytes()
	D_bigInt := new(big.Int).SetBytes(D_bytes)

	coorD_BS, _ := base64.URLEncoding.DecodeString(privKey_ECDH.D)
	coorD_bigInt := new(big.Int).SetBytes(coorD_BS)
	if coorD_bigInt.Cmp(D_bigInt) != 0 {
		t.Log("Values of the D coordinate do not match after exporting and conversion.")
	}

	t.Log("BSD1: ", D_bytes)
	t.Log("BSD2: ", coorD_BS)
	pubKey_BS := privKeyCasted.PublicKey().Bytes()

	/* It means 'this elliptic curve point is specified in uncompressed format', that is, the x and y
	*  coordinates are given explicitly.The other alternatives are:
	*  02, which means 'this is a compressed point, where we give the x coordinate explicitly; of the two possible y coordinates that are compatible with that x coordinate, select the one with the 0 lsbit.
	*  03, which is the same, except you select the y coordinate with a 1 lsbit. The compressed formats are about half as long (saving space), but requires more computation (a modular square root) to use if you perform an operation that requires the y coordinate.
	* This is why the first byte is removed.
	 */
	t.Log("BSX1: ", pubKey_BS[1:33])
	t.Log("BSY1: ", pubKey_BS[33:])

	orginalX_BS, _ := base64.URLEncoding.DecodeString(pubKey_ECDH.X)
	t.Log("BSX2: ", orginalX_BS)

	orginalY_BS, _ := base64.URLEncoding.DecodeString(pubKey_ECDH.Y)
	t.Log("BSY2: ", orginalY_BS)

	fmt.Println(fmt.Sprint(privKeyCasted.Curve())) // now I can get the string.
}

func Test_GenerateKeyPair(t *testing.T) {
	privKey_ECDSA, pubKey_ECDSA, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !slices.Contains(privKey_ECDSA.Key_ops, "sign") {
		t.Errorf("Private ECDSA Keys must contain as a key option, 'sign'")
	}

	if privKey_ECDSA.Kty != "EC" {
		t.Errorf("ECDSA keys must have the kty parameter set to 'EC'")
	}

	if privKey_ECDSA.Crv != "P-256" {
		t.Errorf("Something went wrong. All keys are currently suppose to use the curve 'P-256.'")
	}

	if privKey_ECDSA.X == "" {
		t.Errorf("Private keys must have the X coordinate set.")
	}

	if privKey_ECDSA.Y == "" {
		t.Errorf("Private keys must have the Y coordinate set.")
	}

	if privKey_ECDSA.D == "" {
		t.Errorf("Private keys must have the D coordinate set.")
	}

	if privKey_ECDSA.Kid[4:] != pubKey_ECDSA.Kid[3:] {
		t.Log("Private Key ID: ", privKey_ECDSA.Kid[3:])
		t.Log("Public Key ID: ", pubKey_ECDSA.Kid[2:])
		t.Errorf("ECDSA key ids are a mismatch. Key ids for the key pair should match but for the prefix")
	}

	if pubKey_ECDSA.X == "" {
		t.Errorf("Public keys must have the Y coordinate set.")
	}

	if pubKey_ECDSA.Y == "" {
		t.Errorf("Public keys must have the Y coordinate set.")
	}

	if pubKey_ECDSA.D != "" {
		t.Errorf("Only the private key of a key pair must have the D coordinate set.")
	}

	privK, err := privKey_ECDSA.ExportKeyAsGoType()
	if err != nil {
		t.Errorf(err.Error())
	}

	if _, ok := privK.(*ecdsa.PrivateKey); !ok {
		t.Errorf("A generated ECDSA private key must be convertible to the *ecdsa.PrivateKey type in Go.")
	}

	pubKey, err := pubKey_ECDSA.ExportKeyAsGoType()
	if err != nil {
		t.Errorf("pubKey_ECSA coudl not be exported as Go type: %s \n", err.Error())
	}

	if _, ok := pubKey.(*ecdsa.PublicKey); !ok {
		t.Errorf("A generated ECDSA public key must be convertible to the *ecdsa.PublicKey type in Go.")
	}

}

func Test_Equal(t *testing.T) {
	// Check creation, conversion and equivalence of a private/public ECDSA key
	privJWK_ecdsa, pubJWK_ecdsa, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Errorf(err.Error())
	}

	privJWKCopy := &JWK{
		Key_ops: make([]string, len(privJWK_ecdsa.Key_ops)),
		Kty:     privJWK_ecdsa.Kty,
		Kid:     privJWK_ecdsa.Kid,
		Crv:     privJWK_ecdsa.Crv,
		X:       privJWK_ecdsa.X,
		Y:       privJWK_ecdsa.Y,
		D:       privJWK_ecdsa.D,
	}

	for i, val := range privJWK_ecdsa.Key_ops {
		privJWKCopy.Key_ops[i] = val
	}

	if !privJWK_ecdsa.Equal(privJWKCopy) {
		t.Errorf("Exported ECDSA private key doesn't equal itself after copying.")
	}

	pubJWK_Copy := &JWK{
		Key_ops: make([]string, len(pubJWK_ecdsa.Key_ops)),
		Kty:     pubJWK_ecdsa.Kty,
		Kid:     pubJWK_ecdsa.Kid,
		Crv:     pubJWK_ecdsa.Crv,
		X:       pubJWK_ecdsa.X,
		Y:       pubJWK_ecdsa.Y,
		D:       pubJWK_ecdsa.D,
	}

	for i, val := range pubJWK_ecdsa.Key_ops {
		pubJWK_Copy.Key_ops[i] = val
	}

	if !pubJWK_ecdsa.Equal(pubJWK_Copy) {
		t.Errorf("Exported ECDSA public key doesn't equal itself after copying.")
	}

	// Check creation, conversion and equivalence of a private/public ECDH key
	privJWK_ecdh, pubJWK_ecdh, err := GenerateKeyPair(ECDH)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Copy & then test the private key
	privJWKCopy = &JWK{
		Key_ops: make([]string, len(privJWK_ecdh.Key_ops)),
		Kty:     privJWK_ecdh.Kty,
		Kid:     privJWK_ecdh.Kid,
		Crv:     privJWK_ecdh.Crv,
		X:       privJWK_ecdh.X,
		Y:       privJWK_ecdh.Y,
		D:       privJWK_ecdh.D,
	}

	for i, val := range privJWK_ecdh.Key_ops {
		privJWKCopy.Key_ops[i] = val
	}

	if !privJWK_ecdh.Equal(privJWKCopy) {
		t.Errorf("Exported ECDH private key doesn't equal itself after copying.")
	}

	// Copy & then test the public key
	pubJWK_Copy = &JWK{
		Key_ops: make([]string, len(pubJWK_ecdh.Key_ops)),
		Kty:     pubJWK_ecdh.Kty,
		Kid:     pubJWK_ecdh.Kid,
		Crv:     pubJWK_ecdh.Crv,
		X:       pubJWK_ecdh.X,
		Y:       pubJWK_ecdh.Y,
		D:       pubJWK_ecdh.D,
	}

	for i, val := range pubJWK_ecdh.Key_ops {
		pubJWK_Copy.Key_ops[i] = val
	}

	if !pubJWK_ecdh.Equal(pubJWK_Copy) {
		t.Errorf("Exported ECDH public key doesn't equal itself after copying.")
	}
}

func Test_SignWithKey(t *testing.T) {
	const NUMBER_OF_CASES = 10

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < NUMBER_OF_CASES; i++ {
		length := r1.Intn(2000)
		var randomByteSlice []byte
		randomBuffer := bytes.NewBuffer(randomByteSlice)
		for i := 0; i < length; i++ {
			randomData := r1.Intn(255)
			err := binary.Write(randomBuffer, binary.BigEndian, uint8(randomData))
			if err != nil {
				t.Fatal(err.Error())
			}
		}
		randomDataToSign := randomBuffer.Bytes()
		privJWK_ecdsa1, pubJWK_ecdsa1, _ := GenerateKeyPair(ECDSA)

		signature, err := privJWK_ecdsa1.SignWithKey(randomDataToSign)
		if err != nil {
			t.Error(err.Error())
		}
		verified, err := pubJWK_ecdsa1.CheckAgainstASN1Signature(signature, randomDataToSign)
		if err != nil {
			t.Error(err.Error())
		}
		if !verified {
			t.Error("Signature verification failed")
		}
	}
}

// TEST ECDSA
func Test_SymmetricEncryption(t *testing.T) {
	// Generate first key pair
	privJWK, pubJWK, err := GenerateKeyPair(ECDH)
	if err != nil {
		t.Error(err.Error())
	}

	// Generate second Key pair
	privJWK2, pubJWK2, err := GenerateKeyPair(ECDH)
	if err != nil {
		t.Error(err.Error())
	}

	// Derive the shared secret
	ssJWK1, err := privJWK.GetECDHSharedSecret(pubJWK2)
	if err != nil {
		t.Error(err.Error())
	}
	ssJWK2, err := privJWK2.GetECDHSharedSecret(pubJWK)
	if err != nil {
		t.Error(err.Error())
	}

	dataPoints := [][]byte{
		[]byte("Byte Slice 1 To Encrypt"),
		[]byte("One, Two, three"),
		[]byte(""),
		[]byte("                           "),
	}

	for _, data := range dataPoints {
		ct, err := ssJWK1.SymmetricEncrypt(data)
		if err != nil {
			t.Error(err.Error())
		}

		pt, err := ssJWK2.SymmetricDecrypt(ct)
		if err != nil {
			t.Error(err.Error())
		}

		ct2, err := ssJWK2.SymmetricEncrypt(data)
		if err != nil {
			t.Error(err.Error())
		}

		pt2, err := ssJWK1.SymmetricDecrypt(ct2)
		if err != nil {
			t.Error(err.Error())
		}

		if string(pt) != string(data) ||
			string(pt2) != string(data) ||
			string(pt) != string(pt2) {
			t.Error("Symmetric encryption | decryption is failing.")
		}
	}
}

// t.Log () equivalent to fmt.Print line but concurrently safe
// t.Fail() will show that a test case failed
// t. FailNow() safely exit (test?) without continuing
// t.Error() = t.Log() + t.Fail()
// t.Fatal() = t.Log() + t.FailNow()
