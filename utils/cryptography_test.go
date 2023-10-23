package utils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"math/rand"
	"slices"
	"testing"
	"time"
)

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

func Test_PublicKeyEquals() {

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

// func Test_SymmetricEncrypt(t *testing.T){}

// func

// func main() {
// 	test1()
// 	test2()
// 	err := errors.New("Two")
// 	errToPrint := fmt.Errorf("One: %w", err)
// 	fmt.Println(errToPrint)
// }

// TEST ECDSA
func test2() {

	privJWK_ecdsa1, pubJWK_ecdsa1, _ := GenerateKeyPair(ECDSA)
	//privJWK_ecdsa2, pubJWK_ecdsa2, _ := utils.GenerateKeyPair(utils.ECDSA)

	data := []byte("GSIGN ME!")

	signature, err := privJWK_ecdsa1.SignWithKey(data)
	if err != nil {
		panic(err.Error())
	}

	// _, errTest := pubJWK_ecdsa1.SignWithKey(data)
	// if errTest != nil {
	// 	fmt.Println(errTest.Error())
	// }

	// _, errTest2 := utils.SignData(pubJWK_ecdsa1, data)
	// if errTest2 != nil {
	// 	panic(errTest2)
	// }

	verified, err := pubJWK_ecdsa1.CheckAgainstASN1Signature(signature, data)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("verified? ", verified)

}

// TEST ECDH
func test1() {
	// Generate First Key pair
	privJWK, pubJWK, err := GenerateKeyPair(ECDH)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Generate Second Key pair
	privJWK2, pubJWK2, err := GenerateKeyPair(ECDH)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Derive the two shared secrets
	ssJWK1, err := privJWK.GetECDHSharedSecret(pubJWK2)
	if err != nil {
		fmt.Println(err.Error())
	}
	ssJWK2, err := privJWK2.GetECDHSharedSecret(pubJWK)
	if err != nil {
		fmt.Println(err.Error())
	}
	//bs, _ := json.MarshalIndent(sJWK1, "", "  ")
	//fmt.Println(string(bs))
	//fmt.Println(ssJWK1)
	//fmt.Println(ssJWK2)

	//Test 1
	data := []byte("HIT ME!")
	ct, err := ssJWK1.SymmetricEncrypt(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("ct: ", string(ct))
	pt, err := ssJWK2.SymmetricDecrypt(ct)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("pt:", string(pt))

	//Test 2
	data2 := []byte("BABY ONE MORE TIME!")
	ct2, err := ssJWK2.SymmetricEncrypt(data2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("ct2: ", string(ct2))
	pt2, err := ssJWK1.SymmetricDecrypt(ct2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("pt2:", string(pt2))

	if string(pt) == string(data) && string(pt2) == string(data2) {
		fmt.Println("r u getting it?")
	} else {
		fmt.Println("try harder...")
	}
}

// t.Log () equivalent to fmt.Print line but concurrently safe
// t.Fail() will show that a test case failed
// t. FailNow() safely exit (test?) without continuing
// t.Error() = t.Log() + t.Fail()
// t.Fatal() = t.Log() + t.FailNow()
