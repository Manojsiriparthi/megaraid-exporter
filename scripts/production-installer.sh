#!/bin/bash
# MegaRAID Exporter Production Installer
# This script automates the production installation process

set -e

# Configuration
BINARY_NAME="megaraid-exporter"
SERVICE_NAME="megaraid-exporter"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/megaraid-exporter"
LOG_DIR="/var/log/megaraid-exporter"
GITHUB_REPO="yourusername/megaraid-exporter"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

# Detect OS
detect_os() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        OS=$NAME
        VERSION=$VERSION_ID
    else
        log_error "Cannot detect operating system"
        exit 1
    fi
    
    log_info "Detected OS: $OS $VERSION"
}

# Install system dependencies
install_dependencies() {
    log_info "Installing system dependencies..."
    
    if [[ $OS == *"Ubuntu"* ]] || [[ $OS == *"Debian"* ]]; then
        apt-get update
        apt-get install -y curl wget systemd sudo unzip
    elif [[ $OS == *"CentOS"* ]] || [[ $OS == *"Red Hat"* ]]; then
        yum update -y
        yum install -y curl wget systemd sudo unzip
    else
        log_warning "Unknown OS, please install manually: curl, wget, systemd, sudo, unzip"
    fi
    
    log_success "System dependencies installed"
}

# Install MegaCLI
install_megacli() {
    log_info "Checking MegaCLI installation..."
    
    if command -v megacli64 &> /dev/null; then
        log_success "MegaCLI64 is already installed"
        megacli64 -v
        return
    fi
    
    log_info "Installing MegaCLI64..."
    
    cd /tmp
    
    # Download MegaCLI (you may need to update this URL)
    if [[ ! -f "8-07-14_MegaCLI.zip" ]]; then
        log_info "Please download MegaCLI from Broadcom website manually"
        log_info "URL: https://docs.broadcom.com/docs-and-downloads/raid-controllers/"
        log_error "MegaCLI installation required before proceeding"
        exit 1
    fi
    
    unzip -o 8-07-14_MegaCLI.zip
    
    if [[ $OS == *"Ubuntu"* ]] || [[ $OS == *"Debian"* ]]; then
        dpkg -i Linux/MegaCli-8.07.14-1.noarch.deb || true
        apt-get install -f -y
    elif [[ $OS == *"CentOS"* ]] || [[ $OS == *"Red Hat"* ]]; then
        rpm -ivh Linux/MegaCli-8.07.14-1.noarch.rpm || true
    fi
    
    # Verify installation
    if command -v megacli64 &> /dev/null; then
        log_success "MegaCLI64 installed successfully"
        megacli64 -v
    else
        log_error "MegaCLI64 installation failed"
        exit 1
    fi
}

# Download and install megaraid-exporter
install_exporter() {
    log_info "Installing MegaRAID Exporter..."
    
    # Get latest release
    LATEST_VERSION=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | cut -d'"' -f4)
    
    if [[ -z "$LATEST_VERSION" ]]; then
        log_error "Failed to fetch latest version"
        exit 1
    fi
    
    log_info "Latest version: $LATEST_VERSION"
    
    # Download binary
    DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/megaraid-exporter-linux-amd64"
    
    cd /tmp
    wget -O "$BINARY_NAME" "$DOWNLOAD_URL"
    
    # Install binary
    chmod +x "$BINARY_NAME"
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    
    # Verify installation
    if [[ -f "$INSTALL_DIR/$BINARY_NAME" ]]; then
        log_success "Binary installed to $INSTALL_DIR/$BINARY_NAME"
        "$INSTALL_DIR/$BINARY_NAME" version
    else
        log_error "Binary installation failed"
        exit 1
    fi
}

# Create configuration
create_configuration() {
    log_info "Creating configuration..."
    
    # Create directories
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$LOG_DIR"
    
    # Create configuration file
    cat > "$CONFIG_DIR/config.yaml" << EOF
# MegaRAID Exporter Production Configuration
port: 9272
log_level: "info"
megacli_path: "/usr/sbin/megacli64"
command_timeout: 30s

metrics:
  collect_interval: 30s
  enabled_collectors:
    - controller
    - arrays
    - drives
    - bbu
    - events

events:
  lookback_hours: 24
  severity_levels:
    - "critical"
    - "warning"

advanced:
  max_concurrent_commands: 2
  cache_duration: 15s
EOF
    
    # Set permissions
    chmod 640 "$CONFIG_DIR/config.yaml"
    chown root:root "$CONFIG_DIR/config.yaml"
    
    log_success "Configuration created at $CONFIG_DIR/config.yaml"
}

# Create systemd service
create_service() {
    log_info "Creating systemd service..."
    
    cat > "/etc/systemd/system/$SERVICE_NAME.service" << EOF
[Unit]
Description=MegaRAID Prometheus Exporter
Documentation=https://github.com/$GITHUB_REPO
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
ExecStart=$INSTALL_DIR/$BINARY_NAME --config $CONFIG_DIR/config.yaml
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$SERVICE_NAME

# Environment
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$LOG_DIR /tmp

# Resource limits
LimitNOFILE=2048
LimitNPROC=1024
MemoryLimit=1G
CPUQuota=200%

[Install]
WantedBy=multi-user.target
EOF
    
    # Reload systemd and enable service
    systemctl daemon-reload
    systemctl enable "$SERVICE_NAME"
    
    log_success "Systemd service created and enabled"
}

# Setup monitoring
setup_monitoring() {
    log_info "Setting up monitoring..."
    
    # Create health check script
    cat > "/usr/local/bin/megaraid-health-check" << 'EOF'
#!/bin/bash
URL="http://localhost:9272"
LOG_FILE="/var/log/megaraid-exporter/health-check.log"

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> $LOG_FILE
}

if ! systemctl is-active --quiet megaraid-exporter; then
    log "ERROR: Service is not running"
    exit 1
fi

if ! curl -s --connect-timeout 5 "$URL/health" > /dev/null; then
    log "ERROR: HTTP endpoint not responding"
    exit 1
fi

METRIC_COUNT=$(curl -s "$URL/metrics" | grep -c "megaraid_")
if [ "$METRIC_COUNT" -lt 5 ]; then
    log "WARNING: Low metric count: $METRIC_COUNT"
    exit 1
fi

log "OK: Health check passed, metrics: $METRIC_COUNT"
exit 0
EOF
    
    chmod +x "/usr/local/bin/megaraid-health-check"
    
    # Add cron job
    (crontab -l 2>/dev/null; echo "*/5 * * * * /usr/local/bin/megaraid-health-check") | crontab -
    
    # Setup log rotation
    cat > "/etc/logrotate.d/megaraid-exporter" << EOF
$LOG_DIR/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0644 root root
    postrotate
        systemctl reload megaraid-exporter 2>/dev/null || true
    endscript
}
EOF
    
    log_success "Monitoring setup completed"
}

# Test installation
test_installation() {
    log_info "Testing installation..."
    
    # Start service
    systemctl start "$SERVICE_NAME"
    sleep 5
    
    # Check service status
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_success "Service is running"
    else
        log_error "Service failed to start"
        systemctl status "$SERVICE_NAME"
        exit 1
    fi
    
    # Test endpoints
    if curl -s --connect-timeout 10 "http://localhost:9272/health" > /dev/null; then
        log_success "Health endpoint is responding"
    else
        log_error "Health endpoint is not responding"
        exit 1
    fi
    
    # Test metrics
    METRIC_COUNT=$(curl -s "http://localhost:9272/metrics" | grep -c "megaraid_" || echo "0")
    if [[ $METRIC_COUNT -gt 0 ]]; then
        log_success "Metrics endpoint is working ($METRIC_COUNT metrics)"
    else
        log_warning "No metrics found - check MegaRAID controller presence"
    fi
    
    log_success "Installation test completed"
}

# Main installation function
main() {
    echo "=========================================="
    echo "MegaRAID Exporter Production Installer"
    echo "=========================================="
    
    check_root
    detect_os
    install_dependencies
    install_megacli
    install_exporter
    create_configuration
    create_service
    setup_monitoring
    test_installation
    
    echo
    log_success "Installation completed successfully!"
    echo
    echo "Next steps:"
    echo "1. Review configuration: $CONFIG_DIR/config.yaml"
    echo "2. Check service status: systemctl status $SERVICE_NAME"
    echo "3. View logs: journalctl -u $SERVICE_NAME -f"
    echo "4. Test metrics: curl http://localhost:9272/metrics"
    echo "5. Configure firewall and monitoring as needed"
    echo
    echo "Service management commands:"
    echo "  Start:   systemctl start $SERVICE_NAME"
    echo "  Stop:    systemctl stop $SERVICE_NAME"
    echo "  Restart: systemctl restart $SERVICE_NAME"
    echo "  Status:  systemctl status $SERVICE_NAME"
    echo "  Logs:    journalctl -u $SERVICE_NAME -f"
}

# Run main function
main "$@"
