Self-signed certificate
===================================

# Test with openssl

```shell
make help
```

## Create certificates

```shell=
make ca
make certs
```

## Run test

Open terminal, run a server.

```shell=
make runs
```

Open the other terminal, run test clients. See the `Makefile` for testcases.

```shell=
make runc-no-verify        # Succ
make runc-verify-no-ca     # Fail
make runc-verify-host-err  # Fail
make runc-verify-succ      # Succ
```

## Test mTLS

```shell=
make runs-verify
```

```shell=
make runc-verify-succ   # Fail
make runc-with-cert     $ Succ
```

# httpsrv

See `cmd/httpsrv/`. Build with `make bin`.

## HTTP

```shell=
./bin/httpsrv server -insecure

./bin/httpsrv client -url http://localhost:4433/hello
```


## HTTPS

```shell=
./bin/httpsrv server -cert pki/certs/server/test-server.pem -key pki/certs/server/test-server-key.pem

./bin/httpsrv client -cacert pki/CA_test/cacert.pem 
```

## HTTPS skip verify

```shell=
./bin/httpsrv server -cert pki/certs/server/test-server.pem -key pki/certs/server/test-server-key.pem

./bin/httpsrv client -insecure
```

## mTLS

The mutual TLS will verify certificate on both side, thus, the client side must provide client certificate and key.

```shell=
./bin/httpsrv server -cert pki/certs/server/test-server.pem -key pki/certs/server/test-server-key.pem -verify -cacert pki/CA_test/cacert.pem

./bin/httpsrv client -cacert pki/CA_test/cacert.pem -cert pki/certs/client/test-client.pem -key pki/certs/client/test-client-key.pem
```

## AutoTLS

AutoTLS is a convenient way to use TLS as encrypted transport, without managing complicated CA/certificate jobs. This is done by creating temeperary certificate on server side, and the client side skip server certificate verification.

```shell=
./bin/httpsrv server -auto-tll

./bin/httpsrv client -insecure
```

# TODO

- [ ] In `httpsrv`, how to verify client certificate hostname.
- [ ] Build Go tool to generate ca, certificate for future reuse.
- [X] The directory `TODO/` is a standalone program, integrate it with this project.
    + Solved: Integerate into `httpsrv`.
- [X] Add DNS(localhost) in SANs to certificates, to allow client visit `https://localhost:4433/`.
    + Stuck in problem: When adds SANs to certificates, the ssl client won't use CN to match hostname??
    + Solved: According to RFC6215, if SANs extention does exist, match SANs, else match subject CN.


