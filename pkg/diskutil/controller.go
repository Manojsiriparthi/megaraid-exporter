package diskutil

import (
	"strings"
)

// ControllerStat represents the statistics of a RAID controller
type ControllerStat struct {
	AdapterIndex           int    `json:"adapter_index"`
	ProductName            string `json:"product_name"`
	SerialNumber           string `json:"serial_number"`
	FWVersion              string `json:"fw_version"`
	BIOSVersion            string `json:"bios_version"`
	ControllerTemperature  int    `json:"controller_temperature"`
	ROCTemperature         int    `json:"roc_temperature"`
	MemorySize             int    `json:"memory_size"`
	MemoryType             string `json:"memory_type"`
	ControllerStatus       string `json:"controller_status"`
	AlarmState             string `json:"alarm_state"`
	RebuildRate            int    `json:"rebuild_rate"`
	PatrolReadRate         int    `json:"patrol_read_rate"`
	BGIRate                int    `json:"bgi_rate"`
	CCRate                 int    `json:"cc_rate"`
	ReconstructionRate     int    `json:"reconstruction_rate"`
	ClusterActive          string `json:"cluster_active"`
	ClusterSupported       string `json:"cluster_supported"`
	MaxDrivesPerSpan       int    `json:"max_drives_per_span"`
	MaxSpansPerArray       int    `json:"max_spans_per_array"`
}

func (c *ControllerStat) parseLine(line string) error {
	if strings.HasPrefix(line, keyCtrlProductName) {
		productName, err := parseFiled(line, keyCtrlProductName, typeString)
		if err != nil {
			return err
		}
		c.ProductName = productName.(string)
	} else if strings.HasPrefix(line, keyCtrlSerialNumber) {
		serialNumber, err := parseFiled(line, keyCtrlSerialNumber, typeString)
		if err != nil {
			return err
		}
		c.SerialNumber = serialNumber.(string)
	} else if strings.HasPrefix(line, keyCtrlFWVersion) {
		fwVersion, err := parseFiled(line, keyCtrlFWVersion, typeString)
		if err != nil {
			return err
		}
		c.FWVersion = fwVersion.(string)
	} else if strings.HasPrefix(line, keyCtrlTemperature) {
		temperature, err := parseFiled(line, keyCtrlTemperature, typeInt)
		if err != nil {
			return err
		}
		c.ControllerTemperature = temperature.(int)
	} else if strings.HasPrefix(line, keyCtrlROCTemperature) {
		rocTemperature, err := parseFiled(line, keyCtrlROCTemperature, typeInt)
		if err != nil {
			return err
		}
		c.ROCTemperature = rocTemperature.(int)
	} else if strings.HasPrefix(line, keyCtrlMemorySize) {
		memorySize, err := parseFiled(line, keyCtrlMemorySize, typeInt)
		if err != nil {
			return err
		}
		c.MemorySize = memorySize.(int)
	} else if strings.HasPrefix(line, keyCtrlAlarmState) {
		alarmState, err := parseFiled(line, keyCtrlAlarmState, typeString)
		if err != nil {
			return err
		}
		c.AlarmState = alarmState.(string)
	} else if strings.HasPrefix(line, keyCtrlRebuildRate) {
		rebuildRate, err := parseFiled(line, keyCtrlRebuildRate, typeInt)
		if err != nil {
			return err
		}
		c.RebuildRate = rebuildRate.(int)
	}
	return nil
}
