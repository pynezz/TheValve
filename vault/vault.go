package vault

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

// The vault contains several sections which are owned by different people
// Each section contains a secret key and the owner of the secret key + the entries of that person

// The entries are the actual secrets, which are stored in the vault as key:value pairs

// For each section, the owner is authenticated with a secret key, which is stored in the vault as a argon2 hash
// Each section is encrypted with a different secret key, which is stored in the vault as a argon2 hash

// Custom types for easy conversion to base64
type (
	byteSlice []byte
	str       string
)

var (
	// Define the schema version, in case we change this in the future
	schema = 0.1

	// Not in use
	// Define magicNumber and separator as byte slices
	magicNumber = []byte{0x24, 0xbe, 0xef}             // Equivalent to $beef
	separator   = []byte{0x24, 0x2d, 0x2d, 0x2d, 0x24} // Equivalent to $---$
)

// The Vault is the main struct, which contains all the sections
// The Vault contains all the sections (the file)
type Vault struct {
	// The section is the value, the owner is the key. The owner is hashed with CRC32 for quick lookup
	hashmap map[string]section
}

// The section contains the owner and the entries
// The owner authenticates the section
// Example:
// vault: {
// 		"schema": 0.1,
//		section: {
//			"owner": "John Doe",		// CRC32 hashed for quick lookup (ex: John Doe -> 0x6A3811C6)
//			"salt": "1234567890",
//			"entries": {
// 				"nonce": "1234567890",	// In the future, we might want to add a nonce for each entry
// 				"content": {
//					"data": "23jdj1AKd2dFL1x/1diksd21DdkkLD2dogwtg23(...)234ijdf2CD+fASfaks65321jcK="
//				}
//			}
//		},
//		section: {
//			"owner": "Jane Doe",
//			"salt": "1234567890",
//			"entries": {
// 				"nonce": "1234567890",	// In the future, we might want to add a nonce for each entry
// 				"content": {
//					"data": "23jdj1AKd2dFL1x/1diksd21DdkkLD2dogwtg23(...)234ijdf2CD+fASfaks65321jcK="
//				}
//			}
//		},
// }

// Decryted, the entries would look like this:
// ...
//
//	"entries": {
//			"nonce": "1234567890",
//			"content": {
//				"api_key": "1234567890",
//				"other_secret": "abcdefghijklmnopqrstuvwxyz1234567890",
//			}
//		}
//
// ...
type section struct {
	owner  owner
	salt   byteSlice // The salt is used to regenerate the argon2 hash
	nonce  byteSlice // The nonce (16 bytes) is used to decrypt the section and is stored as a base64 encoded string
	argon2 params    // The argon2 parameters
}

// The owner represents a person who owns a section
type owner struct {
	name    string           // The name of the owner
	entries map[string]entry // The entries are stored as entryname:entry in encrypted form
}

// The struct containing the entry
// Example:
//
//	"entry": {
//		"owner": "John Doe",
//		"content": {
//			"api_key: "1234567890",
//		}
//	}
type entry struct {
	owner string
	// In the content map, the key is the name of the entry, the value is the stored data in the entry,
	content map[string]string // both are encrypted and represented as a base64 encoded string
}

var errMsg = map[string]string{
	"1": "No path provided - You need to provide a path to the file containing the vault",
	"2": "Error reading the old secret key",
	"3": "Error comparing the old secret key",
	"4": "The old secret key does not match the stored secret key",
	"5": "Error reading the new secret key",
	"6": "Error when creating base64 encoded hash",
	"7": "Passwords don't match",
}

// TODO: This doesn't seem right
func (s *section) newEntry(owner string, content map[string]string) {
	v := entry{owner: owner, content: content}

	v.content = content
	v.owner = owner

	// Add the entry to the section
	s.owner.entries[owner] = v

	// Write the entry to the vault file
	WriteEntryToFile(&v)
}

func (v *Vault) GetVaultCount() int {
	// Get the count of the vault
	return len(v.hashmap)
}

func (s str) StrToBase64() string {
	// Convert a string to a base64 encoded string
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func (b byteSlice) BytesToBase64() string {
	// Convert a byte slice to a base64 encoded string
	return base64.StdEncoding.EncodeToString(b)
}

// Create a new section
func (v *Vault) NewSection(sectionName string, ownerName string, password string, salt byteSlice, nonce byteSlice) string {
	// Create a new section
	s := section{
		salt:  byteSlice(salt.BytesToBase64()),
		nonce: byteSlice(nonce.BytesToBase64()),
		owner: owner{
			name:    ownerName,
			entries: make(map[string]entry),
		},
	}

	// Add the section to the vault
	v.hashmap[sectionName] = s

	return "success"
}

func (o *owner) ChangeSecretKey(newSecretKey string) error {
	// Change the secret key of the owner
	reader := bufio.Reader{}

	// First, the owner will have to enter the old password
	fmt.Println("Enter the old secret key: ")
	oldSecretKey, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading the old secret key")
	}

	// We check if it matches the stored secret key
	match, err := generateFromPassword(oldSecretKey, o)
	if err != nil {
		errors.New(errMsg["3"])
	}

	// If not, we return
	if !match {
		fmt.Println("The old secret key does not match the stored secret key")
		return errors.New(errMsg["4"])
	}

	// If it matches, we ask for the new secret key
	fmt.Println("Enter the new secret key: ")
	newSecretKey, err = reader.ReadString('\n')
	if err != nil {
		return errors.New(errMsg["5"])
	}

	fmt.Println("Enter the new secret key again: ")
	verifyNewSecretKey, err := reader.ReadString('\n')

	if newSecretKey != verifyNewSecretKey {
		fmt.Println("That didn't work. Passwords don't match. Try again.")
		return errors.New(errMsg["7"])
	}

	o.secretKey = NewArgon2().InitArgon(newSecretKey).EncodedHash

	return nil
}

func (s *section) decrypt() {
	// Decrypt the section
}

func (s *section) encrypt() {
	// Encrypt the section
}

func (v *entry) deleteEntry(s *section, entryName string) {
	// Delete an entry from the section
	for i, e := range s.entries {
		if e.owner == v.owner {
			delete(s.entries, i)
		}
	}
}

func verifyAuthSession(s *section, owner *owner) bool {
	// Verify the authentication session
	return true
}

// Create a function to create a new vault
func (s *section) createEntry() *entry {
	if valid := verifyAuthSession(s, &s.owner); !valid {
		fmt.Println("[!] Not authenticated")
		return nil
	}

	// Create a new entry
	v := entry{owner: s.owner.name, content: make(map[string]string)}
	v.content = make(map[string]string)

	return &v
}

func WriteEntryToFile(v *entry) {
	// Write the entry to the vault file
}

func (v *Vault) AuthUser(password string) section {
	// Authenticate the user
	return section{}
}

func ReadVaultFile(path string) *Vault {
	// Read the vault from a file

	if path == "" {
		fmt.Println("[!] No path provided - You need to provide a path to the file containing the vault")
	}
	return nil
}

// Write the vault to a file
func (v *Vault) WriteVaultToFile() {
	// Permissions for the file
	permissions := 0644

	// Start with the magicNumber
	byteArray := make([]byte, 0)
	byteArray = append(byteArray, magicNumber...)

	// Append the separator
	byteArray = append(byteArray, separator...)

	// TODO: Convert your Vault data into a byte slice and append it to byteArray
	// For example, if your Vault has a Name field, you could append it like this:
	// byteArray = append(byteArray, []byte(v.Name)...)
	// Don't forget to add separators as needed between different parts of your Vault data

	// Append another separator at the end (optional, depending on your format needs)
	byteArray = append(byteArray, separator...)

	// Write the byteArray to a file
	err := os.WriteFile("vault.steel", byteArray, os.FileMode(permissions))
	if err != nil {
		// Handle error
		panic(err)
	}
}
