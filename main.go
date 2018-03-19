package main

import "fmt"
import "crypto/sha512"
import "encoding/base64"
import "log"
import "net/http"
import "time"
import "os"
import "sync/atomic"
import "strconv"
import "strings"

const TIME_TO_SLEEP = 5 //TODO change to 5
const PASSWORD_PARAM_NAME = "password"

const HASH_ENDPOINT_NAME = "/hash"
const HASH_WITH_SLASH_ENDPOINT_NAME = "/hash/"
const SHUTDOWN_ENDPOINT_NAME = "/shutdown"

type PasswordInfo struct {
	Id int
	PasswordHash string
}

var inShutdownMode bool = false
var handlingAHashRequest bool = false
var hashId int32 = 0
var seenPasswords map[int]PasswordInfo


func main() {

	seenPasswords = make(map[int]PasswordInfo)

	http.HandleFunc(HASH_ENDPOINT_NAME, hash)
	http.HandleFunc(HASH_WITH_SLASH_ENDPOINT_NAME, hash)
	http.HandleFunc(SHUTDOWN_ENDPOINT_NAME, registerShutdown)

	go checkForShutdownAndExit()
	log.Fatal(http.ListenAndServe(":80", nil))
}

func resetVariablesToStartingValues() { // TODO hack
	inShutdownMode = false
	handlingAHashRequest = false
	hashId = 0
	seenPasswords = make(map[int]PasswordInfo)
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

func getIdPointerFromPath(path string) *int {
	pathWithoutHash := strings.Replace(path, HASH_ENDPOINT_NAME+"/","", 1)
	// TODO make this more robust nad handle garbage after the endpoint?
	i, err := strconv.Atoi(pathWithoutHash)
	if (err == nil) {
		return &i
	}

	return (*int)(nil)
}


func hash(w http.ResponseWriter, r *http.Request) {
	handlingAHashRequest = true

	if inShutdownMode {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	id := 0

	var idPtr = getIdPointerFromPath(r.URL.Path)
	if (idPtr != (*int)(nil)) {
		id = *idPtr
	}

	if (id != 0) {
		passwordStruct := seenPasswords[id]
		if (passwordStruct.Id == 0) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			fmt.Fprintf(w, passwordStruct.PasswordHash)
			return
		}
	}

	// TODO for the same password value we give different ids.

	atomic.AddInt32(&hashId, 1)
	fmt.Fprintln(w, hashId)
	w.(http.Flusher).Flush()

	r.ParseForm()
	password := r.FormValue(PASSWORD_PARAM_NAME)
	// TODO what if we don't see password?

	time.Sleep(sleepTimeSeconds() * time.Second)
	hashedPassword := hashPassword(password)
	fmt.Fprintf(w, hashedPassword)
	int32HashId := int(hashId)
	seenPasswords[int32HashId] = PasswordInfo{
		int32HashId,
		hashedPassword,
	}
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
