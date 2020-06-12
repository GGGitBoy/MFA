package main

import (
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	err := qrcode.WriteFile("otpauth://totp/kisexu@gmail.com?secret=DpI45HCEBCJK6HG7", qrcode.Medium, 256, "qr.png")
	if err != nil {
		fmt.Println("write error")
	}
}
