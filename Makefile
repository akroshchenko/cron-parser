build:
	go build -o cron-parser

test:
	go test ./...

cover:
	go test ./... -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html
	open coverage.html || true

clean:
	@gitd clean -fdx &>/dev/null || echo "Cannot clean up" && exit 1

.PHONY: build test clean