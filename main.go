package main

import "fmt"
import "log"
import "net/http"
import "time"
import "os"
import "sync/atomic"
import "strconv"
import "strings"

const TIME_TO_SLEEP = 5
const PASSWORD_PARAM_NAME = "password"

const HASH_ENDPOINT_NAME = "/hash"
const HASH_WITH_SLASH_ENDPOINT_NAME = "/hash/"
const SHUTDOWN_ENDPOINT_NAME = "/shutdown"
const STATS_ENDPOINT_NAME = "/stats"
const CLEAN_SHUTDOWN_CODE = 0
const SHUTDOWN_WAIT_CHECK = 1

var inShutdownMode bool = false
var hashId int32 = 0
var seenPasswords map[int]Password
var stats Statistics

func main() {

	seenPasswords = make(map[int]Password)
	stats = *new(Statistics)

	http.HandleFunc(HASH_ENDPOINT_NAME, hash)
	http.HandleFunc(HASH_WITH_SLASH_ENDPOINT_NAME, hash)
	http.HandleFunc(SHUTDOWN_ENDPOINT_NAME, registerShutdown)
	http.HandleFunc(STATS_ENDPOINT_NAME, statisticsGet)

	go checkForShutdownAndExit()
	log.Fatal(http.ListenAndServe(":80", nil))
}

func checkForShutdownAndExit() {
	for {
		time.Sleep(SHUTDOWN_WAIT_CHECK * time.Second)
		if inShutdownMode {
			time.Sleep(sleepTimeSeconds() * time.Second)
			os.Exit(CLEAN_SHUTDOWN_CODE)
		}
	}
}

func statisticsGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, stats.statsOutput())
}

func registerShutdown(w http.ResponseWriter, r *http.Request) {
	inShutdownMode = true
}

func getIdPointerFromPath(path string) *int {
	pathWithoutHash := strings.Replace(path, HASH_WITH_SLASH_ENDPOINT_NAME, "", 1)
	i, err := strconv.Atoi(pathWithoutHash)
	if err == nil {
		return &i
	}

	return (*int)(nil)
}

func hashFromId(startNanos int64, w http.ResponseWriter, id int) {
	passwordStruct := seenPasswords[id]
	if passwordStruct.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {

		fmt.Fprintf(w, passwordStruct.PasswordHash)
		captureStatistics(startNanos)
		return
	}
}

func hash(w http.ResponseWriter, r *http.Request) {
	startNanos := time.Now().UnixNano()

	if inShutdownMode {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	id := 0

	var idPtr = getIdPointerFromPath(r.URL.Path)
	if idPtr != (*int)(nil) {
		id = *idPtr
	}

	if id != 0 {
		hashFromId(startNanos, w, id)
	} else {
		hashCreate(startNanos, w, r)
	}
}

func hashCreate(startNanos int64, w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	password := r.FormValue(PASSWORD_PARAM_NAME)
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	atomic.AddInt32(&hashId, 1)
	fmt.Fprintln(w, hashId)
	w.(http.Flusher).Flush()

	time.Sleep(sleepTimeSeconds() * time.Second)
	passwordObj := *new(Password)
	passwordObj.hashPassword(password)
	passwordObj.Id = int(hashId)
	seenPasswords[passwordObj.Id] = passwordObj

	// TODO should we capture both hash create requests and hash retrieval requests in the stats
	captureStatistics(startNanos)
}

func captureStatistics(startNanos int64) {
	endNanos := time.Now().UnixNano()
	stats.incrementTotal()
	stats.incrementCumulativeTime(int((endNanos - startNanos) / 1000000))
}

func sleepTimeSeconds() time.Duration {
	return time.Duration(TIME_TO_SLEEP)
}

func resetVariablesToStartingValues() { // TODO hack only used for testing
	inShutdownMode = false
	hashId = 0
	seenPasswords = make(map[int]Password)
	stats = *new(Statistics)
}
