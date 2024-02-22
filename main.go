package main

import (
	"fmt"
	"os"

	"github.com/pynezz/thevalve/tui"
	vault "github.com/pynezz/thevalve/vault"
)

// This project is a vault securing API keys and other sensitive data

func main() {
	// fmt.Println("\033[0;33mHello, World!\033[0m")

	tui := tui.NewTUI()
	choice := tui.Greet()

	switch choice {
	case 1:
		fmt.Println("Authenticating...")
	case 2:
		fmt.Println("Listing all owners...")
	case 0:
		fmt.Println("Exiting...")
		os.Exit(0)
	}

	pass := "password"
	// salt := "saltsalt"

	p := vault.NewArgon2().InitArgon(pass)

	printableKey := p.GetPrintableKeyWithSalt(p.Salt)
	match, err := vault.ComparePasswordAndHash(pass, printableKey)

	fmt.Printf("Generated argon2 hash \033[0;35m%s\033[0m\n", printableKey)

	fmt.Printf("Generated new hash from the same password \033[1;36m%s\033[0m\n", pass)
	// newHash := generateFromPassword(pass, &p.params)

	// fmt.Println(generateFromPassword(pass, &p.params))

	fmt.Println("Comparing the password to the hash...")
	// match, err := comparePasswordAndHash(pass, printableKey)

	// $argon2i$v=19$m=65536,t=3,p=2$R2NObUdpc0kwcVZNMlNpcQ$Midu7c5xtCBH9k/6Ow7hvHvD6QyUsKGXzpD0mfGW21o
	// $argon2i$v=19$m=65536,t=3,p=2$skeMa4CdlIpFhOcnozC0UA$FsvlMdGeIK7HCSy4jfS+mvH5QwS6GeHzfKxHWHsZU5U
	if err != nil {
		fmt.Println(err)
	}

	if match {
		fmt.Println("\033[0;32mPasswords match\033[0m")
	} else {
		fmt.Println("\033[0;31mPasswords don't match\033[0m")
	}
}
