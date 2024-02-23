genbuf:
	buf lint && buf generate

run:
	go run cmd/server/main.go

test:
	go test -race ./...
