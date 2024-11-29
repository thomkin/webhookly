#!/bin/bash

set -e

BUILD_DIR="$(pwd)"
EXECUTABLE="$BUILD_DIR/webhookly"

# build the current go code
go build -o $EXECUTABLE .

# create a service file under /etc/systemd
SERVICE_FILE="/etc/systemd/system/webhookly.service"
echo "[Unit]
Description=Webhookly
After=network.target

[Service]
User=webhookly
ExecStart=$EXECUTABLE
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target" > $SERVICE_FILE

# enable and start the service
systemctl enable webhookly
systemctl start webhookly
