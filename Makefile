PACKAGES_PATH = $(shell go list -f '{{ .Dir }}' ./...)
APPLICATION_PATH = "$(GOPATH)/github.com/johanaggu/unittest"

.PHONY: all
all: require tidy fmt goimports vet  mocks test

.PHONY: require
require:
	@type "goimports" > /dev/null 2>&1 \
		|| (echo 'goimports not found: to install it, run "go install golang.org/x/tools/cmd/goimports@latest"'; exit 1)
	@type "staticcheck" > /dev/null 2>&1 \
		|| (echo 'staticcheck not found: to install it, run "go install honnef.co/go/tools/cmd/staticcheck@latest"'; exit 1)

.PHONY: tidy
tidy:
	@echo "=> Executing go mod tidy"
	@go mod tidy

.PHONY: fmt
fmt:
	@echo "=> Executing go fmt"
	@go fmt ./...

.PHONY: goimports
goimports:
	@echo "=> Executing goimports"
	@goimports -w $(PACKAGES_PATH)

.PHONY: vet
vet:
	@echo "=> Executing go vet"
	@go vet ./...

.PHONY: staticcheck
staticcheck:
	@echo "=> Executing staticcheck"
	@staticcheck ./...

.PHONY: test
test:
	@echo "=> Running unit tests"
	@go test ./... -covermode=atomic -coverpkg=./... -count=1 -race -shuffle=on -short

.PHONY: test-cover
test-cover:
	@echo "=> Running unit tests and generating report"
	@go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1 -race -shuffle=on -short -v
	@go tool cover -html=/tmp/coverage.out

.PHONY: test-all
test-all:
	@echo "=> Running all tests"
	@docker-compose up -d --remove-orphans
	@go test ./... -covermode=atomic -coverpkg=./... -v -count=1 -race -shuffle=on;\
	exit_code=$$?;\
	docker-compose down -v;\
 	exit $$exit_code

.PHONY: test-all-cover
test-all-cover:
	@echo "=> Running all tests and generating report"
	@docker-compose up -d --remove-orphans
	@go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1 -race -shuffle=on;\
    exit_code=$$?;\
    docker-compose down -v;\
    if [ $$exit_code -ne 0 ]; then \
      	exit $$exit_code;\
    else \
       go tool cover -html=/tmp/coverage.out; \
    fi

.PHONY: up
up:
	@docker-compose up -d --remove-orphans -V

.PHONY: mocks
mocks:
	@echo "=> Creating all mocks"
	@go generate ./...

.PHONY: sqlc
sqlc:
	@echo "=> Executing sqlc generate "
	@sqlc generate

