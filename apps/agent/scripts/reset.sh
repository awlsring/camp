#!/bin/bash

sudo systemctl daemon-reload
sudo systemctl enable campd.service
sudo systemctl start campd.service