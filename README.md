
# Password Hash API

## Purpose

This is a go server that will hash a password using SHA 512.

## Startup

To run the server on ubuntu: `sudo go run main.go password.go statistics.go`

You must start it as root because it binds to port 80, a privileged port. 

## Client actions

To hash a password: `curl --data "password=angryMonkey1" http://localhost/hash`

This will pass back an id (which the client can save) and, eventually, the hashed value of the password.

To retrieve a cached value of a password: `curl http://localhost/hash/1` where `1` is the previously retrieved id. If the id was not previously returned, a 404 error will be returned.

To shut the server down (finishes existing requests, but new requests will be sent a 503 error): `curl http://localhost/shutdown`

## Development notes

run `go fmt <filename>` on any files.

If you are creating any http tests, make sure you run `resetVariablesToStartingValues` before any tests to reset the global variables.

You can run tests by running `go test`
