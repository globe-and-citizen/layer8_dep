package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	
	"syscall/js"
	std_utils "globe-and-citizen/layer8/utils"
)

func main() {
	close := make(chan struct{}, 0)

	js.Global().Set("scram", map[string]interface{}{
		"first": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) != 4 {
				panic(errors.New("all arguments are required: (password, iterationCount, salt, combinedNonce)"))
			}

			passwordArg := args[0].String()
			iterationCountArg := args[1].Int()
			saltArg := args[2].String()
			combinedNonceArg := args[3].String()

			if passwordArg == "" || iterationCountArg == 0 || saltArg == "" || combinedNonceArg == "" {
				panic(errors.New("all arguments are required: (password, iterationCount, salt, combinedNonce)"))
			}

			salt, err := base64.StdEncoding.DecodeString(saltArg)
			if err != nil {
				panic(err)
			}

			saltedPassword := std_utils.HI(passwordArg, string(salt), iterationCountArg)
			clientKey := std_utils.HmacSha256(saltedPassword, []byte("Client Key"))
			storedKey := sha256.Sum256(clientKey)
			serverKey := std_utils.HmacSha256(saltedPassword, []byte("Server Key"))
			authMessage := "n=" + combinedNonceArg
			clientSignature := std_utils.HmacSha256(storedKey[:], []byte(authMessage))
			clientProof, err := std_utils.XOR(clientKey, clientSignature)
			if err != nil {
				panic(err)
			}
			serverSignature := std_utils.HmacSha256(serverKey, []byte(authMessage))
			return map[string]interface{}{
				"proof":     base64.StdEncoding.EncodeToString(clientProof),
				"signature": base64.StdEncoding.EncodeToString(serverSignature),
			}
		}),
		"final": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			signature1Arg := args[0].String()
			signature2Arg := args[1].String()

			if signature1Arg == "" || signature2Arg == "" {
				panic(errors.New("all arguments are required: (signature1, signature2)"))
			}

			signature1, err := base64.StdEncoding.DecodeString(signature1Arg)
			if err != nil {
				panic(err)
			}
			signature2, err := base64.StdEncoding.DecodeString(signature2Arg)
			if err != nil {
				panic(err)
			}

			return hmac.Equal(signature1, signature2)
		}),
	})

	<-close
}
