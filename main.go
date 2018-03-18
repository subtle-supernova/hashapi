package main

import "fmt"
import "crypto/sha512"
import "encoding/base64"
import "log"
import "net/http"
import "time"
import "os"
import "sync/atomic"

const TIME_TO_SLEEP = 5 //TODO change to 5
const PASSWORD_PARAM_NAME = "password"

const HASH_ENDPOINT_NAME = "/hash"
const SHUTDOWN_ENDPOINT_NAME = "/shutdown"

var inShutdownMode bool = false
var handlingAHashRequest bool = false
var hashId uint64 = 0

// TODO gofmt

func main() {

	http.HandleFunc(HASH_ENDPOINT_NAME, hash)
	http.HandleFunc(SHUTDOWN_ENDPOINT_NAME, registerShutdown)

	go checkForShutdownAndExit()
	log.Fatal(http.ListenAndServe(":80", nil))
}

func resetVariablesToStartingValues() { // TODO hack
	inShutdownMode = false
	handlingAHashRequest = false
	hashId = 0
}

func checkForShutdownAndExit() {
	for {
		time.Sleep(1 * time.Second)
		if inShutdownMode {
			if handlingAHashRequest {
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

	if inShutdownMode {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	atomic.AddUint64(&hashId, 1)
	fmt.Fprintln(w, hashId)
	w.(http.Flusher).Flush()

	r.ParseForm()
	password := r.FormValue(PASSWORD_PARAM_NAME)
	// TODO what if we don't see password?

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
