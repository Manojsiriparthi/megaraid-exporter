package diskutil

import (
	"strings"
)

// VirtualDriveStat represents the statistics of a virtual drive (RAID array)
type VirtualDriveStat struct {
	TargetId            int    `json:"target_id"`
	Name                string `json:"name"`
	RAID_Level          string `json:"raid_level"`
	Size                string `json:"size"`
	State               string `json:"state"`
	StripSize           string `json:"strip_size"`
	NumberOfDrives      int    `json:"number_of_drives"`
	SpanDepth           int    `json:"span_depth"`
	DefaultCachePolicy  string `json:"default_cache_policy"`
	CurrentCachePolicy  string `json:"current_cache_policy"`
	AccessPolicy        string `json:"access_policy"`
	DiskCachePolicy     string `json:"disk_cache_policy"`
	OngoingProgresses   string `json:"ongoing_progresses"`
	BadBlocksExist      string `json:"bad_blocks_exist"`
	IsVDCached          string `json:"is_vd_cached"`
}

func (v *VirtualDriveStat) parseLine(line string) error {
	if strings.HasPrefix(line, keyVdTargetId) {
		targetId, err := parseFiled(line, keyVdTargetId, typeInt)
		if err != nil {
			return err
		}
		v.TargetId = targetId.(int)
	} else if strings.HasPrefix(line, keyVdName) {
		name, err := parseFiled(line, keyVdName, typeString)
		if err != nil {
			return err
		}
		v.Name = name.(string)
	} else if strings.HasPrefix(line, keyVdRAIDLevel) {
		raidLevel, err := parseFiled(line, keyVdRAIDLevel, typeString)
		if err != nil {
			return err
		}
		v.RAID_Level = raidLevel.(string)
	} else if strings.HasPrefix(line, keyVdSize) {
		size, err := parseFiled(line, keyVdSize, typeString)
		if err != nil {
			return err
		}
		v.Size = size.(string)
	} else if strings.HasPrefix(line, keyVdState) {
		state, err := parseFiled(line, keyVdState, typeString)
		if err != nil {
			return err
		}
		v.State = state.(string)
	} else if strings.HasPrefix(line, keyVdStripSize) {
		stripSize, err := parseFiled(line, keyVdStripSize, typeString)
		if err != nil {
			return err
		}
		v.StripSize = stripSize.(string)
	} else if strings.HasPrefix(line, keyVdNumberOfDrives) {
		numberOfDrives, err := parseFiled(line, keyVdNumberOfDrives, typeInt)
		if err != nil {
			return err
		}
		v.NumberOfDrives = numberOfDrives.(int)
	} else if strings.HasPrefix(line, keyVdBadBlocksExist) {
		badBlocksExist, err := parseFiled(line, keyVdBadBlocksExist, typeString)
		if err != nil {
			return err
		}
		v.BadBlocksExist = badBlocksExist.(string)
	}
	return nil
}
