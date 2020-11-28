# Contributing to urlChecker

## Welcome!

Thank you for considering working on this project, and this is where you can get started. 

## Environment Setup

urlChecker is a small and simple program helping checking dead links, and it is written in Golang, please make sure you have the newest version of Golang installed on your computer.

## Usage

```bash
git clone https://github.com/isabellaliu77/urlChecker.git

cd urlchecker
go run urlChecker.go test/urls.txt
go run urlChecker.go test/urls.txt test/urls2.txt
```

##  Formatter

Please use the command below to format any code changes before sending any PR: 

``
gofmt -w filename.go
``
## Linter

Please use the command below to check linters before sending any PR: 
``
golint filename.go

## Testing

To test the program `go test -v`

## Code coverage

To check the code coverage ` go test -cover`