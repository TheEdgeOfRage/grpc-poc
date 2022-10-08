# gRPC streaming PoC

Setup:

```
make generate-proto
make build
```

Run:

```
./grpcserver dev-dune-qes-results results.json
./httpserver
```

The parameters for the grpc server are the S3 bucket name and key.
The above one is present on the Dune AWS account.

You can also prepend the above commands with `/usr/bin/time -f '%M' -q`,
if you have time installed (not the built-in shell command)
