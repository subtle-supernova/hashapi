package main

import "fmt"
import "crypto/sha512"
import "encoding/base64"
import "log"
import "net/http"
import "time"
import "os"

const TIME_TO_SLEEP = 5 //TODO change to 5
const PASSWORD_PARAM_NAME = "password"

const HASH_ENDPOINT_NAME = "/hash"
const SHUTDOWN_ENDPOINT_NAME = "/shutdown"

var inShutdownMode bool = false
var handlingAHashRequest bool = false
func main() {

	http.HandleFunc(HASH_ENDPOINT_NAME, hash)
	http.HandleFunc(SHUTDOWN_ENDPOINT_NAME, registerShutdown)

	go checkForShutdownAndExit()
	log.Fatal(http.ListenAndServe(":80", nil))
}

func checkForShutdownAndExit() { 
	for {
		time.Sleep(1 * time.Second)
		if(inShutdownMode) {
			if (handlingAHashRequest) {
				time.Sleep(sleepTimeSeconds() * time.Second)
			}
			os.Exit(0)
		}
	}
}


func registerShutdown(w http.ResponseWriter, r *http.Request) {
	inShutdownMode = newShutdownValue()
}

func hash(w http.ResponseWriter, r *http.Request) {
	handlingAHashRequest = true
	if (inShutdownMode) {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	r.ParseForm()
	password := r.FormValue(PASSWORD_PARAM_NAME)
	// fmt.Fprintln(w, password)

	time.Sleep(sleepTimeSeconds() * time.Second)
	fmt.Fprintf(w, hashPassword(password))
	handlingAHashRequest = false
}

func hashPassword(passwd string) string {
	hash := sha512.Sum512([]byte(passwd))

	return base64.StdEncoding.EncodeToString(hash[:])
}

func sleepTimeSeconds() time.Duration {
	return time.Duration(TIME_TO_SLEEP)
}

func newShutdownValue() bool {
	return true
}
