# link-tracker

This repository serves as a starting point for building web services.
It uses minimal dependencies, implements idiomatic code and follows as much as possible the following best practices:

- https://golang.org/doc/effective_go
- https://github.com/golang/go/wiki/CodeReviewComments
- https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html

It also includes a makefile whose default target run the following tools:

- go mod tidy
- go fmt
- goimports
- go test

## Build locally
After cloning the repository, execute the following command under the root directory

```shell
make run
```

## Create a link

`curl -POST http://localhost:8080/link -d '{"link":"http://www.google.com", "password":"123"}'`

## Open a link
You can either use cURL and follow the redirection with -L or opening a browser and navigate to the link. For a link with id 1 please visit: http://localhost:8080/link/1?password=123
