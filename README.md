# pihole-nodes-sync

Project used to sync pihole nodes from a central server.

<!--toc:start-->

- [pihole-nodes-sync](#pihole-nodes-sync)
  - [Functionalities](#functionalities)
  - [Getting Started](#getting-started)
  - [Deployment](#deployment)
    - [Docker compose](#docker-compose)
    - [Linux](#linux)
  - [Contribution](#contribution)
  - [MakeFile](#makefile)
  - [Improvements](#improvements)
  <!--toc:end-->

## Functionalities

- Make a high available DNS Server

## Getting Started

## Deployment

### Docker compose

```yaml
services:
  pihole_sync:
    image: ghcr.io/diegoramil/pihole-nodes-sync:latest
    container_name: pihole_sync
    restart: always
    environment:
      - TIMEOUT=<timeout_in_seconds>
      - PASSWORD=<your_parent_pwd>
      - BASE_URL=<your_parent_url>
      - CHILD_URLS=<your_node_url>
      - CHILD_PASSWORDS=<your_node_pwd>
      - SYNC_HOURS=<sync_hours_integer>
      - UPDATE_GRAVITY=<true|false>
      - PROFILE=<debug,normal>
```

### Linux

- clone the repository

```bash
git clone git@github.com:diegoramil/pihole-nodes-sync.git
cd pihole-nodes-sync
```

- change the env variables in the `.env` file

```bash
cp .env.example .env
```

- build and run the application

```bash
go mod tidy
make build
./main
```

## Contribution

- Fork the repository
- Create a new branch (`git checkout -b feature/your-feature`)
- Make your changes
- Commit your changes (`git commit -m 'Add some feature'`)
- Push to the branch (`git push origin feature/your-feature`)
- Open a pull request

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

- [x] enable the user to use multiple nodes instead of only one
- [x] change panics to errors
- [ ] add tests
- [ ] add documentation
- [x] add dev docs
- [x] upload docker image to registry
- [ ] reuse token and make a mechanism in order to reuse and update it
- [x] schedule the sync process
- [x] create a "debug" mode
