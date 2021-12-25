.DEFAULT:
build: test lint
	goreleaser release --snapshot --rm-dist --skip-publish

test:
	go test --cover -covermode=count -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./... -v

format:
	go fmt ./...

format-check:
	gofmt -l cmd/.. .