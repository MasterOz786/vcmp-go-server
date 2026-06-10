.PHONY: test tidy

test:
	go test ./...

tidy:
	go mod tidy
