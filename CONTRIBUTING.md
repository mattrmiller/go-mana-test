# Contributing

## Setup
1. Clone this project with the Go tooling
`go get -u github.com/mattrmiller/go-mana-test`

## Adding New Go Packages
We use [dep](https://github.com/golang/dep/) as our dependency manager.

This is a service, so in order to produce reproducible builds we check in our `vendor` directory.

Adding a new package:
```bash
dep ensure -add github.com/foo/bar
```

From there, when you need to update to a new version. Just update to [Gopkg.toml](./Gopkg.toml) file following the [version schema](https://github.com/golang/dep/blob/master/docs/Gopkg.toml.md).
```bash
dep ensure -update
```

## Tests
Please provide tests for any core logic changes. To run tests:

```bash
make test
```
* To run only tests from a package: `go test <path>`
* To run only a single test: `go test -run <test-name> <path>` (e.g. test-name=`MyTestFunc` path=`./api`)
* To run tests in debug mode: `go test -debug` (in case is passing a path, debug flag must be the last argument)

## Linting
Please make sure your code changes pass all linting rules. To run lint:

```bash
make lint
```

## Pull Requeests
- Please open an Issue stating intent with all PRs
- Submit a PR for review
- Your name will be added below to Contributors
