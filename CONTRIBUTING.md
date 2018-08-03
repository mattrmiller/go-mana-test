# Contributing

## Setup
1. Clone this project with the Go tooling
`go get -u github.com/mattrmiller/go-mana-test`

## Adding New Go Packages
We use [dep](https://github.com/golang/dep/) as our dependency manager.

This is a service, so in order to produce reproducible builds we check in our `vendor` directory.

Adding a new package:
```
dep ensure -add github.com/foo/bar
```

From there, when you need to update to a new version. Just update to [Gopkg.toml](./Gopkg.toml) file following the [version schema](https://github.com/golang/dep/blob/master/docs/Gopkg.toml.md).
```
dep ensure -update
```

## Tests
Please provide tests for any core logic changes. Such as the API routes.

* Test recursively, all go tests, in current folder:
```
alias go-testall='go test $(go list ./... | grep -v vendor)'
```
* To run only tests from a package: `go test <path>`
* To run only a single test: `go test -run <test-name> <path>` (e.g. test-name=`MyTestFunc` path=`./api`)
* To run tests in debug mode: `go test -debug` (in case is passing a path, debug flag must be the last argument)

## Pull Requeests
- Please make sure all changed files are run through gofmt
- Submit a PR for review
- Your name will be added below to Contributors
