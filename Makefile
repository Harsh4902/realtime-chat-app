PACKAGE=github.com/Harsh4902/realtime-chat-app
CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist

.PHONY: build-local
build-local:
	make build-app
	make build-client


.PHONY: build-app
build-app:
	go build -o ${DIST_DIR}/app ${PACKAGE}/app

.PHONY: build-client
build-client:
	go build -o ${DIST_DIR}/client ${PACKAGE}/cli-client
