// This is a password hash struct.
package main

import "crypto/sha512"
import "encoding/base64"

type Password struct {
	Id           int
	PasswordHash string
}

func (p *Password) hashPassword(passwd string) {
	hash := sha512.Sum512([]byte(passwd))

	p.PasswordHash = base64.StdEncoding.EncodeToString(hash[:])
}
