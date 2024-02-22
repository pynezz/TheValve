# TheValve

- [TheValve](#thevalve)
  - [Security Features](#security-features)
  - [Installation](#installation)
  - [Dependencies](#dependencies)
  - [Technical Details](#technical-details)
    - [Program Flow](#program-flow)
      - [First time setup](#first-time-setup)
      - [Returning user](#returning-user)
    - [Storage Format](#storage-format)


TheValve is a tool written in Go, made to safely convserve sensitive pieces of data for the forseeable future while keeping it easy accessible.

## Security Features

When you open TheValve, you will be assured that your data have been kept in its pristine condition. It does this by leveraging AES-256-GCM encryption, which is a widely accepted standard for encrypting data. TheValve also uses memory resistant Argon2 key derivation to generate a key from your password, which makes it harder for attackers to brute force your password. TheValve also uses a random salt for each password, which makes it harder for attackers to use precomputed tables to crack your password.

## Installation

```bash
go get github.com/pynezz/TheValve
```

## Dependencies

Go version: 1.21.6

- Argon2
- AES-256-GCM
- Golang.org/x/crypto

## Technical Details

As expected, the password is not stored anywhere. By leveraging Argon2 key derivation, the hash of the password is not stored anywhere either. TheValve uses a random salt for each password, and computes the hash, which in turn is used to generate the key for the AES-256-GCM encryption. TheValve uses a random nonce for each encryption, which makes it harder for attackers to use precomputed tables to crack your password.

### Program Flow

#### First time setup

1. The user enters a username and password.
2. TheValve generates a random salt.
3. The hash of the password is computed using the salt and Argon2 parameters.
4. A random nonce is generated.
5. The user enters the data to be stored.
6. When the user is done, the data is encrypted with AES-256-GCM using the hash as the key and the nonce as the nonce.
7. TheValve stores the username, salt, nonce, and the encrypted data.
8. TheValve exits.

#### Returning user

1. The user enters the username and password.
2. TheValve looks up the salt and the nonce stored for the username.
3. The hash of the password is computed using the salt and Argon2 parameters.
4. TheValve validates that the generated hash matches the supplied plaintext password.
5. TheValve decrypts the data using the hash as the key and the nonce as the nonce.
6. TheValve displays the decrypted data.
7. The user can now choose to update the data, add a new entry, change the password, or exit.
8. TheValve stores the username, salt, nonce, and the encrypted data.
9. TheValve exits.

### Storage Format

TheValve stores the data in a file named `thevalve.json` in the current working directory. The file is a JSON file with the following format:

- section: owner
  - salt: random salt
  - nonce: random nonce
  - entries: array of entries
    - entryname: entry
    - data: encrypted data as base64 string

The section contains the owner and the entries
The owner authenticates the section
Example:

```json
"vault": {
    "schema": 0.1,
    "section": {
        "owner": "John Doe",// CRC32 hashed for quick lookup (ex: John Doe -> 0x6A3811C6)
        "salt": "1234567890",
        "entries": {
            "nonce": "1234567890",  // In the future, we might want to add a nonce for each entry
            "content": {
                "data": "0A3jdj1d2DKJvx+01diksdCDvKW2dk02dogdk2D(...)1sS4ijdf2CD+fO/Sfaks6532dVF="
            }
        }
    },
    "section": {
        "owner": "Jane Doe",
        "salt": "1234567890",
        "entries": {
            "nonce": "1234567890",
            "content": {
                "data": "23jdj1AKd2dFL1x/1diksd21DdkkLD2d0gwtg23(...)234ijdf2CD+fASfaks65321jcK="
            }
        }
    },
}
```

Decryted, the entries would look like this:

```json
...
"entries": {
    "nonce": "1234567890",
    "content": {
        "api_key": "1234567890",
        "other_secret": "abcdefghijklmnopqrstuvwxyz1234567890",
    }
}
...
```
