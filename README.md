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

Using two terminals to run `make runs` and `make runc`.

```shell=
# Success test: client verify server certificate success
make runs
make runc

# Fail test: client verify server certificate fail
make runs
make runc-err-cert

# Fail test: client verify server certificate fail, with mismatched hostname
make runs
make runc-err-cn

# Fail test: server verify client certificate fail
make runs-verify
make runc

# Success test: Mutual verifications success
make runs-verify
make runc-with-cert
```

# TODO

The directory `TODO/` is a standalone program, integrate it with this project.
