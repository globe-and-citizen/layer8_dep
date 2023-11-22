package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	std_utils "globe-and-citizen/layer8/utils"
	"syscall/js"
)

func main() {
	close := make(chan struct{}, 0)

	js.Global().Set("scram", map[string]interface{}{
		"first": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, resolve__reject []js.Value) interface{} {
				var (
					resolve = resolve__reject[0]
					reject  = resolve__reject[1]
				)

				if len(args) != 4 {
					reject.Invoke(js.Global().Get("Error").New("all arguments are required: (password, iterationCount, salt, combinedNonce)"))
					return nil
				}

				passwordArg := args[0].String()
				iterationCountArg := args[1].Int()
				saltArg := args[2].String()
				combinedNonceArg := args[3].String()

				if passwordArg == "" || iterationCountArg == 0 || saltArg == "" || combinedNonceArg == "" {
					reject.Invoke(js.Global().Get("Error").New("all arguments are required: (password, iterationCount, salt, combinedNonce)"))
					return nil
				}

				salt, err := base64.StdEncoding.DecodeString(saltArg)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return nil
				}

				saltedPassword := std_utils.HI(passwordArg, string(salt), iterationCountArg)
				clientKey := std_utils.HmacSha256(saltedPassword, []byte("Client Key"))
				storedKey := sha256.Sum256(clientKey)
				serverKey := std_utils.HmacSha256(saltedPassword, []byte("Server Key"))
				authMessage := "n=" + combinedNonceArg
				clientSignature := std_utils.HmacSha256(storedKey[:], []byte(authMessage))
				clientProof, err := std_utils.XOR(clientKey, clientSignature)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return nil
				}
				serverSignature := std_utils.HmacSha256(serverKey, []byte(authMessage))
				resolve.Invoke(map[string]interface{}{
					"proof":     base64.StdEncoding.EncodeToString(clientProof),
					"signature": base64.StdEncoding.EncodeToString(serverSignature),
				})
				return nil
			}))
		}),
		"final": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, resolve__reject []js.Value) interface{} {
				var (
					resolve = resolve__reject[0]
					reject  = resolve__reject[1]

					signature1Arg = args[0].String()
					signature2Arg = args[1].String()
				)

				if signature1Arg == "" || signature2Arg == "" {
					reject.Invoke(js.Global().Get("Error").New("all arguments are required: (signature1, signature2)"))
					return nil
				}

				signature1, err := base64.StdEncoding.DecodeString(signature1Arg)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return nil
				}
				signature2, err := base64.StdEncoding.DecodeString(signature2Arg)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return nil
				}

				resolve.Invoke(hmac.Equal(signature1, signature2))
				return nil
			}))
		}),
	})

	<-close
}
