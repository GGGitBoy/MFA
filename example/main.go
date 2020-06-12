/*******************************************************
	File Name: main.go
	Author: An
	Mail:lijian@cmcm.com
	Created Time: 14/11/25 - 10:24:49
	Modify Time: 14/11/25 - 10:24:49
 *******************************************************/
package main

import (
	"fmt"

	"googleAuthenticator"
)

func createSecret(ga *googleAuthenticator.GAuth) string {
	secret, err := ga.CreateSecret(16)
	if err != nil {
		return ""
	}
	return secret
}

func getCode(ga *googleAuthenticator.GAuth, secret string) string {
	code, err := ga.GetCode(secret)
	if err != nil {
		return "*"
	}
	return code
}

func verifyCode(ga *googleAuthenticator.GAuth, secret, code string) bool {
	// 1:30sec
	ret, err := ga.VerifyCode(secret, code, 1)
	if err != nil {
		return false
	}
	return ret
}

func main() {
	// if len(os.Args) != 2 {
	// 	return
	// }

	// secret := os.Args[1]
	//secret := "IU7B5Q3VBL55Q645"

	ga := googleAuthenticator.NewGAuth()

	secret, _ := ga.CreateSecret()
	fmt.Println(secret)

	fmt.Println(getCode(ga, secret))

	// for i := 0; i < 10; i++ {
	// 	log.Println(createSecret(ga))
	// }
}
