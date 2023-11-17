# temporal-workflow-scheduler

## Development

### Requirements

- Go 1.19+
- GNU Make
- git

### Local

Use https://github.com/temporalio/docker-compose and eventually the reverse proxy https://ngrok.com/ when not running on the same host. 
When not running local, create a `.env` file at the root and set the TEMPORAL_HOSTPORT environment variable with the hostport information from grok. For example syntax ```ngrok tcp 7233``` and look for the `Forwarding tcp connection`

To use the Temporal Cluster UI from a remote host, you can also use ```ngrok http 8080``` then access it with the provided https Forwarding URL.

### Build & Execute

When running on a macos/arm host:

```bash
make clean compile-only
./bin/darwin/worker
```

To compile, run tests and linter
```bash
make
```
