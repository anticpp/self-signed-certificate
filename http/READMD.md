
# HTTP

```shell=
./chttp server -insecure

./chttp client -url http://localhost:4433/hello
```


# HTTPS

```shell=
./chttp server -cert ../server/test-server.pem -key ../server/test-server-key.pem

./chttp client -cacert ../CA_test/cacert.pem 
```

# HTTPS skip verify

```shell=
./chttp server -cert ../server/test-server.pem -key ../server/test-server-key.pem

./chttp client -insecure
```

# mTLS

```shell=
./chttp server -cert ../server/test-server.pem -key ../server/test-server-key.pem -verify -cacert ../CA_test/cacert.pem
```
