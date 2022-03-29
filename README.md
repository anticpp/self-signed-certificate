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
make runc-no-verify        # Succ
make runc-verify-no-ca     # Fail
make runc-verify-host-err  # Fail
make runc-verify-succ      # Succ
```

# Test mTLS

```shell=
make runs-verify
```

```shell=
make runc-verify-succ   # Fail
make runc-with-cert     $ Succ
```

# TODO

[ ] The directory `TODO/` is a standalone program, integrate it with this project.
[X] Add DNS(localhost) in SANs to certificates, to allow client visit `https://localhost:4433/`.
    Stuck in problem: When adds SANs to certificates, the ssl client won't use CN to match hostname??
    Solved: According to RFC6215, if SANs extention does exist, match SANs, else match subject CN.
