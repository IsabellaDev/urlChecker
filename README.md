# urlChecker

Go implementation of the url checking

# Versioning
This is the 0.2 version of this tool. 

# Usage

```bash
git clone https://github.com/isabellaliu77/urlChecker.git

cd urlchecker
go run urlChecker.go test/urls.txt
go run urlChecker.go test/urls.txt test/urls2.txt
```

# Features
- It can check the file passed through arguments to get and check the urls inside it. 

- Multiple files an be passed through arguments to get checked. 

- The result will show different colors and notifications based on different 
accessibility of each url. 

- The version of the tool can be checked through command line arguments, supporting both Windows and Unix style. (-v, --version, /v)

- Support for parallelization for the tool to run efficiently

- Only headers are requested when the urls are being checked. 
- Allows passing glob pattern as an argument
  ```go run urlChecker.go -g [glob-pattern]
  ```
  For example,
  ```go run urlChecker.go -g *.txt
  ```

- Allows passing -j, --json flag, which causes the program to output JSON. The JSON output should look like [ { "url": 'https://www.google.com', "status": 200 }, { "url": 'https://bad-link.com', "status": 404 } ], and allows output be piped into a file, eg.: output.txt

```go run urlChecker.go test/urls2.txt -j
```
```go run urlChecker.go test/urls2.txt --json > output.txt
```