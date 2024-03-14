genbuf:
	buf lint && buf generate

run:
	go run cmd/server/main.go

drun:
	go run cmd/server/database/main.go

arun: 
	go run cmd/server/auth/main.go

srun: 
	go run cmd/server/auth/main.go

test:
	go test -race ./...
