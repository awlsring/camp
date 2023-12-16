SRC_DIR := $(shell git rev-parse --show-toplevel)
MODEL_DIR := $(SRC_DIR)/model

APPS := local campd

all: $(APPS)

codegen-common:
	@echo "Building the common model"
	cd ${MODEL_DIR}/common && smithy format model && smithy build

build-common:
	mkdir -p $(SRC_DIR)/build
	go mod tidy

$(APPS): codegen-common
	@echo "Building $@"
	$(MAKE) -C cmd/$@

.PHONY: all