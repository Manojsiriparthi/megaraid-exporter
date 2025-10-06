package diskutil

import (
	"strings"
)

// BatteryBackupStat represents the statistics of a battery backup unit
type BatteryBackupStat struct {
	AdapterIndex              int    `json:"adapter_index"`
	BatteryType               string `json:"battery_type"`
	BatteryState              string `json:"battery_state"`
	ChargeStatus              string `json:"charge_status"`
	ChargeLevel               int    `json:"charge_level"`
	Temperature               int    `json:"temperature"`
	DesignCapacity            int    `json:"design_capacity"`
	FullChargeCapacity        int    `json:"full_charge_capacity"`
	CycleCount                int    `json:"cycle_count"`
	ReplacementRequired       string `json:"replacement_required"`
	PackMissing               string `json:"pack_missing"`
	VoltageStatus             string `json:"voltage_status"`
	CurrentStatus             string `json:"current_status"`
	CapacitanceStatus         string `json:"capacitance_status"`
	LearnCycleActive          string `json:"learn_cycle_active"`
	NextLearnTime             string `json:"next_learn_time"`
	ManufactureDate           string `json:"manufacture_date"`
	SerialNumber              string `json:"serial_number"`
	FirmwareVersion           string `json:"firmware_version"`
}

func (b *BatteryBackupStat) parseLine(line string) error {
	if strings.HasPrefix(line, keyBbuBatteryType) {
		batteryType, err := parseFiled(line, keyBbuBatteryType, typeString)
		if err != nil {
			return err
		}
		b.BatteryType = batteryType.(string)
	} else if strings.HasPrefix(line, keyBbuBatteryState) {
		batteryState, err := parseFiled(line, keyBbuBatteryState, typeString)
		if err != nil {
			return err
		}
		b.BatteryState = batteryState.(string)
	} else if strings.HasPrefix(line, keyBbuChargeStatus) {
		chargeStatus, err := parseFiled(line, keyBbuChargeStatus, typeString)
		if err != nil {
			return err
		}
		b.ChargeStatus = chargeStatus.(string)
	} else if strings.HasPrefix(line, keyBbuAbsoluteStateOfCharge) {
		chargeLevel, err := parseFiled(line, keyBbuAbsoluteStateOfCharge, typeInt)
		if err != nil {
			return err
		}
		b.ChargeLevel = chargeLevel.(int)
	} else if strings.HasPrefix(line, keyBbuTemperature) {
		temperature, err := parseFiled(line, keyBbuTemperature, typeInt)
		if err != nil {
			return err
		}
		b.Temperature = temperature.(int)
	} else if strings.HasPrefix(line, keyBbuDesignCapacity) {
		designCapacity, err := parseFiled(line, keyBbuDesignCapacity, typeInt)
		if err != nil {
			return err
		}
		b.DesignCapacity = designCapacity.(int)
	} else if strings.HasPrefix(line, keyBbuFullChargeCapacity) {
		fullChargeCapacity, err := parseFiled(line, keyBbuFullChargeCapacity, typeInt)
		if err != nil {
			return err
		}
		b.FullChargeCapacity = fullChargeCapacity.(int)
	} else if strings.HasPrefix(line, keyBbuCycleCount) {
		cycleCount, err := parseFiled(line, keyBbuCycleCount, typeInt)
		if err != nil {
			return err
		}
		b.CycleCount = cycleCount.(int)
	} else if strings.HasPrefix(line, keyBbuReplacementRequired) {
		replacementRequired, err := parseFiled(line, keyBbuReplacementRequired, typeString)
		if err != nil {
			return err
		}
		b.ReplacementRequired = replacementRequired.(string)
	} else if strings.HasPrefix(line, keyBbuPackMissing) {
		packMissing, err := parseFiled(line, keyBbuPackMissing, typeString)
		if err != nil {
			return err
		}
		b.PackMissing = packMissing.(string)
	} else if strings.HasPrefix(line, keyBbuVoltageStatus) {
		voltageStatus, err := parseFiled(line, keyBbuVoltageStatus, typeString)
		if err != nil {
			return err
		}
		b.VoltageStatus = voltageStatus.(string)
	} else if strings.HasPrefix(line, keyBbuLearnCycleActive) {
		learnCycleActive, err := parseFiled(line, keyBbuLearnCycleActive, typeString)
		if err != nil {
			return err
		}
		b.LearnCycleActive = learnCycleActive.(string)
	} else if strings.HasPrefix(line, keyBbuNextLearnTime) {
		nextLearnTime, err := parseFiled(line, keyBbuNextLearnTime, typeString)
		if err != nil {
			return err
		}
		b.NextLearnTime = nextLearnTime.(string)
	}
	return nil
}
