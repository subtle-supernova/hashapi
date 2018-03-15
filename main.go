package main

import "fmt"
import "crypto/sha512"
import "encoding/base64"
import "os"

func main() {
	argsWithoutProg := os.Args[1:]
	// TODO check length of args
	//if argsWithoutProg.size != 1
	password := argsWithoutProg[0]
	fmt.Println(hash(password))
}

func hash(passwd string) string {
	hash := sha512.Sum512([]byte(passwd))
	hash2 := hash[:]

	str := base64.StdEncoding.EncodeToString(hash2)

	return str
}


