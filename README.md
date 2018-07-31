# go-mana-test
The Make Apis Nice Again. Testing framework

[![Go Report Card](https://goreportcard.com/badge/github.com/mattrmiller/go-mana-test)](https://goreportcard.com/report/github.com/mattrmiller/go-mana-tet)
[![godocs](https://img.shields.io/badge/godocs-reference-blue.svg)](https://godoc.org/github.com/mattrmiller/go-mana-test)

# Mission Statement
The Make Apis Nice Again. This is a testing framework to help you quickly, and neatly write end to end tests for your Api.


# Install
```
go get github.com/mattrmiller/go-mana-test
```

```
import (
    bedrock "github.com/mattrmiller/go-mana-test"
)
```

# Rules For Contributing
- Please make sure all changed files are run through gofmt
- Submit a PR for review
- Your name will be added below to Contributors

# Helpful Bash Additions
Test recursively, all go tests, in current folder:
```
alias go-testall='go test $(go list ./... | grep -v vendor)'
```

# Author
[Matthew R. Miller](https://github.com/mattrmiller)

# Contributors
[Matthew R. Miller](https://github.com/mattrmiller)

# License
[MIT License](LICENSE)
