all: format build

ci: format build lint

format:
	go fmt ./...

build:
	go build -v .

build-examples:
	@for dir in examples/*; do \
		echo "Building $$dir"; \
		cd $$dir; \
		go build -v .; \
		cd -; \
	done

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install go.uber.org/nilaway/cmd/nilaway@latest

lint:
	golangci-lint run
	nilaway ./...

upgrade:
	go get -u
	go mod tidy