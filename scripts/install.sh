#!/bin/bash

set -e

SERVICE_NAME="megaraid-exporter"
BINARY_PATH="/usr/local/bin/megaraid-exporter"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
CONFIG_DIR="/etc/megaraid-exporter"

echo "Installing MegaRAID Exporter..."

# Copy binary
sudo cp bin/megaraid-exporter ${BINARY_PATH}
sudo chmod +x ${BINARY_PATH}

# Create config directory
sudo mkdir -p ${CONFIG_DIR}
sudo cp config/config.yaml ${CONFIG_DIR}/ 2>/dev/null || true

# Create systemd service
sudo tee ${SERVICE_FILE} > /dev/null <<EOF
[Unit]
Description=MegaRAID Prometheus Exporter
After=network.target
Wants=network.target

[Service]
Type=simple
User=root
ExecStart=${BINARY_PATH} --config=${CONFIG_DIR}/config.yaml
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and enable service
sudo systemctl daemon-reload
sudo systemctl enable ${SERVICE_NAME}

echo "Installation completed!"
echo "Start service: sudo systemctl start ${SERVICE_NAME}"
echo "Check status: sudo systemctl status ${SERVICE_NAME}"
