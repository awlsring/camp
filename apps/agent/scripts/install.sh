#!/bin/bash

SRC_DIR=$(git rev-parse --show-toplevel)

cp $SRC_DIR/target/debug/geth-agent /opt/campd/campd