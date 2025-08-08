package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Generate hashes for test passwords
	passwords := []string{"admin123", "cashier123", "password123"}
	
	for _, password := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Error generating hash for %s: %v\n", password, err)
			continue
		}
		fmt.Printf("Password '%s': %s\n", password, string(hash))
	}
}
