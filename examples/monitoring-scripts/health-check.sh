#!/bin/bash
# MegaRAID Health Check Script
# This script demonstrates how to use the exporter without Prometheus

set -e

EXPORTER_URL="http://localhost:9272/metrics"
TEMP_FILE="/tmp/megaraid_metrics"

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Check if exporter is running
echo "=== MegaRAID Health Check ==="
if ! curl -s "$EXPORTER_URL" > "$TEMP_FILE"; then
    print_status $RED "ERROR: MegaRAID exporter is not responding on $EXPORTER_URL"
    exit 1
fi

print_status $GREEN "✓ Exporter is responding"

# Check controller status
CONTROLLER_ERRORS=$(grep 'megaraid_controller_status{' "$TEMP_FILE" | grep -v ' 0$' | wc -l)
if [ "$CONTROLLER_ERRORS" -gt 0 ]; then
    print_status $RED "✗ CRITICAL: Controller errors detected ($CONTROLLER_ERRORS)"
    grep 'megaraid_controller_status{' "$TEMP_FILE" | grep -v ' 0$'
else
    print_status $GREEN "✓ All controllers are healthy"
fi

# Check array status
ARRAY_ERRORS=$(grep 'megaraid_array_status{' "$TEMP_FILE" | grep -v ' 0$' | wc -l)
if [ "$ARRAY_ERRORS" -gt 0 ]; then
    print_status $RED "✗ CRITICAL: Array issues detected ($ARRAY_ERRORS)"
    grep 'megaraid_array_status{' "$TEMP_FILE" | grep -v ' 0$'
else
    print_status $GREEN "✓ All arrays are optimal"
fi

# Check drive status
FAILED_DRIVES=$(grep 'megaraid_drive_status{' "$TEMP_FILE" | grep ' 1$' | wc -l)
REBUILDING_DRIVES=$(grep 'megaraid_drive_status{' "$TEMP_FILE" | grep ' 2$' | wc -l)

if [ "$FAILED_DRIVES" -gt 0 ]; then
    print_status $RED "✗ CRITICAL: $FAILED_DRIVES failed drives detected"
    grep 'megaraid_drive_status{' "$TEMP_FILE" | grep ' 1$'
elif [ "$REBUILDING_DRIVES" -gt 0 ]; then
    print_status $YELLOW "⚠ WARNING: $REBUILDING_DRIVES drives rebuilding"
    grep 'megaraid_drive_status{' "$TEMP_FILE" | grep ' 2$'
else
    print_status $GREEN "✓ All drives are online"
fi

# Check temperatures
HIGH_TEMP_CONTROLLERS=$(grep 'megaraid_controller_temperature{' "$TEMP_FILE" | awk '$2 > 70' | wc -l)
HIGH_TEMP_DRIVES=$(grep 'megaraid_drive_temperature{' "$TEMP_FILE" | awk '$2 > 60' | wc -l)

if [ "$HIGH_TEMP_CONTROLLERS" -gt 0 ] || [ "$HIGH_TEMP_DRIVES" -gt 0 ]; then
    print_status $YELLOW "⚠ WARNING: High temperature detected"
    [ "$HIGH_TEMP_CONTROLLERS" -gt 0 ] && grep 'megaraid_controller_temperature{' "$TEMP_FILE" | awk '$2 > 70'
    [ "$HIGH_TEMP_DRIVES" -gt 0 ] && grep 'megaraid_drive_temperature{' "$TEMP_FILE" | awk '$2 > 60'
else
    print_status $GREEN "✓ All temperatures are normal"
fi

# Check BBU status
BBU_ERRORS=$(grep 'megaraid_bbu_status{' "$TEMP_FILE" | grep -v ' 0$' | wc -l)
if [ "$BBU_ERRORS" -gt 0 ]; then
    print_status $YELLOW "⚠ WARNING: BBU issues detected ($BBU_ERRORS)"
    grep 'megaraid_bbu_status{' "$TEMP_FILE" | grep -v ' 0$'
else
    print_status $GREEN "✓ BBU status is healthy"
fi

# Summary
echo
echo "=== Summary ==="
TOTAL_CONTROLLERS=$(grep 'megaraid_controller_info{' "$TEMP_FILE" | wc -l)
TOTAL_ARRAYS=$(grep 'megaraid_array_info{' "$TEMP_FILE" | wc -l)
TOTAL_DRIVES=$(grep 'megaraid_drive_info{' "$TEMP_FILE" | wc -l)

echo "Controllers: $TOTAL_CONTROLLERS"
echo "Arrays: $TOTAL_ARRAYS" 
echo "Drives: $TOTAL_DRIVES"

# Cleanup
rm -f "$TEMP_FILE"

# Exit with appropriate code
if [ "$CONTROLLER_ERRORS" -gt 0 ] || [ "$ARRAY_ERRORS" -gt 0 ] || [ "$FAILED_DRIVES" -gt 0 ]; then
    exit 2  # Critical issues
elif [ "$REBUILDING_DRIVES" -gt 0 ] || [ "$HIGH_TEMP_CONTROLLERS" -gt 0 ] || [ "$HIGH_TEMP_DRIVES" -gt 0 ] || [ "$BBU_ERRORS" -gt 0 ]; then
    exit 1  # Warning issues
else
    exit 0  # All good
fi
