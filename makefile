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


LATEST_LOG_DIR := ./logs/latest
OLD_LOG_DIR := ./logs/old

start:
	@echo "\n🚀 Starting servers...\ncancel: 'ctrl + c'\n"
	@mkdir -p $(LATEST_LOG_DIR)
	@mkdir -p $(OLD_LOG_DIR)
	$(eval NOW := $(shell date --utc --iso-8601=seconds))
	@[ -n "$(wildcard $(LATEST_LOG_DIR)/*)" ] && mv $(LATEST_LOG_DIR)/* $(OLD_LOG_DIR)/ || true
	@touch $(LATEST_LOG_DIR)/database-$(NOW).log
	@touch $(LATEST_LOG_DIR)/auth-$(NOW).log
	@touch $(LATEST_LOG_DIR)/resource-$(NOW).log
	@touch $(LATEST_LOG_DIR)/ui-$(NOW).log
	@trap 'kill -TERM $$(jobs -p)' INT; \
	make drun > $(LATEST_LOG_DIR)/database-$(NOW).log & \
	make arun > $(LATEST_LOG_DIR)/auth-$(NOW).log & \
	make srun > $(LATEST_LOG_DIR)/resource-$(NOW).log & \
	make uirun > $(LATEST_LOG_DIR)/ui-$(NOW).log;

rm-old-logs:
	rm -rf $(OLD_LOG_DIR)/*

test:
	go test -race ./...
