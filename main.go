// This server hashes arbitrary passwords and provides statistics on how many password hash requests have been made.
// See https://github.com/mooreds/hashapi for more info
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

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
var stats Statistics

var concurrentPaswordHash sync.Map

func findPort() string {
	port, exists := os.LookupEnv("HASHAPI_PORT")
	if exists {
		return ":" + port
	}
	return ":8888"
}

func main() {

	stats = *new(Statistics)
	concurrentPaswordHash = *new(sync.Map)

	http.HandleFunc(HASH_ENDPOINT_NAME, hash)
	http.HandleFunc(HASH_WITH_SLASH_ENDPOINT_NAME, hash)
	http.HandleFunc(SHUTDOWN_ENDPOINT_NAME, registerShutdown)
	http.HandleFunc(STATS_ENDPOINT_NAME, statisticsGet)

	go checkForShutdownAndExit()
	port := findPort()
	log.Println("Starting server on " + port + "...")
	log.Fatal(http.ListenAndServe(findPort(), nil))
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

	idPtr, err := getIdPointerFromPath(r.URL.Path)
	if err == nil {
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
	concurrentPasswordStruct, ok := concurrentPaswordHash.Load(id)
	if !ok {
		log.Println(fmt.Sprintf("Nothing found for id %v", id))
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		fmt.Fprintf(w, concurrentPasswordStruct.(Password).PasswordHash)
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

	newId := atomic.AddInt32(&hashId, 1)
	fmt.Fprintln(w, newId)

	go hashCreate(startNanos, newId, password)
}

func hashCreate(startNanos int64, id int32, password string) {
	time.Sleep(sleepTimeSeconds() * time.Second)

	passwordObj := *new(Password)
	passwordObj.hashPassword(password)
	passwordObj.Id = int(id)
	concurrentPaswordHash.Store(passwordObj.Id, passwordObj)

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
			log.Println("Shutting down the server")
			os.Exit(CLEAN_SHUTDOWN_CODE)
		}
	}
}

func getIdPointerFromPath(path string) (*int, error) {
	pathWithoutHash := strings.Replace(path, HASH_WITH_SLASH_ENDPOINT_NAME, "", 1)
	i, err := strconv.Atoi(pathWithoutHash)
	if err == nil {
		return &i, nil
	}

	return nil, err
}

func badRequestResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func resetVariablesToStartingValues() { // TODO hack only used for testing
	inShutdownMode = false
	hashId = 0
	stats.Total = 0
	stats.CumulativeTime = 0
}
