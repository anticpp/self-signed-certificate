Self-signed certificate
===================================

# Usage

```shell
make help
```

# Create certificates

Before testing, create necessary CA/server/client certificates.

```shell=
make ca
make certs
```

# Test

Open terminal, run a server.

```shell=
make runs
```

Open the other terminal, run test clients. See the `Makefile` for testcases.

```shell=
make runc-no-verify
make runc-verify-no-ca
make runc-verify-host-err
make runc-verify-succ
```

# Test mTLS

```shell=
make runs-verify
make runc-with-cert
```

# TODO

The directory `TODO/` is a standalone program, integrate it with this project.
