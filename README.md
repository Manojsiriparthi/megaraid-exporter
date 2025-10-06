# MegaRAID Exporter

A standalone Prometheus exporter for LSI MegaRAID controllers that provides detailed metrics about RAID arrays, physical drives, and controller status using MegaCLI64.

## Features

- Monitor RAID array status and health
- Track physical drive metrics (temperature, errors, wear, predictive failures)
- Controller information and detailed statistics
- Battery backup unit (BBU) monitoring with charge cycles
- Event log monitoring for proactive alerts
- Foreign configuration detection
- **Standalone operation** - works without Prometheus installation
- Prometheus-compatible metrics format accessible via HTTP

## Prerequisites

- LSI MegaRAID controller with `megacli64` utility installed
- Go 1.19 or later (for building from source)
- Root/Administrator privileges (required for hardware access)

**Note**: Prometheus is NOT required to run this exporter. The exporter runs as a standalone HTTP server and exposes metrics that can be accessed directly or scraped by any monitoring system.

## Installation

### Install MegaCLI64

#### Option 1: Download from LSI/Broadcom
```bash
# Download MegaCLI from LSI/Broadcom website
# Extract and install
sudo dpkg -i megacli_*.deb  # For Debian/Ubuntu
# or
sudo rpm -ivh megacli-*.rpm  # For RHEL/CentOS
```

#### Option 2: Package Manager (if available)
```bash
# Ubuntu/Debian
sudo apt-get install megacli

# CentOS/RHEL (with EPEL)
sudo yum install megacli
```

### Install MegaRAID Exporter

#### Option 1: Download Binary
```bash
# Download latest release
wget https://github.com/yourusername/megaraid-exporter/releases/latest/download/megaraid-exporter-linux-amd64
chmod +x megaraid-exporter-linux-amd64
sudo mv megaraid-exporter-linux-amd64 /usr/local/bin/megaraid-exporter
```

#### Option 2: Build from Source
```bash
git clone https://github.com/yourusername/megaraid-exporter.git
cd megaraid-exporter
go build -o megaraid-exporter ./cmd/exporter
sudo mv megaraid-exporter /usr/local/bin/
```

## Usage

### Quick Start (Standalone)
```bash
# Run the exporter (no Prometheus required)
sudo megaraid-exporter

# Access metrics directly via HTTP
curl -s http://localhost:9272/metrics

# Check specific metrics
curl -s http://localhost:9272/metrics | grep megaraid_controller
curl -s http://localhost:9272/metrics | grep megaraid_drive_temperature
```

### Command Line Options
```
--port          Port to listen on (default: 9272)
--megacli-path  Path to megacli64 binary (default: /usr/sbin/megacli64)
--log-level     Log level: debug, info, warn, error (default: info)
--config        Path to configuration file
--timeout       Command timeout in seconds (default: 30)
```

### Examples of Direct Access
```bash
# Basic health check
curl -s http://localhost:9272/metrics | grep -E "(controller_status|array_status|drive_status)"

# Temperature monitoring
curl -s http://localhost:9272/metrics | grep temperature

# Drive errors
curl -s http://localhost:9272/metrics | grep errors_total

# BBU status
curl -s http://localhost:9272/metrics | grep bbu

# Format for monitoring scripts
curl -s http://localhost:9272/metrics | awk '/megaraid_drive_temperature/{print $1, $2}'
```

### Configuration File
Create `/etc/megaraid-exporter/config.yaml`:
```yaml
port: 9272
megacli_path: "/usr/sbin/megacli64"
log_level: "info"
command_timeout: 30s
metrics:
  collect_interval: 30s
  enabled_collectors:
    - controller
    - arrays
    - drives
    - bbu
    - events
```

## Systemd Service

Create `/etc/systemd/system/megaraid-exporter.service`:
```ini
[Unit]
Description=MegaRAID Prometheus Exporter (MegaCLI)
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/megaraid-exporter --config /etc/megaraid-exporter/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable megaraid-exporter
sudo systemctl start megaraid-exporter

# Verify it's working
curl -s http://localhost:9272/metrics | head -20
```

## Integration Options

### Option 1: Direct HTTP Access (No Prometheus Required)
```bash
# Shell script monitoring
#!/bin/bash
METRICS=$(curl -s http://localhost:9272/metrics)
CONTROLLER_STATUS=$(echo "$METRICS" | grep 'megaraid_controller_status' | awk '{print $2}')
if [ "$CONTROLLER_STATUS" != "0" ]; then
    echo "ALERT: Controller error detected!"
fi
```

### Option 2: Prometheus Integration (Optional)
Add to your `prometheus.yml`:
```yaml
scrape_configs:
  - job_name: 'megaraid'
    static_configs:
      - targets: ['localhost:9272']
    scrape_interval: 30s
    scrape_timeout: 25s
```

### Option 3: Other Monitoring Systems
The exporter works with any system that can scrape HTTP endpoints:
- Telegraf (InfluxDB)
- Datadog
- New Relic
- Custom monitoring scripts
- Grafana Agent

## Metrics

### Controller Metrics
- `megaraid_controller_info` - Controller information (model, firmware, driver)
- `megaraid_controller_status` - Controller status (0=OK, 1=Error)
- `megaraid_controller_temperature` - Controller temperature in Celsius
- `megaraid_controller_memory_size` - Controller memory size in MB
- `megaraid_controller_memory_errors` - Memory error count

### Array Metrics
- `megaraid_array_info` - Virtual drive information
- `megaraid_array_status` - Array status (0=Optimal, 1=Degraded, 2=Failed, 3=Offline)
- `megaraid_array_size_bytes` - Array size in bytes
- `megaraid_array_stripe_size` - Stripe size in KB
- `megaraid_array_read_policy` - Read policy (0=Normal, 1=ReadAhead, 2=Adaptive)
- `megaraid_array_write_policy` - Write policy (0=WriteThrough, 1=WriteBack)

### Drive Metrics
- `megaraid_drive_info` - Physical drive information
- `megaraid_drive_status` - Drive status (0=Online, 1=Failed, 2=Rebuilding, 3=Missing)
- `megaraid_drive_temperature` - Drive temperature in Celsius
- `megaraid_drive_errors_total` - Total drive errors
- `megaraid_drive_predictive_failures` - Predictive failure count
- `megaraid_drive_smart_errors` - SMART error count
- `megaraid_drive_rebuild_progress` - Rebuild progress percentage

### BBU Metrics
- `megaraid_bbu_status` - Battery status (0=OK, 1=Error, 2=Missing)
- `megaraid_bbu_charge_percent` - Battery charge percentage
- `megaraid_bbu_temperature` - Battery temperature in Celsius
- `megaraid_bbu_cycle_count` - Battery charge cycle count
- `megaraid_bbu_voltage` - Battery voltage
- `megaraid_bbu_remaining_time` - Estimated remaining backup time in minutes

### Event Metrics
- `megaraid_events_total` - Total event count by severity
- `megaraid_critical_events` - Critical events in last 24 hours
- `megaraid_warning_events` - Warning events in last 24 hours

## Monitoring Examples

### Simple Health Check Script
```bash
#!/bin/bash
# healthcheck.sh
URL="http://localhost:9272/metrics"

# Check if exporter is running
if ! curl -s "$URL" > /dev/null; then
    echo "ERROR: MegaRAID exporter is not responding"
    exit 1
fi

# Check controller status
CONTROLLER_OK=$(curl -s "$URL" | grep 'megaraid_controller_status{' | grep ' 0$' | wc -l)
if [ "$CONTROLLER_OK" -eq 0 ]; then
    echo "WARNING: Controller issues detected"
fi

# Check for failed drives
FAILED_DRIVES=$(curl -s "$URL" | grep 'megaraid_drive_status{' | grep ' 1$' | wc -l)
if [ "$FAILED_DRIVES" -gt 0 ]; then
    echo "CRITICAL: $FAILED_DRIVES failed drives detected"
fi

echo "Health check completed"
```

### Temperature Monitoring
```bash
#!/bin/bash
# temp_monitor.sh
METRICS=$(curl -s http://localhost:9272/metrics)

echo "=== Temperature Report ==="
echo "$METRICS" | grep 'megaraid.*_temperature{' | while read line; do
    TEMP=$(echo "$line" | awk '{print $2}')
    DEVICE=$(echo "$line" | grep -o 'controller="[^"]*"' | cut -d'"' -f2)
    echo "Device: $DEVICE, Temperature: ${TEMP}°C"
done
```

## Troubleshooting

### Common Issues

1. **Exporter Not Starting**
   ```bash
   # Check if port is available
   netstat -tlnp | grep 9272
   
   # Run with debug logging
   sudo megaraid-exporter --log-level debug
   ```

2. **No Metrics Returned**
   ```bash
   # Test MegaCLI directly
   sudo megacli64 -adpCount -NoLog
   
   # Check exporter logs
   journalctl -u megaraid-exporter -f
   ```

3. **Permission Denied**
   ```bash
   # Ensure running as root
   sudo megaraid-exporter
   ```

4. **megacli64 Not Found**
   ```bash
   # Check if MegaCLI is installed
   which megacli64
   
   # Or specify custom path
   megaraid-exporter --megacli-path /opt/MegaRAID/MegaCli/MegaCli64
   ```

### Testing the Exporter
```bash
# Basic connectivity test
curl -I http://localhost:9272/metrics

# Get sample metrics
curl -s http://localhost:9272/metrics | head -50

# Check for specific metrics
curl -s http://localhost:9272/metrics | grep -c "megaraid_"

# Validate Prometheus format
curl -s http://localhost:9272/metrics | promtool check metrics
```

## Development

### Building
```bash
make build
```

### Testing
```bash
make test
```

## Why MegaCLI over StorCLI?

- **More Detailed Information**: MegaCLI provides more granular drive and controller statistics
- **Better Event Monitoring**: Superior event log access and filtering
- **Predictive Analytics**: Enhanced SMART data and predictive failure detection
- **Wider Compatibility**: Works with older MegaRAID controllers
- **Detailed BBU Information**: More comprehensive battery monitoring

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Support

- Create an issue for bugs or feature requests
- Check existing issues before creating new ones
- Provide system information and logs when reporting issues
- Include MegaCLI version: `megacli64 -v`
- Test with direct curl access before reporting Prometheus issues