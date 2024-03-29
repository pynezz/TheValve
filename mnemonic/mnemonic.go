package mnemonic

func generateMnemonic() string {
	// Mnemonics are generated by hashing a random number with HMAC-SHA512
	// we could use the password as the key and the random number as the message
	// and then hash the result with SHA512

	// The result is then used to generate a mnemonic by splitting the hash into 12 words
	// and mapping the binary representation of the split hash values to a wordlist
	// Wordlists are available at several places, like the EFF wordlist (https://www.eff.org/dice)

	// Maybe a bit overkill for this project, but it's a fun nice-to-have feature
	return ""
}
