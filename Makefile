TODO_LIST_BINARY=todoListApp
IMAGE_NAME=todolist
IMAGE_TAG=latest

# run: as running app
run:
	go run main.go

## test: Run all test in this app
test:
	@echo "All tests are running..."
	go test -v ./...
	@echo "Test finished"

## test: Run all test with clean cache in this app
test_nocache:
	@echo "Clean all cache..."
	go clean -testcache
	@echo "All tests are running..."
	go test -v ./...
	@echo "Test finished"

## test_cover: Run all test with coverage
test_cover:
	@echo "All test are running with coverage..."
	go test ./... -v -cover

## test: Run all test with clean cache and coverage
test_cover_nocache:
	@echo "Clean all cache..."
	go clean -testcache
	@echo "All tests are running..."
	go test ./... -v -cover
	@echo "Test finished"

# build: build app todo-list to binary file
build:
	@echo "Building binary todo list.."
	env GOOS=linux CGO_ENABLED=0 go build -o ${TODO_LIST_BINARY} ./
	@echo "Done!"
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"