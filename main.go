package main

import "fmt"
import "crypto/sha512"
import "encoding/base64"
import "log"
import "net/http"
import "time"

const TIME_TO_SLEEP = 5 //TODO change to 5
const PASSWORD_PARAM_NAME = "password"

func main() {

	http.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		password := r.FormValue(PASSWORD_PARAM_NAME)
		// fmt.Fprintln(w, password)

		time.Sleep(sleepTimeSeconds() * time.Second)
		fmt.Fprintf(w, hash(password))
	})

	log.Fatal(http.ListenAndServe(":80", nil))

}

func hash(passwd string) string {
	hash := sha512.Sum512([]byte(passwd))

	hash2 := hash[:]

	str := base64.StdEncoding.EncodeToString(hash2)

	return str
}

func sleepTimeSeconds() time.Duration {
	return time.Duration(TIME_TO_SLEEP)
}
