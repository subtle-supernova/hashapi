package main

import "fmt"
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
const CLEAN_SHUTDOWN_CODE = 0
const SHUTDOWN_WAIT_CHECK = 1


var inShutdownMode bool = false
var handlingAHashRequest bool = false
var hashId int32 = 0
var seenPasswords map[int]Password
var stats Statistics

// TODO refactor, put some stuff into packages

func main() {

	seenPasswords = make(map[int]Password)
	stats = *new(Statistics)

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
	seenPasswords = make(map[int]Password)
	stats = *new(Statistics)
}

func checkForShutdownAndExit() {
	for {
		time.Sleep(SHUTDOWN_WAIT_CHECK * time.Second)
		if inShutdownMode {
			if handlingAHashRequest {
				time.Sleep(sleepTimeSeconds() * time.Second)
			}
			os.Exit(CLEAN_SHUTDOWN_CODE)
		}
	}
}

func statisticsGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, stats.statsOutput())
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
	handlingAHashRequest = true // TODO this is wrong, need to update to lock it down better
	startNanos := time.Now().UnixNano()

	if inShutdownMode {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	id := 0

// TODO stats should include any hashes returnbed, but not 404 errors
// Use this: https://tour.golang.org/concurrency/9

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


	r.ParseForm()
	password := r.FormValue(PASSWORD_PARAM_NAME)
	if (password == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	atomic.AddInt32(&hashId, 1)
	fmt.Fprintln(w, hashId)
	w.(http.Flusher).Flush()

	time.Sleep(sleepTimeSeconds() * time.Second)
	passwordObj := *new(Password)
	passwordObj.hashPassword(password)
	fmt.Fprintf(w, passwordObj.PasswordHash)
	passwordObj.Id = int(hashId)
	seenPasswords[passwordObj.Id] = passwordObj

	endNanos := time.Now().UnixNano()
	stats.incrementTotal()
	stats.incrementCumulativeTime(int((endNanos - startNanos)/1000000))
	handlingAHashRequest = false
}


func sleepTimeSeconds() time.Duration {
	return time.Duration(TIME_TO_SLEEP)
}

func newShutdownValue() bool {
	return true
}
