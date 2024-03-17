include .env.local

ENV_PATH=/workspaces/.env.local

genbuf:
	buf lint && buf generate

set-ui-env:
	ln -fn ./.env.local ./web/.env.local

run:
	go run cmd/server/main.go

drun:
	go run cmd/server/database/main.go -source ${ENV_PATH}

arun: 
	go run cmd/server/auth/main.go -source ${ENV_PATH}

srun: 
	go run cmd/server/resource/main.go -source ${ENV_PATH}

uirun:
	cd web && npx next start -p ${UI_SERVER_PORT}

uidev:
	cd web && npx next dev -p ${UI_SERVER_PORT}

buildui:
	cd web && npx next build

start:
	@echo "Starting servers..."
	make drun & \
	make arun & \
	make srun & \
	make uirun; 


test:
	go test -race ./...
