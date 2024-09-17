# casestudy

Simple binary that will parse some HTML pages and extract all links referenced on those pages.

## Install

Just clone the repository and build the binary

```
~/ git clone https://github.com/nicolasverle/casestudy.git
# ...
~/ cd casestudy/
~/casestudy/ make release
go fmt ./...
go vet ./...
CGO_ENABLED=0 go build -ldflags "-s -w" -a -trimpath  -o bin/linkextractor main.go
~/casestudy/ ./bin/linkextractor --help
Command that will parse a set URLs and extract all the links from their contents

Usage:
  linkextractor [flags]

Flags:
  -h, --help            help for linkextractor
  -o, --output string   format of the output (default "stdout")
```
