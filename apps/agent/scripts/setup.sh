#!/bin/bash

SRC_DIR=$(git rev-parse --show-toplevel)

sudo groupadd campd
sudo useradd -r -g campd -s /sbin/nologin -d /etc/campd campd
sudo chown -R campd:campd /etc/campd/

mkdir -p /opt/campd
sudo chown -R campd:campd /opt/campd
sudo chmod 755 /opt/campd
cp $SRC_DIR/apps/agent/resources/campd.service /etc/systemd/system/campd.service