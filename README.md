# link-tracker

This repository serves as a starting point for building web services. It uses minimal dependencies, implements idiomatic
code and follows as much as possible the following best practices:

- https://golang.org/doc/effective_go
- https://github.com/golang/go/wiki/CodeReviewComments
- https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html

It also includes a makefile whose default target run the following tools:

- `go mod tidy`
- `go fmt`
- `goimports`
- `go test`

## Go installation

- Download Go tools from `https://golang.org/dl`
- If you are on a Unix-like system using bash, add the following lines to your .profile:
    - `export GOPATH=$HOME/go`
    - `export PATH=$PATH:$GOPATH/bin`

Note: Despite being mentioned in many blog posts, there is no need to define GOROOT.

## Build locally

After cloning the repository, execute the following command under the root directory

```shell
make run
```

> The main branch does not use any framework. It relies on simple primitives built on top of the standard library.
> [Here](https://github.com/emacampolo/link-tracker/tree/feature/gin-implementation) you will find the same program
> using the gin framework.
> It is worth noting that despite changing the entire transport implementation, the domain layer remains intact.

## Business Problem

This project provides a web service for creating and redirecting URLs upon password validation. It also tracks how many
times a URL is visited with a simple counter. As implementation note, the service is not thread safe, and the repository
is volatile, meaning that the default implementation stores the data in memory.

## Create a link

`curl -POST http://localhost:8080/link -d '{"link":"https://www.google.com", "password":"123"}'`

## Open a link

You can either use cURL and follow the redirection with -L or opening a browser and navigate to the link. For a link
with id 1 please visit: http://localhost:8080/link/1?password=123
