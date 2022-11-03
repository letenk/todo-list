# todo-list
A RestApi Todo List for take home test in skyshi.com

## Documentation

[Postman Documentation](https://documenter.getpostman.com/view/12132212/2s8YRqmWJb)
[ERD Documentation](https://dbdiagram.io/d/635f77a35170fb6441c7f5f2)
[Image Docker Hub](https://hub.docker.com/repository/docker/letenk/todolist)

## Installation

1. Clone the project

```bash
  git clone https://github.com/letenk/todo-list.git
```

2. Go to the project directory

```go
  cd todo-list
```

3. Export environment variable into terminal

```go
export MYSQL_USER="root"
export MYSQL_PASSWORD="root"
export MYSQL_HOST="127.0.0.1"
export MYSQL_PORT="3306"
export MYSQL_DBNAME="todo4"
```
**Note: for value to each environment variable please customize with yours**

4. Start the server

```go
go run main.go
```

5. This app can be accessed in local with url: `http://localhost:3030`

## Run Test
Here can use `Makefile` for shortcut syntax to run each test.
**Note: If not export environment variable. Please export first!**
### Run All Test 
- Run with `Makefile`
```go
make test
```

- Run without `MakeFile`
```go
go test -v ./...
```

### Run All Test No Cache
- Run with `Makefile`
```go
make test_nocache
```

- Run without `MakeFile`
```go
go clean -testcache
go test -v ./...
```
### Run All Test With Code Coverage
- Run with `Makefile`
```go
make test_cover
```

- Run without `MakeFile`
```go
go test ./... -v -cover
```
### Run All Test With Code Coverage No Cache
- Run with `Makefile`
```go
make test_cover_nocache
```

- Run without `MakeFile`
```go
go clean -testcache
go test ./... -v -cover
```

