```markdown
# Production Setup Guide

This guide provides step-by-step instructions for deploying MegaRAID Exporter in production environments.

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Pre-Installation Steps](#pre-installation-steps)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Service Setup](#service-setup)
6. [Security Configuration](#security-configuration)
7. [Monitoring Setup](#monitoring-setup)
8. [Testing & Validation](#testing--validation)
9. [Maintenance](#maintenance)
10. [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements
- **Operating System**: Linux (Ubuntu 18.04+, RHEL 7+, CentOS 7+, SLES 12+)
- **CPU**: 1 core minimum, 2+ cores recommended
- **Memory**: 512MB minimum, 1GB+ recommended
- **Disk Space**: 1GB minimum for installation and logs
- **Network**: Access to target servers and monitoring infrastructure

### Hardware Requirements
- LSI/Broadcom MegaRAID controller
- Active RAID arrays (recommended for full functionality)
- Root/sudo access for hardware queries

### Software Dependencies
```bash
# Check if required commands are available
which curl || echo "Install curl: sudo apt-get install curl"
which systemctl || echo "Systemd is required"
which sudo || echo "Sudo is required"
```

## Pre-Installation Steps

### Step 1: System Updates
```bash
# Update package lists
sudo apt-get update

# Install essential packages
sudo apt-get install -y curl wget systemd sudo

# For RHEL/CentOS
# sudo yum update -y
# sudo yum install -y curl wget systemd sudo
```

### Step 2: Install MegaCLI64
```bash
# Option A: Download from Broadcom (recommended)
cd /tmp
wget https://docs.broadcom.com/docs-and-downloads/raid-controllers/raid-controllers-common-files/8-07-14_MegaCLI.zip

# Extract and install
unzip 8-07-14_MegaCLI.zip
cd Linux
sudo dpkg -i MegaCli-8.07.14-1.noarch.deb

# For RHEL/CentOS
# sudo rpm -ivh MegaCli-8.07.14-1.noarch.rpm

# Option B: Package manager (if available)
# sudo apt-get install megacli
```

### Step 3: Verify MegaCLI Installation
```bash
# Test MegaCLI installation
sudo megacli64 -v
sudo megacli64 -adpCount -NoLog

# Expected output should show controller count
# Controllers found: 1
```

### Step 4: Create System Users (Optional but Recommended)
```bash
# Create dedicated user for the exporter
sudo useradd -r -s /bin/false -d /var/lib/megaraid-exporter megaraid-exporter

# Create necessary directories
sudo mkdir -p /var/lib/megaraid-exporter
sudo mkdir -p /var/log/megaraid-exporter
sudo chown megaraid-exporter:megaraid-exporter /var/lib/megaraid-exporter
sudo chown megaraid-exporter:megaraid-exporter /var/log/megaraid-exporter
```

## Installation

### Step 1: Download MegaRAID Exporter
```bash
# Create installation directory
cd /opt
sudo mkdir -p megaraid-exporter
cd megaraid-exporter

# Download latest release
LATEST_VERSION=$(curl -s https://api.github.com/repos/yourusername/megaraid-exporter/releases/latest | grep tag_name | cut -d '"' -f 4)
echo "Latest version: $LATEST_VERSION"

# Download binary
sudo wget https://github.com/yourusername/megaraid-exporter/releases/download/${LATEST_VERSION}/megaraid-exporter-linux-amd64

# Make executable and move to system path
sudo chmod +x megaraid-exporter-linux-amd64
sudo mv megaraid-exporter-linux-amd64 /usr/local/bin/megaraid-exporter

# Verify installation
/usr/local/bin/megaraid-exporter version
```

### Step 2: Alternative - Build from Source
```bash
# Install Go (if not already installed)
cd /tmp
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Clone and build
git clone https://github.com/yourusername/megaraid-exporter.git
cd megaraid-exporter
make build
sudo cp megaraid-exporter /usr/local/bin/
```

## Configuration

### Step 1: Create Configuration Directory
```bash
sudo mkdir -p /etc/megaraid-exporter
sudo mkdir -p /var/log/megaraid-exporter
```

### Step 2: Create Main Configuration File
```bash
sudo tee /etc/megaraid-exporter/config.yaml << 'EOF'
# MegaRAID Exporter Production Configuration

# Server settings
port: 9272
log_level: "info"

# MegaCLI settings
megacli_path: "/usr/sbin/megacli64"
command_timeout: 30s

# Metrics collection
metrics:
  collect_interval: 30s
  enabled_collectors:
    - controller
    - arrays
    - drives
    - bbu
    - events

# Event monitoring
events:
  lookback_hours: 24
  severity_levels:
    - "critical"
    - "warning"

# Performance tuning
advanced:
  max_concurrent_commands: 2
  cache_duration: 15s
  skip_drive_states:
    - "Unconfigured(bad)"
EOF
```

### Step 3: Set Configuration Permissions
```bash
sudo chmod 640 /etc/megaraid-exporter/config.yaml
sudo chown root:root /etc/megaraid-exporter/config.yaml
```

### Step 4: Test Configuration
```bash
# Test the configuration
sudo /usr/local/bin/megaraid-exporter --config /etc/megaraid-exporter/config.yaml --log-level debug &
sleep 5

# Test metrics endpoint
curl -s http://localhost:9272/metrics | head -10

# Stop test process
sudo pkill megaraid-exporter
```

## Service Setup

### Step 1: Create Systemd Service File
```bash
sudo tee /etc/systemd/system/megaraid-exporter.service << 'EOF'
[Unit]
Description=MegaRAID Prometheus Exporter
Documentation=https://github.com/yourusername/megaraid-exporter
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
ExecStart=/usr/local/bin/megaraid-exporter --config /etc/megaraid-exporter/config.yaml
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=megaraid-exporter

# Environment
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/megaraid-exporter /tmp

# Resource limits
LimitNOFILE=2048
LimitNPROC=1024
MemoryLimit=1G
CPUQuota=200%

[Install]
WantedBy=multi-user.target
EOF
```

### Step 2: Enable and Start Service
```bash
# Reload systemd configuration
sudo systemctl daemon-reload

# Enable service to start at boot
sudo systemctl enable megaraid-exporter

# Start the service
sudo systemctl start megaraid-exporter

# Check service status
sudo systemctl status megaraid-exporter
```

### Step 3: Verify Service is Working
```bash
# Check if service is active
sudo systemctl is-active megaraid-exporter

# Check service logs
sudo journalctl -u megaraid-exporter -f --lines 50

# Test metrics endpoint
curl -s http://localhost:9272/health
curl -s http://localhost:9272/metrics | grep megaraid_controller_info
```

## Security Configuration

### Step 1: Firewall Configuration
```bash
# Using UFW (Ubuntu)
sudo ufw allow from 10.0.0.0/8 to any port 9272 comment 'MegaRAID Exporter - Internal'
sudo ufw allow from 192.168.0.0/16 to any port 9272 comment 'MegaRAID Exporter - Private'

# Using firewalld (RHEL/CentOS)
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="10.0.0.0/8" port protocol="tcp" port="9272" accept'
sudo firewall-cmd --reload

# Using iptables (manual)
sudo iptables -A INPUT -s 10.0.0.0/8 -p tcp --dport 9272 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 9272 -j DROP
```

### Step 2: Log Rotation Setup
```bash
sudo tee /etc/logrotate.d/megaraid-exporter << 'EOF'
/var/log/megaraid-exporter/*.log {
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
```

### Step 3: File Permissions Audit
```bash
# Verify file permissions
ls -la /usr/local/bin/megaraid-exporter
ls -la /etc/megaraid-exporter/
ls -la /var/log/megaraid-exporter/

# Expected permissions:
# /usr/local/bin/megaraid-exporter: -rwxr-xr-x root root
# /etc/megaraid-exporter/config.yaml: -rw-r----- root root
```

## Monitoring Setup

### Step 1: Create Health Check Script
```bash
sudo tee /usr/local/bin/megaraid-health-check << 'EOF'
#!/bin/bash
# MegaRAID Exporter Health Check

URL="http://localhost:9272"
LOG_FILE="/var/log/megaraid-exporter/health-check.log"

# Function to log with timestamp
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> $LOG_FILE
}

# Check service status
if ! systemctl is-active --quiet megaraid-exporter; then
    log "ERROR: Service is not running"
    exit 1
fi

# Check HTTP endpoint
if ! curl -s --connect-timeout 5 "$URL/health" > /dev/null; then
    log "ERROR: HTTP endpoint not responding"
    exit 1
fi

# Check for metrics
METRIC_COUNT=$(curl -s "$URL/metrics" | grep -c "megaraid_")
if [ "$METRIC_COUNT" -lt 5 ]; then
    log "WARNING: Low metric count: $METRIC_COUNT"
    exit 1
fi

log "OK: Health check passed, metrics: $METRIC_COUNT"
exit 0
EOF

sudo chmod +x /usr/local/bin/megaraid-health-check
```

### Step 2: Setup Cron for Health Checks
```bash
# Add cron job for health monitoring
(sudo crontab -l 2>/dev/null; echo "*/5 * * * * /usr/local/bin/megaraid-health-check") | sudo crontab -
```

### Step 3: Configure Prometheus (if using)
```bash
# Create prometheus configuration snippet
sudo tee /etc/prometheus/megaraid-exporter.yml << 'EOF'
# Add this to your main prometheus.yml scrape_configs section
- job_name: 'megaraid-exporter'
  static_configs:
    - targets: ['localhost:9272']
  scrape_interval: 30s
  scrape_timeout: 25s
  metrics_path: /metrics
EOF
```

## Testing & Validation

### Step 1: Functional Testing
```bash
# Test all endpoints
curl -I http://localhost:9272/
curl -I http://localhost:9272/health
curl -I http://localhost:9272/metrics

# Test specific metrics
curl -s http://localhost:9272/metrics | grep -E "(controller|array|drive|bbu)"

# Test with different log levels
sudo systemctl stop megaraid-exporter
sudo /usr/local/bin/megaraid-exporter --config /etc/megaraid-exporter/config.yaml --log-level debug &
sleep 5
curl -s http://localhost:9272/metrics | head -20
sudo pkill megaraid-exporter
sudo systemctl start megaraid-exporter
```

### Step 2: Performance Testing
```bash
# Test concurrent requests
for i in {1..10}; do
    curl -s http://localhost:9272/metrics > /dev/null &
done
wait

# Monitor resource usage
top -p $(pgrep megaraid-exporter)
```

### Step 3: Reliability Testing
```bash
# Test service restart
sudo systemctl restart megaraid-exporter
sleep 10
curl -s http://localhost:9272/health

# Test configuration reload
sudo systemctl reload megaraid-exporter
sleep 5
curl -s http://localhost:9272/health
```

## Maintenance

### Daily Tasks
```bash
# Check service status
sudo systemctl status megaraid-exporter

# Check recent logs
sudo journalctl -u megaraid-exporter --since "24 hours ago" | tail -20

# Verify metrics are being generated
curl -s http://localhost:9272/metrics | grep -c megaraid_
```

### Weekly Tasks
```bash
# Review log files
sudo tail -100 /var/log/megaraid-exporter/health-check.log

# Check disk space
df -h /var/log/megaraid-exporter/

# Update check (manual)
# Check https://github.com/yourusername/megaraid-exporter/releases for updates
```

### Monthly Tasks
```bash
# Security updates
sudo apt-get update && sudo apt-get upgrade

# Log cleanup (if not using logrotate)
sudo find /var/log/megaraid-exporter/ -name "*.log" -mtime +30 -delete

# Configuration backup
sudo tar czf /backup/megaraid-exporter-config-$(date +%Y%m%d).tar.gz /etc/megaraid-exporter/
```

## Troubleshooting

### Common Issues and Solutions

#### 1. Service Won't Start
```bash
# Check service status and logs
sudo systemctl status megaraid-exporter
sudo journalctl -u megaraid-exporter -n 50

# Common causes:
# - Configuration file syntax error
# - MegaCLI not found
# - Port already in use
# - Insufficient permissions
```

#### 2. No Metrics Returned
```bash
# Test MegaCLI directly
sudo megacli64 -adpCount -NoLog

# Check if controllers are detected
sudo megacli64 -AdpAllInfo -aALL -NoLog

# Verify configuration
sudo /usr/local/bin/megaraid-exporter --config /etc/megaraid-exporter/config.yaml --log-level debug
```

#### 3. High CPU Usage
```bash
# Check for too frequent collection
# Edit /etc/megaraid-exporter/config.yaml
# Increase collect_interval: 60s

# Check for too many concurrent commands
# Reduce max_concurrent_commands: 1

sudo systemctl restart megaraid-exporter
```

#### 4. Memory Issues
```bash
# Check memory usage
ps aux | grep megaraid-exporter

# Add memory limits to systemd service
# MemoryLimit=512M

# Enable caching to reduce command executions
# cache_duration: 30s
```

### Log Analysis Commands
```bash
# Service logs
sudo journalctl -u megaraid-exporter -f

# Error logs only
sudo journalctl -u megaraid-exporter -p err

# Performance logs
sudo journalctl -u megaraid-exporter --since "1 hour ago" | grep -i "slow\|timeout\|error"

# Health check logs
sudo tail -f /var/log/megaraid-exporter/health-check.log
```

### Emergency Procedures

#### Service Recovery
```bash
# If service is stuck
sudo systemctl kill megaraid-exporter
sudo systemctl reset-failed megaraid-exporter
sudo systemctl start megaraid-exporter

# If system is unresponsive
sudo killall -9 megaraid-exporter
sudo systemctl start megaraid-exporter
```

#### Configuration Rollback
```bash
# Restore from backup
sudo cp /backup/megaraid-exporter-config-YYYYMMDD.tar.gz /tmp/
cd /tmp && sudo tar xzf megaraid-exporter-config-YYYYMMDD.tar.gz
sudo cp etc/megaraid-exporter/config.yaml /etc/megaraid-exporter/
sudo systemctl restart megaraid-exporter
```

## Production Checklist

Before going live, ensure:

- [ ] MegaCLI64 installed and tested
- [ ] Service starts automatically at boot
- [ ] Firewall configured properly
- [ ] Log rotation configured
- [ ] Health monitoring setup
- [ ] Backup procedures in place
- [ ] Documentation updated
- [ ] Team trained on basic operations
- [ ] Emergency contacts identified
- [ ] Performance baselines established

## Support Contacts

- **Technical Issues**: Create GitHub issue with logs
- **Security Issues**: security@yourcompany.com
- **Production Support**: ops@yourcompany.com

---

**Next Steps**: After successful production deployment, consider setting up alerting rules and Grafana dashboards for comprehensive monitoring.
```
