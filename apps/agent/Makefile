SRC_DIR := $(shell git rev-parse --show-toplevel)
PKG_DIR := $(SRC_DIR)/packages
CUR_DIR := $(shell pwd)
GO_SDK_DST := $(PKG_DIR)/camp_agent
OPEN_API_MODEL := $(SRC_DIR)/models/agent/build/smithy/openapi/openapi/CampAgent.openapi.json
SERVER_SDK_DST := $(PKG_DIR)/camp_agent_server
CLIENT_SDK_DST := $(PKG_DIR)/camp_agent_client
SERVER_SDK_SRC := $(SRC_DIR)/models/agent/build/smithy/rust/rust-server-codegen
CLIENT_SDK_SRC := $(SRC_DIR)/models/agent/build/smithy/rust/rust-client-codegen

all: build

codegen:
	@echo "Building the agent model"
	cd ${SRC_DIR}/models/agent && smithy build

	@echo "Creating rs package location"
	mkdir -p $(SERVER_SDK_DST) $(CLIENT_SDK_DST)
	cp -av $(SERVER_SDK_SRC)/* $(SERVER_SDK_DST)/
	cp -av $(CLIENT_SDK_SRC)/* $(CLIENT_SDK_DST)/

	@echo "Generating agent go client via ogen"
	mkdir -p $(GO_SDK_DST)
	go mod download golang.org/x/mod # tpm hack, fix this later
	go mod download golang.org/x/tools # tpm hack, fix this later
	go run github.com/ogen-go/ogen/cmd/ogen --target $(GO_SDK_DST) -package campagent --clean $(OPEN_API_MODEL)

build: codegen clippy
	cargo build

install: build
	${SRC_DIR}/apps/agent/scripts/install.sh
	${SRC_DIR}/apps/agent/scripts/reset.sh

run: codegen
	cargo run

clippy: codegen
	cargo clippy

test: codegen
	cargo test

doc-open: codegen
	cargo doc --no-deps --open

clean:
	cargo clean || echo "Unable to run cargo clean"

distclean: clean
	rm -rf $(SERVER_SDK_DST) $(CLIENT_SDK_DST) Cargo.lock

.PHONY: all