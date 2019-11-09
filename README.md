## corn

This is used to generate a project with the followings:

- Echo framework - https://echo.labstack.com
- SQLBoiler database - https://github.com/volatiletech/sqlboiler

## Prerequisites

- PostgreSQL or MySQL if you choose to use a database.
- Go programming language, version 1.3.x or newer with go mod
- Ensure `$GOPATH/bin` is in your `$PATH`. Example: `PATH=$PATH:$GOPATH/bin`

## Installation

- `go get github.com/jameshuynh/corn`
- `corn new <project-name> --database={postgresql|mysql}`
- Start using it: `cd <project-name> && go run main.go`
- Visit http://localhost:3000
