
# Hash api

To run the server on ubuntu: `sudo go run main.go`

You must start it as root because it binds to port 80, a privileged port. 

# Client actions

To hash a password: `curl  --data "password=angryMonkey1" http://localhost/hash`

To shut the server down (finishes existing requests, but new requests will be sent a 503 error): `curl http://localhost/shutdown`

# Development notes

run `go fmt <filename>` on any files.

If you are creating any http tests, make sure you run `resetVariablesToStartingValues` before any tests to reset the global variables.

You can run tests by running `go test`
