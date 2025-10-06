package diskutil

import (
	"fmt"
	"strconv"
	"strings"
)

// parseFiled parses a field value from a line based on the key and expected type
func parseFiled(line, key, fieldType string) (interface{}, error) {
	// Remove the key prefix and trim whitespace
	value := strings.TrimSpace(strings.TrimPrefix(line, key))
	
	switch fieldType {
	case typeInt:
		// Extract numeric value, handling cases like "123 C" or "456 MB"
		parts := strings.Fields(value)
		if len(parts) == 0 {
			return 0, fmt.Errorf("empty value for key %s", key)
		}
		
		// Try to parse the first part as integer
		intVal, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("failed to parse int value '%s' for key %s: %v", parts[0], key, err)
		}
		return intVal, nil
		
	case typeString:
		return value, nil
		
	case typeFloat:
		// Extract numeric value for float parsing
		parts := strings.Fields(value)
		if len(parts) == 0 {
			return 0.0, fmt.Errorf("empty value for key %s", key)
		}
		
		floatVal, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0.0, fmt.Errorf("failed to parse float value '%s' for key %s: %v", parts[0], key, err)
		}
		return floatVal, nil
		
	default:
		return nil, fmt.Errorf("unsupported field type: %s", fieldType)
	}
}

// ParsePhysicalDriveInfo parses MegaCLI output for physical drive information
func ParsePhysicalDriveInfo(output string) ([]*PhysicalDriveStat, error) {
	var drives []*PhysicalDriveStat
	var currentDrive *PhysicalDriveStat
	
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check if this is the start of a new drive section
		if strings.Contains(line, "Enclosure Device ID:") && currentDrive != nil {
			drives = append(drives, currentDrive)
			currentDrive = nil
		}
		
		// Initialize new drive if we encounter enclosure device ID
		if strings.HasPrefix(line, keyPdEnclosureDeviceId) && currentDrive == nil {
			currentDrive = &PhysicalDriveStat{}
		}
		
		// Parse the line if we have a current drive
		if currentDrive != nil {
			if err := currentDrive.parseLine(line); err != nil {
				return nil, fmt.Errorf("failed to parse line '%s': %v", line, err)
			}
		}
	}
	
	// Add the last drive if exists
	if currentDrive != nil {
		drives = append(drives, currentDrive)
	}
	
	return drives, nil
}

// ParseBatteryInfo parses MegaCLI output for battery backup unit information
func ParseBatteryInfo(output string) ([]*BatteryBackupStat, error) {
	var batteries []*BatteryBackupStat
	var currentBattery *BatteryBackupStat
	
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Initialize new battery section
		if strings.Contains(line, "BBU status for Adapter:") || strings.Contains(line, "Battery Type:") {
			if currentBattery != nil {
				batteries = append(batteries, currentBattery)
			}
			currentBattery = &BatteryBackupStat{}
		}
		
		if currentBattery != nil {
			if err := currentBattery.parseLine(line); err != nil {
				return nil, fmt.Errorf("failed to parse battery line '%s': %v", line, err)
			}
		}
	}
	
	if currentBattery != nil {
		batteries = append(batteries, currentBattery)
	}
	
	return batteries, nil
}

// ParseControllerInfo parses MegaCLI output for controller information
func ParseControllerInfo(output string) ([]*ControllerStat, error) {
	var controllers []*ControllerStat
	var currentController *ControllerStat
	
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Initialize new controller section
		if strings.Contains(line, "Adapter #") || strings.HasPrefix(line, keyCtrlProductName) {
			if currentController != nil {
				controllers = append(controllers, currentController)
			}
			currentController = &ControllerStat{}
		}
		
		if currentController != nil {
			if err := currentController.parseLine(line); err != nil {
				return nil, fmt.Errorf("failed to parse controller line '%s': %v", line, err)
			}
		}
	}
	
	if currentController != nil {
		controllers = append(controllers, currentController)
	}
	
	return controllers, nil
}

// ParseVirtualDriveInfo parses MegaCLI output for virtual drive information
func ParseVirtualDriveInfo(output string) ([]*VirtualDriveStat, error) {
	var virtualDrives []*VirtualDriveStat
	var currentVD *VirtualDriveStat
	
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Initialize new virtual drive section
		if strings.Contains(line, "Virtual Drive:") || strings.HasPrefix(line, keyVdTargetId) {
			if currentVD != nil {
				virtualDrives = append(virtualDrives, currentVD)
			}
			currentVD = &VirtualDriveStat{}
		}
		
		if currentVD != nil {
			if err := currentVD.parseLine(line); err != nil {
				return nil, fmt.Errorf("failed to parse virtual drive line '%s': %v", line, err)
			}
		}
	}
	
	if currentVD != nil {
		virtualDrives = append(virtualDrives, currentVD)
	}
	
	return virtualDrives, nil
}
