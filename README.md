
# Password Hash API

## Purpose

This is a go API that will hash a password using SHA 512.

## Startup

To run the server on ubuntu: `go build && sudo ./hashapi`

You must start it as root because it binds to port 80, a privileged port. 

## Client actions

To hash a password: `curl --data "password=angryMonkey1" http://localhost/hash`

This will pass back an id (which the client can save).

To retrieve the hashed password: `curl http://localhost/hash/1` where `1` is the previously retrieved id. If the id was not previously returned, or if the password is not yet generated, a 404 error will be returned.

To shut the server down (finishes existing requests, but new requests will be sent a 503 error): `curl http://localhost/shutdown`

To access statistics about the number of passwords hashed and the average amount of time taken: `curl http://localhost/stats`

## Development notes

To run the server on ubuntu: `sudo go run main.go password.go statistics.go`

run `go fmt <filename>` on any changed files before checkin.

If you are creating any http tests, make sure you run `resetVariablesToStartingValues` before any tests to reset the global variables. You also need to add any other global variables to that function.

You can run all tests by running `go test` or a particular test by running `go test -run TestYYY` 

## TODO

Uses port 80 and requires sudo to run the server
No error message on invalid/nonexistent id (but returns appropriate error msg)
Uses the same handler for both hash and hash/id routes with some confusing/unnecessary logic based on the retrieved id for discerning between the GET and POST paths: basically attempts to retrieve an id on both routes before checking for the request type.
shutdown route does not use shutdown method on the server but instead sleeps in a loop and checks whether a shutdown has been received and then calls os.Exit!
Shutdown also does not wait for hashing go-routines to finish
getIdPointerFromPath function (used to parse the id from the GET request) returns a (*int)(nil) pointer, instead of just returning a tuple (int, error).
Stats reading not protected by read mutex
the password hash map is not mutex protected at all
Race condition on creating new hash IDs. uses atomic.AddInt32 to protect incrementing the hash value, but ignores the output of the function which is what should be used thereafter. So he's safely incrementing the value, but then using the pointer to the value afterwards - at the same time another request could come in and increment the value.
Chooses to protect each piece of the statistic with the same mutex, but in separate functions. This means at any time you have a race condition where the total count of requests doesn't match up with the total time.
Tons of test-specific race conditions due to resetVariablesToStartingValues method resetting global state.
Sprintf format is incorrent in several places
Uses port 80 locally
Ignores errors in places
