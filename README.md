HTTPS server with sel-signed certificate.

## Re-generate certificates.

```
cd certificate/
./main
```

### Run HTTPS server

```
cd https/
./main
```

### Run client

```
cd httpc/
./main

```

### Run client with `curl`

```
cd httpc/
curl --cacert ./cert.pem https://localhost:4001/hello
```
