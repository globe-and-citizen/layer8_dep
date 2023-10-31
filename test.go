package main

import (
	"errors"
	"fmt"
	"globe-and-citizen/layer8/utils"
)

func main() {
	test1()
	test2()
	err := errors.New("Two")
	errToPrint := fmt.Errorf("One: %w", err)
	fmt.Println(errToPrint)
}

// TEST ECDSA
func test2() {

	privJWK_ecdsa1, pubJWK_ecdsa1, _ := utils.GenerateKeyPair(utils.ECDSA)
	//privJWK_ecdsa2, pubJWK_ecdsa2, _ := utils.GenerateKeyPair(utils.ECDSA)

	data := []byte("GSIGN ME!")

	siganture, err := privJWK_ecdsa1.SignWithKey(data)
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

	verified, err := pubJWK_ecdsa1.CheckAgainstASN1Signature(siganture, data)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("verified? ", verified)

}

// TEST ECDH
func test1() {
	// Generate First Key pair
	privJWK, pubJWK, err := utils.GenerateKeyPair(utils.ECDH)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Generate Second Key pair
	privJWK2, pubJWK2, err := utils.GenerateKeyPair(utils.ECDH)
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
