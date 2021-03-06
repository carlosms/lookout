# lookout [![Build Status](https://travis-ci.org/src-d/lookout.svg)](https://travis-ci.org/src-d/lookout) [![GoDoc](https://godoc.org/gopkg.in/src-d/lookout?status.svg)](https://godoc.org/github.com/src-d/lookout)

A service for assisted code review, that allows running custom code Analyzers on pull requests.

# Installation

`go get github.com/src-d/lookout`

## Dependencies

The included [`./docker-compose.yml`](./docker-compose.yml) allows to start all dependencies using [Docker Compose](https://docs.docker.com/compose/) 

* [bblfshd](https://github.com/bblfsh/bblfshd), on `localhost:9432`
* [PostgreSQL](https://www.postgresql.org/), on `localhost:5432` password `example`

Clone the repository, or download [`./docker-compose.yml`](./docker-compose.yml), and run:

```bash
docker-compose up
```


# Example

## SDK

If you are developing an Analyzer, please check [SDK documentation](./sdk/README.md).

It includes a curl-style binary that allows to trigger Analyzers directly, without launching a full lookout server.

## Server

To trigger the analysis on an actual pull request of a GitHub repository do:

1. Start an analyzer
Any of the analyzers or a default dummy one, included in this repository
    ```
    go build -o analyzer ./cmd/dummy
    ./analyzer serve
    ```
1. Start a lookout server
    1. With posting analysis results on GitHub
        - Obtain [GitHub access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)
        - Run `lookout serve --github-token <token> --github-user <user> <repository>`
    1. Without posting analysis results (only printing)
        - `lookout serve --dry-run <repository>`
1. Create a new pull requires in the repository


# Contribute

[Contributions](https://github.com/src-d/lookout/issues) are more than welcome, if you are interested please take a look to
our [Contributing Guidelines](CONTRIBUTING.md).

# Code of Conduct

All activities under source{d} projects are governed by the [source{d} code of conduct](https://github.com/src-d/guide/blob/master/.github/CODE_OF_CONDUCT.md).

# License
Affero GPL v3.0, see [LICENSE](LICENSE).

SDK package in `./sdk` is released under the terms of the [Apache License v2.0](./sdk/LICENSE)
