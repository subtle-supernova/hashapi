
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
