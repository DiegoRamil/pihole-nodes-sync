# pihole-nodes-sync

Project used to sync pihole nodes from a central server.

## Functionalities

- Make a high available DNS Server

## Getting Started

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```

## Improvements

- [ ] enable the user to use multiple nodes instead of only one
- [x] change panics to errors
- [ ] add tests
- [ ] add documentation
- [ ] add dev docs
- [x] upload docker image to registry
- [ ] reuse token and make a mechanism in order to reuse and update it
- [x] schedule the sync process
- [ ] create a "debug" mode
