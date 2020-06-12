package main

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

func main() {

	fmt.Println(qrcode.Encode("otpauth://totp/kisexu@gmail.com?secret=DpI45HCEBCJK6HG7", qrcode.Medium, 256))
	// err := qrcode.WriteFile("otpauth://totp/kisexu@gmail.com?secret=DpI45HCEBCJK6HG7", qrcode.Medium, 256, "qr.png")
	// if err != nil {
	// 	fmt.Println("write error")
	// }
}
