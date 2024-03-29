## corn

This is used to generate a project with the followings:

- Echo framework - https://echo.labstack.com
- SQLBoiler database - https://github.com/volatiletech/sqlboiler

## Prerequisites

- PostgreSQL database (MySQL will be added later)
- Go programming language, version 1.3.x or newer with go mod
- Ensure `$GOPATH/bin` is in your `$PATH`. Example: `PATH=$PATH:$GOPATH/bin`

## Installation

- `go get github.com/jameshuynh/corn`

## Usage

### 1. Create a new app

```bash
corn new <project_name> --database=postgresql
```

After that you can run this project by executing:

```bash
go run main.go
```

And visit http://localhost:3000

### 2. Create database

```bash
corn db:create
```

### 3. Drop database

```bash
corn db:drop
```

### 4. Generate a new model

```bash
corn g model post title description:text
```

### 5. Generate a blank new migration file

```bash
corn g migration create_cars
```

### 6. Run migration

```bash
corn db:migrate
```

### 7. Rollback a migration

```bash
corn db:rollback
```
