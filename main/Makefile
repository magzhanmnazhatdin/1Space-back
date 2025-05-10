# Makefile for Go project testing

.PHONY: test unit integration clean cover

test:
	go test -v ./...

unit:
	go test -v ./test/unit/...

integration:
	go test -v ./test/integration/...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in your browser to view test coverage report"

clean:
	rm -f coverage.out coverage.html
