// This server hashes arbitrary passwords and provides statistics on how many password hash requests have been made. See https://github.com/mooreds/hashapi for more info
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
var createdPasswords map[int]Password
var stats Statistics

func main() {

	createdPasswords = make(map[int]Password)
	stats = *new(Statistics)

	http.HandleFunc(HASH_ENDPOINT_NAME, hash)
	http.HandleFunc(HASH_WITH_SLASH_ENDPOINT_NAME, hash)
	http.HandleFunc(SHUTDOWN_ENDPOINT_NAME, registerShutdown)
	http.HandleFunc(STATS_ENDPOINT_NAME, statisticsGet)

	go checkForShutdownAndExit()
	log.Fatal(http.ListenAndServe(":80", nil))
}

func statisticsGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, stats.statsOutput())
}

func registerShutdown(w http.ResponseWriter, r *http.Request) {
	inShutdownMode = true
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

	if r.Method == "GET" && id != 0 {
		hashFromId(w, id)
	} else if r.Method == "POST" && id == 0 {
		hashIdAndCreate(startNanos, w, r)
	} else {
		badRequestResponse(w)
		return
	}
}

func hashFromId(w http.ResponseWriter, id int) {
	passwordStruct := createdPasswords[id]
	if passwordStruct.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		fmt.Fprintf(w, passwordStruct.PasswordHash)
		return
	}
}

func hashIdAndCreate(startNanos int64, w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	password := r.FormValue(PASSWORD_PARAM_NAME)
	if password == "" {
		badRequestResponse(w)
		return
	}

	atomic.AddInt32(&hashId, 1)
	fmt.Fprintln(w, hashId)

	go hashCreate(startNanos, hashId, password)
}

func hashCreate(startNanos int64, id int32, password string) {
	time.Sleep(sleepTimeSeconds() * time.Second)

	passwordObj := *new(Password)
	passwordObj.hashPassword(password)
	passwordObj.Id = int(hashId)
	createdPasswords[passwordObj.Id] = passwordObj

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

func checkForShutdownAndExit() {
	for {
		time.Sleep(SHUTDOWN_WAIT_CHECK * time.Second)
		if inShutdownMode {
			time.Sleep(sleepTimeSeconds() * time.Second)
			os.Exit(CLEAN_SHUTDOWN_CODE)
		}
	}
}

func getIdPointerFromPath(path string) *int {
	pathWithoutHash := strings.Replace(path, HASH_WITH_SLASH_ENDPOINT_NAME, "", 1)
	i, err := strconv.Atoi(pathWithoutHash)
	if err == nil {
		return &i
	}

	return (*int)(nil)
}

func badRequestResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func resetVariablesToStartingValues() { // TODO hack only used for testing
	inShutdownMode = false
	hashId = 0
	createdPasswords = make(map[int]Password)
	stats.Total = 0
	stats.CumulativeTime = 0
}
