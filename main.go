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
import "encoding/json"
import "sync"


const TIME_TO_SLEEP = 5 //TODO change to 5
const PASSWORD_PARAM_NAME = "password"

const HASH_ENDPOINT_NAME = "/hash"
const HASH_WITH_SLASH_ENDPOINT_NAME = "/hash/"
const SHUTDOWN_ENDPOINT_NAME = "/shutdown"
const CLEAN_SHUTDOWN_CODE = 0
const SHUTDOWN_WAIT_CHECK = 1

type PasswordInfo struct {
	Id int
	PasswordHash string
}

type Statistics struct {
	Total int
	CumulativeTime int
	mux sync.Mutex
}

var inShutdownMode bool = false
var handlingAHashRequest bool = false
var hashId int32 = 0
var seenPasswords map[int]PasswordInfo
var stats Statistics

// TODO refactor, put some stuff into packages

func main() {

	seenPasswords = make(map[int]PasswordInfo)
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
	seenPasswords = make(map[int]PasswordInfo)
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

func (s *Statistics) statsOutput() string {
	total := s.Total
	averageTime := 0
	if (total != 0) {
		averageTime = s.CumulativeTime / s.Total
	}
	statsMap := map[string]int{"total": total, "average": averageTime}
	jsonMap, _ := json.Marshal(statsMap)
	return string(jsonMap)
}

func (s *Statistics) incrementTotal() {
	s.mux.Lock()
	s.Total += 1
	s.mux.Unlock()
}

func (s *Statistics) incrementCumulativeTime(timeInMilliseconds int) {
	s.mux.Lock()
	s.CumulativeTime += timeInMilliseconds
	s.mux.Unlock()
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
	hashedPassword := hashPassword(password)
	fmt.Fprintf(w, hashedPassword)
	int32HashId := int(hashId)
	seenPasswords[int32HashId] = PasswordInfo{
		int32HashId,
		hashedPassword,
	}
	endNanos := time.Now().UnixNano()
	stats.incrementTotal()
	stats.incrementCumulativeTime(int((endNanos - startNanos)/1000000))
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
