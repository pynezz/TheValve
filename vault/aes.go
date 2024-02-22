package vault

type Aes struct {
	nonce     []byte
	secretKey []byte // Password
}

func Argon2ToAES(argon2EncodedHash string, secretKey string) {
	// Convert the argon2 hash to an AES key

}
