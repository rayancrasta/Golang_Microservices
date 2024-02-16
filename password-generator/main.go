package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Password in plain text
	plainTextPassword := "verysecret!"

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Plaintext Password: %s\n", plainTextPassword)
	fmt.Printf("Bcrypt Hash: %s\n", string(hashedPassword))
}
