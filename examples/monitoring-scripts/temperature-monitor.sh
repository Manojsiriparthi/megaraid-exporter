#!/bin/bash
# Temperature monitoring script for MegaRAID
# Demonstrates direct metrics access without Prometheus

EXPORTER_URL="http://localhost:9272/metrics"
CRITICAL_TEMP=75
WARNING_TEMP=65

echo "=== MegaRAID Temperature Monitor ==="
echo "Critical threshold: ${CRITICAL_TEMP}°C"
echo "Warning threshold: ${WARNING_TEMP}°C"
echo

# Get metrics
METRICS=$(curl -s "$EXPORTER_URL")
if [ $? -ne 0 ]; then
    echo "ERROR: Cannot connect to exporter at $EXPORTER_URL"
    exit 1
fi

# Controller temperatures
echo "Controller Temperatures:"
echo "$METRICS" | grep 'megaraid_controller_temperature{' | while IFS= read -r line; do
    CONTROLLER=$(echo "$line" | grep -o 'controller="[^"]*"' | cut -d'"' -f2)
    TEMP=$(echo "$line" | awk '{print $2}')
    
    if (( $(echo "$TEMP > $CRITICAL_TEMP" | bc -l) )); then
        echo "  Controller $CONTROLLER: ${TEMP}°C [CRITICAL]"
    elif (( $(echo "$TEMP > $WARNING_TEMP" | bc -l) )); then
        echo "  Controller $CONTROLLER: ${TEMP}°C [WARNING]"
    else
        echo "  Controller $CONTROLLER: ${TEMP}°C [OK]"
    fi
done

echo
echo "Drive Temperatures:"
echo "$METRICS" | grep 'megaraid_drive_temperature{' | while IFS= read -r line; do
    CONTROLLER=$(echo "$line" | grep -o 'controller="[^"]*"' | cut -d'"' -f2)
    ENCLOSURE=$(echo "$line" | grep -o 'enclosure="[^"]*"' | cut -d'"' -f2)
    SLOT=$(echo "$line" | grep -o 'slot="[^"]*"' | cut -d'"' -f2)
    TEMP=$(echo "$line" | awk '{print $2}')
    
    if (( $(echo "$TEMP > $CRITICAL_TEMP" | bc -l) )); then
        echo "  Drive $CONTROLLER:$ENCLOSURE:$SLOT: ${TEMP}°C [CRITICAL]"
    elif (( $(echo "$TEMP > $WARNING_TEMP" | bc -l) )); then
        echo "  Drive $CONTROLLER:$ENCLOSURE:$SLOT: ${TEMP}°C [WARNING]"
    else
        echo "  Drive $CONTROLLER:$ENCLOSURE:$SLOT: ${TEMP}°C [OK]"
    fi
done
