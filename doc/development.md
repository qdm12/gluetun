# Development

## Setup

### Using VSCode and Docker

That should be easier and better than a local setup, although it might use more memory if you're not on Linux.

1. Install [Docker](https://docs.docker.com/install/)
    - On Windows, share a drive with Docker Desktop and have the project on that partition
    - On OSX, share your project directory with Docker Desktop
1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
1. Your dev environment is ready to go!... and it's running in a container :+1:

### Locally

Install [Go](https://golang.org/dl/), [Docker](https://www.docker.com/products/docker-desktop) and [Git](https://git-scm.com/downloads); then:

```sh
go mod download
```

And finally install [golangci-lint](https://github.com/golangci/golangci-lint#install)

## Commands available

```sh
# Build the entrypoint binary
go build cmd/main.go
# Test the entrypoint code
go test ./...
# Lint the code
golangci-lint run
# Build the Docker image
docker build -t qmcgaw/private-internet-access .
```

## Guidelines

The Go code is in the Go file [cmd/main.go](../cmd/main.go) and the [internal directory](../internal), you might want to start reading the main.go file.

See the [Contributing document](.github/CONTRIBUTING.md) for more information on how to contribute to this repository.
