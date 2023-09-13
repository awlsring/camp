SRC_DIR := $(shell git rev-parse --show-toplevel)
PKG_DIR := $(SRC_DIR)/packages
APP_DIR := $(SRC_DIR)/apps
CUR_DIR := $(shell pwd)
OPEN_API_MODEL := $(SRC_DIR)/models/local/build/smithy/openapi/openapi/CampLocal.openapi.json
SDK_DST := $(PKG_DIR)/camp_local

all: build

codegen:
	@echo "Building the local model"
	cd ${SRC_DIR}/models/local && smithy format model && smithy build
	@echo "Generating the local api server and client"
	mkdir -p $(SDK_DST)
	go mod download golang.org/x/tools # tpm hack, fix this later
	go run github.com/ogen-go/ogen/cmd/ogen --target $(SDK_DST) -package camplocal --clean $(OPEN_API_MODEL)
	@echo "Creating static docs"
	mkdir -p $(APP_DIR)/local/swagger
	cp $(OPEN_API_MODEL) $(APP_DIR)/local/swagger/swagger.json
	go generate $(SRC_DIR)/models/local/database

build: codegen
	@echo "Building Camp Local app"
	go mod tidy
	go build -o $(CUR_DIR)/build/camplocal $(CUR_DIR)/main.go

.PHONY: all