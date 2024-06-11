# Temporal.io Playground
This is a playground for Temporal.io. It is a simple project that demonstrates how to use Temporal.io to orchestrate a workflow.
It covers some basic and advanced features of Temporal.io.

## Prerequisites
- Go >1.22
- [Temporal.io Server](https://learn.temporal.io/getting_started/go/dev_environment/#set-up-a-local-temporal-service-for-development-with-temporal-cli)
*OR*
- Docker

## Getting Started
1. Clone the repository
2. Run the Temporal.io Server using the following command:
```shell
temporal server start-dev
```
3. Run the worker using the following command:
```shell
go run worker/main.go
```
4. Run the http gateway using the following command:
```shell
go run gateway/main.go
```

To view the temporal web UI and follow each step in the workflow execution, visit [http://localhost:8233](http://localhost:8233).
Visit [http://localhost:8091/start](http://localhost:8091/start) to start the workflow. By default the timeout is set to 10 seconds.
Go to [http://localhost:8091/signal](http://localhost:8091/signal) to signal the workflow to finish execution.
Alternatively, go to [http://localhost:8091/abort](http://localhost:8091/abort) to abort the workflow execution.
You can query the workflow status by visiting [http://localhost:8091/status](http://localhost:8091/status).
