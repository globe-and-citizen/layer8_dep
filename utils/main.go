package main

//https://chromium.googlesource.com/external/github.com/square/go-jose/+/5f0574b31ad6b505801faec07b6bcf079cbda9e4/jwk.go

import (
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"test/utils"

	"github.com/lestrrat-go/jwx/jwk"
)

func main() {
	// TEST ECDSA
	privJWK_ecdsa, pubJWK_ecdsa, err := utils.GenerateKeyPair(utils.ECDSA)
	priv_bs, _ := json.MarshalIndent(privJWK_ecdsa, "", "  ")
	pub_bs, _ := json.MarshalIndent(pubJWK_ecdsa, "", "  ")

	fmt.Println(string(priv_bs))
	fmt.Println(string(pub_bs))

	testPriv, _ := privJWK_ecdsa.ExportKeyAsGoType()
	if _, ok := testPriv.(*ecdsa.PrivateKey); ok {
		fmt.Println("Private ECDSA seems to export.")
	}

	testPub, _ := pubJWK_ecdsa.ExportKeyAsGoType()
	if _, ok := testPub.(*ecdsa.PublicKey); ok {
		fmt.Println("Public ECDSA seems to export")
	}

	// TEST ECDH
	// Generate First Key pair
	privJWK, pubJWK, err := utils.GenerateKeyPair(utils.ECDH)
	if err != nil {
		fmt.Println(err.Error())
	}

	privKey_unCasted, err := privJWK.ExportKeyAsGoType()
	if err != nil {
		fmt.Println(err.Error())
	}

	var privKey *ecdh.PrivateKey
	if pk, ok := privKey_unCasted.(*ecdh.PrivateKey); ok {
		privKey = pk
	}

	pubKey_unCasted, err := pubJWK.ExportKeyAsGoType()
	if err != nil {
		fmt.Println(err.Error())
	}

	var pubKey *ecdh.PublicKey
	if pk, ok := pubKey_unCasted.(*ecdh.PublicKey); ok {
		pubKey = pk
	}

	// Generate Second Key pair
	privJWK2, pubJWK2, err := utils.GenerateKeyPair(utils.ECDH)
	if err != nil {
		fmt.Println(err.Error())
	}
	privKey_unCasted2, err := privJWK2.ExportKeyAsGoType()
	if err != nil {
		fmt.Println(err.Error())
	}
	var privKey2 *ecdh.PrivateKey
	if pk, ok := privKey_unCasted2.(*ecdh.PrivateKey); ok {
		privKey2 = pk
	}
	pubKey_unCasted2, err := pubJWK2.ExportKeyAsGoType()
	if err != nil {
		fmt.Println(err.Error())
	}
	var pubKey2 *ecdh.PublicKey
	if pk, ok := pubKey_unCasted2.(*ecdh.PublicKey); ok {
		pubKey2 = pk
	}

	// Derive the two shared secrets
	ss, err := privKey.ECDH(pubKey2)
	if err != nil {
		fmt.Println(err.Error())
	}
	ss2, err := privKey2.ECDH(pubKey)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(base64.URLEncoding.EncodeToString(ss))
	fmt.Println(base64.URLEncoding.EncodeToString(ss2))
}

type PublicKey struct { // This is a ecdsa.PublicKey
	elliptic.Curve
	X, Y *big.Int
}

type PubJWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	Kid string `json:"kid"`
}

// {
// 	"kty":"EC",
// 	"crv":"P-256",
// 	"x":"f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
// 	"y":"x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0",
// 	"kid":"Public key used in JWS spec Appendix A.3 example"
// }

func test6() {
	// Create 2 ecdh keys
	privKey1, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//privKey2, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// privKey1, _ := ecdh.P256().GenerateKey(rand.Reader)
	// privKey2, _ := ecdh.P256().GenerateKey(rand.Reader)

	//fmt.Println(privKey1.X)
	//fmt.Println(privKey1.Y)
	base64.URLEncoding.EncodeToString(privKey1.X.Bytes())

	// create JWK
	JWKoriginal := &PubJWK{
		Kty: "EC",
		Crv: privKey1.Params().Name,
		X:   base64.URLEncoding.EncodeToString(privKey1.X.Bytes()),
		Y:   base64.URLEncoding.EncodeToString(privKey1.Y.Bytes()),
		Kid: "Key One",
	}
	fmt.Println("original: ", JWKoriginal)
	JWK_BS, err := json.Marshal(JWKoriginal)
	if err != nil {
		fmt.Println(err.Error())
	}

	// send over wire
	var JWKtransmitted PubJWK
	err = json.Unmarshal(JWK_BS, &JWKtransmitted)
	if err != nil {
		fmt.Println(err.Error())
	}

	// set equal to ecdsa.PublicKey
	pk := new(ecdsa.PublicKey)
	pk.Curve = elliptic.P256()
	xBytes, _ := base64.URLEncoding.DecodeString(JWKtransmitted.X)
	pk.X = new(big.Int).SetBytes(xBytes)
	yBytes, _ := base64.URLEncoding.DecodeString(JWKtransmitted.Y)
	pk.Y = new(big.Int).SetBytes(yBytes)
	finally_pubECDHKey, _ := pk.ECDH()

	// Convert privKey1 to *ecdh.Privatekey
	converted_PrivateKey1, _ := privKey1.ECDH()

	// Get second keypair
	privKey2, _ := ecdh.P256().GenerateKey(rand.Reader)
	pubKey2 := privKey2.PublicKey()

	// call .ECDH() & derive my shared secrets
	ss1, _ := converted_PrivateKey1.ECDH(pubKey2)
	ss2, err := privKey2.ECDH(finally_pubECDHKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("ss1: ", base64.URLEncoding.EncodeToString(ss1))
	fmt.Println("ss2: ", base64.URLEncoding.EncodeToString(ss2))

	// derive a ss
}

func test5() {
	pk := new(ecdsa.PublicKey)
	pk.Curve = elliptic.P256()
	bigNumBS, _ := base64.URLEncoding.DecodeString("f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU")
	pk.X = new(big.Int).SetBytes(bigNumBS)
	bigNumBS2, _ := base64.URLEncoding.DecodeString("x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0")
	pk.Y = new(big.Int).SetBytes(bigNumBS2)
	fmt.Println("r u getting it? ", pk)
	pk.ECDH()

}

func test4() {
	//conversions
	privKey1, _ := ecdh.P256().GenerateKey(rand.Reader)
	pubKey1 := privKey1.PublicKey()
	pubKey1_marshalled, _ := x509.MarshalPKIXPublicKey(pubKey1)
	pubKey1_parsed, _ := x509.ParsePKIXPublicKey(pubKey1_marshalled)

	var forJWKProcessing *ecdsa.PublicKey
	if pk1, ok := pubKey1_parsed.(*ecdsa.PublicKey); ok {
		fmt.Printf("Conversion Success: pubKey1_parsed of type `%T`\n", pk1)
		forJWKProcessing = pk1
		fmt.Println(forJWKProcessing)
	} else {
		fmt.Printf("Conversion Failed: pubKey1_parsed of type `%T` is NOT of type %T\n", pubKey1_parsed, pk1)
	}
	jwkMaybe := make([]byte, 10)
	jwk.ParseRawKey(forJWKProcessing.X.Bytes(), jwkMaybe)

	fmt.Printf("%T, key: %v", jwkMaybe, jwkMaybe)

	// Get second keypair
	// privKey2, _ := ecdh.P256().GenerateKey(rand.Reader)
	// pubKey2 := privKey2.PublicKey()

	// // Derive my shared secrets
	// ss1, _ := privKey1.ECDH(pubKey2)
	// ss2, err := privKey2.ECDH(myNewPubKey)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// fmt.Println("ss1: ", base64.URLEncoding.EncodeToString(ss1))
	// fmt.Println("ss2: ", base64.URLEncoding.EncodeToString(ss2))

}

// If I can get a ecdsa, then I can get a ecdh
func test3() {
	//conversions
	privKey1, _ := ecdh.P256().GenerateKey(rand.Reader)
	pubKey1 := privKey1.PublicKey()
	pubKey1_marshalled, _ := x509.MarshalPKIXPublicKey(pubKey1)
	pubKey1_parsed, _ := x509.ParsePKIXPublicKey(pubKey1_marshalled)

	var myNewPubKey *ecdh.PublicKey
	if pk1, ok := pubKey1_parsed.(*ecdsa.PublicKey); ok {
		fmt.Printf("Conversion Success: pubKey1_parsed of type `%T`", pk1)
		myNewPubKey, _ = pk1.ECDH()
		fmt.Printf("one step closer %T\n", pk1)
		fmt.Printf("one step closer %T\n", myNewPubKey)
	} else {
		fmt.Printf("Conversion Failed: pubKey1_parsed of type `%T` is NOT of type %T\n", pubKey1_parsed, pk1)
	}

	// Get second keypair
	privKey2, _ := ecdh.P256().GenerateKey(rand.Reader)
	pubKey2 := privKey2.PublicKey()

	// Derive my shared secrets
	ss1, _ := privKey1.ECDH(pubKey2)
	ss2, err := privKey2.ECDH(myNewPubKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("ss1: ", base64.URLEncoding.EncodeToString(ss1))
	fmt.Println("ss2: ", base64.URLEncoding.EncodeToString(ss2))

}

func test2() {

	ecdsaKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ecdsaPubKey := ecdsaKey.Public()
	var nextPubKey *ecdsa.PublicKey
	if pk, ok := ecdsaPubKey.(*ecdsa.PublicKey); !ok {
		fmt.Println("not okay")
	} else {
		nextPubKey = pk
		fmt.Println("okay")
	}

	whatAmI, _ := nextPubKey.ECDH()

	fmt.Printf("r u getting it? %T\n", whatAmI)

	privKey1, _ := ecdh.X25519().GenerateKey(rand.Reader)
	pubKey1 := privKey1.PublicKey()

	privKey2, _ := ecdh.X25519().GenerateKey(rand.Reader)
	pubKey2 := privKey2.PublicKey()

	key_BS, err := x509.MarshalPKIXPublicKey(pubKey1) // convert my struct to ASN1 DER
	if err != nil {
		fmt.Println("Failed to MarshalOKIXPublicKey: ", err.Error())
	}
	fmt.Println("base64: ", base64.URLEncoding.EncodeToString(key_BS))

	b64pubKey := base64.URLEncoding.EncodeToString(key_BS)
	return_of_key_BS, _ := base64.URLEncoding.DecodeString(b64pubKey)

	pubKeyUMM, err := x509.ParsePKIXPublicKey(return_of_key_BS) // take my ASN1 DER => struct
	if err != nil {
		fmt.Println("Failed to ParsePKIXPubliKey: ", err.Error())
	}

	var pubKeyUMM_asserted *ecdh.PublicKey
	if pk, ok := pubKeyUMM.(*ecdh.PublicKey); !ok {
		fmt.Printf("expected *ecdsa.PublicKey, got %T\n", pubKeyUMM)
	} else {
		fmt.Printf("expected and got, %T\n", pubKeyUMM)
		pubKeyUMM_asserted = pk
	}

	fmt.Printf("%T, %T <= these shoudl equal", pubKey1, pubKeyUMM)

	//(func(k *PrivateKey) ECDH (remote*PublicKey)([]byte, error)
	ss1, _ := privKey1.ECDH(pubKey2)
	ss2, _ := privKey2.ECDH(pubKeyUMM_asserted)

	fmt.Println("len", len(ss1))
	fmt.Println("ss1: ", base64.URLEncoding.EncodeToString(ss1))
	fmt.Println("ss2: ", base64.URLEncoding.EncodeToString(ss2))
	// ss3, err := privKey2.ECDH(pubKeyUMM)
	// if err != nil {
	// 	fmt.Println("Failed to do ECDH with *ecds.PublicKey: ", err.Error())
	// }

}

func func1() {
	// _, pubKey1, err := GenerateKeyPair(ECDH_ALGO)
	_, pubKey1, err := GenerateKeyPair(ECDH_ALGO)
	if err != nil {
		fmt.Println("Failed to generate key pair: ", err.Error())
	}
	if _, ok := pubKey1.(*ecdh.PublicKey); !ok {
		fmt.Printf("expected crypto.PublicKey, got %T\n", pubKey1)
	} else {
		fmt.Printf("pubKey1: expected & got %T\n", pubKey1)
	}

	key_BS, err := x509.MarshalPKIXPublicKey(pubKey1)
	if err != nil {
		fmt.Println("Failed to MarshalOKIXPublicKey: ", err.Error())
	}
	// //base64.URLEncoding.EncodeToString
	// strB64 := base64.URLEncoding.EncodeToString(key_BS)
	// //fmt.Println(strB64)

	// key_BS, err = base64.URLEncoding.DecodeString(strB64)
	// if err != nil {
	// 	fmt.Println("Failed to DecodeString from b64: ", err.Error())
	// }

	pubKey2, err := x509.ParsePKIXPublicKey(key_BS)
	if err != nil {
		fmt.Println("Failed to ParsePKIXPubliKey: ", err.Error())
	}

	if _, ok := pubKey2.(*crypto.PublicKey); !ok {
		fmt.Printf("expected crypto.PublicKey, got %T\n", pubKey2)
	}

	fmt.Println("well done?")

	// _, pubKey1, err := GenerateKeyPair(ECDH_ALGO)
	// pubJWK, err := EncodePublicKey(pubKey1, "pubKey1")
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("Unable to encode pubKey: %s", err.Error()).Error())
	// }
	// jsonBS, err := pubJWK.ExportJSONBytes()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(string(jsonBS))

	// DecodePublicKey(jsonBS)

	// rawBytes, err := pubJWK.ExportRawBytes()
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("Unable to encode pubKey: %s", err.Error()).Error())
	// }
	// fmt.Println(rawBytes)
}

type Algorithm int

var (
	ECDSA_ALGO Algorithm = 1
	ECDH_ALGO  Algorithm = 2
)

type JWK interface {
	ExportJSONBytes() ([]byte, error)
	//ExportHexString()
}

func (p PubJWK) ExportJSONBytes() ([]byte, error) {
	jsonBS, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal Public JWK to json encoded byte slice.")
	}
	return jsonBS, nil
}

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
		pri, err := ecdh.P256().GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("could not generate key pair: %s", err)
		}
		pub := pri.PublicKey()
		return pri, pub, nil
	default:
		return nil, nil, fmt.Errorf("unknown algorithm")
	}
}

// EncodePublicKey encodes a public key to JWK:
// https://datatracker.ietf.org/doc/html/rfc7517
func EncodePublicKey(pub crypto.PublicKey, keyid string) (JWK, error) {
	switch k := pub.(type) {
	case *ecdsa.PublicKey:
		PubJWK := PubJWK{
			Kty: "EC",
			Crv: k.Params().Name,
			X:   base64.URLEncoding.EncodeToString(k.X.Bytes()),
			Y:   base64.URLEncoding.EncodeToString(k.Y.Bytes()),
			Kid: keyid,
		}
		return PubJWK, nil
	case *ecdh.PublicKey:
		PubJWK := PubJWK{
			Kty: "EC",
			Crv: "P-256",
			X:   base64.URLEncoding.EncodeToString(k.Bytes()[:16]),
			Y:   base64.URLEncoding.EncodeToString(k.Bytes()[16:]),
			Kid: keyid,
		}
		return PubJWK, nil
	default:
		return nil, fmt.Errorf("Unknown public key type: %s", reflect.TypeOf(pub))
	}
}

// DecodePublicKey needs to take the byte slice from pubJWK.ExportJSONBytes and turn it
// back into a go struct.
func DecodePublicKey(JWK_BS []byte) (crypto.PublicKey, error) {
	pubJWK := &PubJWK{}
	json.Unmarshal(JWK_BS, pubJWK)
	fmt.Println("pubJWK.Crv: ", pubJWK.Crv)
	fmt.Println("pubJWK.Kty: ", pubJWK.Kty)
	fmt.Println("pubJWK.X: ", pubJWK.X)
	fmt.Println("pubJWK.Y: ", pubJWK.Y)
	fmt.Println("pubJWK.Kid: ", pubJWK.Kid)
	xBS, err := base64.URLEncoding.DecodeString(pubJWK.X)
	if err != nil {
		fmt.Println("unable to decode the x coordinate")
	}
	yBS, err := base64.URLEncoding.DecodeString(pubJWK.Y)
	if err != nil {
		fmt.Println("unable to decode the y coordinate")
	}

	/*
			TOMORROWS LABOUR: func (k *PrivateKey) ECDH(remote *PublicKey) ([]byte, error)
			The above function is the accepted way to do ECDH.
			Therefore, I need to figure out how to have the function DecodePublicKey(...) *ecdh.PublicKey
			hmmm

			boring.New
			elliptic.Unmarshal()
			elliptic.P256().ScalarMult()
			func (*PrivateKey) ECDH ¶
		    func (k *PrivateKey) ECDH(remote *PublicKey) ([]byte, error)
			maybe x509.ParsePKIXPublicKey? it takes a []byte
	*/
	var pubKey []byte
	pubKey = append(xBS, yBS...)

	fmt.Println(pubKey)

	return nil, fmt.Errorf("temp error...")

	// var (
	// 	key []byte
	// 	err error
	// )
	// algo, err := strconv.Atoi(encoded[len(encoded)-1:])
	// if err != nil {
	// 	return nil, fmt.Errorf("could not decode public key: %s", err)
	// }
	// key, err = hex.DecodeString(encoded[:len(encoded)-5])
	// if err != nil {
	// 	return nil, fmt.Errorf("could not decode public key: %s", err)
	// }
	// switch Algorithm(algo) {
	// case ECDSA_ALGO:
	// 	x := new(big.Int).SetBytes(key[:len(key)/2])
	// 	y := new(big.Int).SetBytes(key[len(key)/2:])
	// 	return &ecdsa.PublicKey{
	// 		Curve: elliptic.P256(),
	// 		X:     x,
	// 		Y:     y,
	// 	}, nil
	// case ECDH_ALGO:
	// 	var pub [32]uint8
	// 	copy(pub[:], key)
	// 	return pub, nil
	// default:
	// 	return nil, fmt.Errorf("unknown algorithm")
	// }
}
