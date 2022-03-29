

```
# Run server
./chttp server -cert ../server/test-server.pem -key ../server/test-server-key.pem

# Run client
./chttp client -cacert ../CA_test/cacert.pem -url https://localhost:4433/hello
```
