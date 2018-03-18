
# Hash api

To run the server: `sudo go run main.go`

You must start it as root because it binds to port 80, a privileged port.

# Client actions

To hash a password: `curl  --data "password=angryMonkey1" http://localhost/hash`

To shut the server down (finishes existing requests, but new requests will be sent a 503 error): `curl http://localhost/shutdown`

