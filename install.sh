#!/bin/bash

set -e

BUILD_DIR="$(pwd)"
EXECUTABLE="$BUILD_DIR/webhookly"
CONFIG_FILE="/etc/webhookly/config.yaml"

CONFIG_DIR="/etc/webhookly"

if [ ! -d "$CONFIG_DIR" ]; then
  mkdir -p "$CONFIG_DIR"
fi

if [ ! -f "$CONFIG_FILE" ]; then
  echo "secret: 1234
cert: /etc/ssl/star.pem
key: /etc/ssl/star.key
handlers:
  refs/heads/main:
    path: /tmp/execute.me.sh" > $CONFIG_FILE
fi

# build the current go code
go build -o $EXECUTABLE .

# create a service file under /etc/systemd
SERVICE_FILE="/etc/systemd/system/webhookly.service"
echo "[Unit]
Description=Webhookly
After=network.target

[Service]
User=webhookly
ExecStart=$EXECUTABLE -config $CONFIG_FILE
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target" > $SERVICE_FILE

# enable and start the service
systemctl enable webhookly
systemctl start webhookly

