package diskutil

// Physical Drive parsing keys
const (
	keyPdEnclosureDeviceId        = "Enclosure Device ID:"
	keyPdDeviceId                 = "Device Id:"
	keyPdSlotNumber               = "Slot Number:"
	keyPdMediaErrorCount          = "Media Error Count:"
	keyPdOtherErrorCount          = "Other Error Count:"
	keyPdPredictiveFailureCount   = "Predictive Failure Count:"
	keyPdPdtype                   = "PD Type:"
	keyPdRawSize                  = "Raw Size:"
	keyPdFirmwareState            = "Firmware state:"
	keyPdInquiryData              = "Inquiry Data:"
	keyPdDriveTemperature         = "Drive Temperature:"
	keyPdSMARTFlag                = "SMART Flag:"
	keyPdSMARTAlertFlagged        = "SMART alert flagged by drive:"
	keyPdLastPredictiveFailureEventSeqNum = "Last Predictive Failure Event Seq Number:"
)

// Battery Backup Unit parsing keys
const (
	keyBbuBatteryType             = "Battery Type:"
	keyBbuBatteryState            = "Battery State:"
	keyBbuChargeStatus            = "Charge Status:"
	keyBbuAbsoluteStateOfCharge   = "Absolute State of charge:"
	keyBbuTemperature             = "Temperature:"
	keyBbuDesignCapacity          = "Design Capacity:"
	keyBbuFullChargeCapacity      = "Full Charge Capacity:"
	keyBbuCycleCount              = "Cycle Count:"
	keyBbuReplacementRequired     = "Replacement required:"
	keyBbuPackMissing             = "Pack is about to fail & should be replaced:"
	keyBbuVoltageStatus           = "Battery Voltage:"
	keyBbuCurrentStatus           = "Battery Current:"
	keyBbuCapacitanceStatus       = "Capacitance Status:"
	keyBbuLearnCycleActive        = "Learn Cycle Active:"
	keyBbuNextLearnTime           = "Next Learn time:"
	keyBbuManufactureDate         = "Manufacture Date:"
	keyBbuSerialNumber            = "Serial Number:"
	keyBbuFirmwareVersion         = "Firmware Version:"
)

// Controller parsing keys
const (
	keyCtrlProductName            = "Product Name:"
	keyCtrlSerialNumber           = "Serial No:"
	keyCtrlFWVersion              = "FW Version:"
	keyCtrlBIOSVersion            = "BIOS Version:"
	keyCtrlTemperature            = "Controller temperature:"
	keyCtrlROCTemperature         = "ROC temperature:"
	keyCtrlMemorySize             = "Memory Size:"
	keyCtrlMemoryType             = "Memory Type:"
	keyCtrlAlarmState             = "Alarm State:"
	keyCtrlRebuildRate            = "Rebuild Rate:"
	keyCtrlPatrolReadRate         = "Patrol Read Rate:"
	keyCtrlBGIRate                = "BGI Rate:"
	keyCtrlCCRate                 = "CC Rate:"
	keyCtrlReconstructionRate     = "Reconstruction Rate:"
	keyCtrlClusterActive          = "Cluster Active:"
	keyCtrlClusterSupported       = "Cluster Supported:"
)

// Virtual Drive parsing keys
const (
	keyVdTargetId                 = "Target Id:"
	keyVdName                     = "Name:"
	keyVdRAIDLevel                = "RAID Level:"
	keyVdSize                     = "Size:"
	keyVdState                    = "State:"
	keyVdStripSize                = "Strip Size:"
	keyVdNumberOfDrives           = "Number Of Drives:"
	keyVdSpanDepth                = "Span Depth:"
	keyVdDefaultCachePolicy       = "Default Cache Policy:"
	keyVdCurrentCachePolicy       = "Current Cache Policy:"
	keyVdAccessPolicy             = "Access Policy:"
	keyVdDiskCachePolicy          = "Disk Cache Policy:"
	keyVdOngoingProgresses        = "Ongoing Progresses:"
	keyVdBadBlocksExist           = "Bad Blocks Exist:"
	keyVdIsVDCached               = "Is VD Cached:"
)

// Parse field types
const (
	typeInt    = "int"
	typeString = "string"
	typeFloat  = "float"
)
