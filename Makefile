SRC_DIR := $(shell git rev-parse --show-toplevel)
APP_DIR := $(SRC_DIR)/apps
CUR_DIR := $(shell pwd)
MODEL_DIR := $(SRC_DIR)/models

APPS := agent local

all: $(APPS)

codegen-common:
	@echo "Building the common model"
	cd ${MODEL_DIR}/common && smithy format model && smithy build

build-common:
	mkdir -p $(SRC_DIR)/build
	go mod tidy

$(APPS): codegen-common
	@echo "Building $@"
	$(MAKE) -C apps/$@

.PHONY: all