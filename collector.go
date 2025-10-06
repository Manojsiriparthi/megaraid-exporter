package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type MegaRAIDCollector struct {
	storCliPath string
	
	// Controller metrics
	controllerStatus *prometheus.Desc
	controllerTemp   *prometheus.Desc
	
	// Virtual Drive metrics
	vdStatus *prometheus.Desc
	vdSize   *prometheus.Desc
	
	// Physical Drive metrics
	pdStatus      *prometheus.Desc
	pdTemp        *prometheus.Desc
	pdMediaErrors *prometheus.Desc
	pdOtherErrors *prometheus.Desc
	pdPredictiveFailures *prometheus.Desc
}

type StorCliResponse struct {
	Controllers []Controller `json:"Controllers"`
}

type Controller struct {
	CommandStatus struct {
		Status string `json:"Status"`
	} `json:"Command Status"`
	ResponseData ControllerData `json:"Response Data"`
}

type ControllerData struct {
	SystemOverview []SystemInfo `json:"System Overview"`
	VDList         []VDInfo     `json:"VD LIST,omitempty"`
	PDList         []PDInfo     `json:"PD LIST,omitempty"`
}

type SystemInfo struct {
	Ctl           int    `json:"Ctl"`
	Model         string `json:"Model"`
	SerialNo      string `json:"SerialNo"`
	Status        string `json:"Status"`
	Temperature   string `json:"ROCtemp"`
}

type VDInfo struct {
	DGVD   string `json:"DG/VD"`
	Type   string `json:"TYPE"`
	State  string `json:"State"`
	Access string `json:"Access"`
	Size   string `json:"Size"`
}

type PDInfo struct {
	EIDSlt string `json:"EID:Slt"`
	DID    int    `json:"DID"`
	State  string `json:"State"`
	DGrp   string `json:"DG"`
	Size   string `json:"Size"`
	Intf   string `json:"Intf"`
	Med    string `json:"Med"`
	SED    string `json:"SED"`
	PI     string `json:"PI"`
	SeSz   string `json:"SeSz"`
	Model  string `json:"Model"`
	Sp     string `json:"Sp"`
	Type   string `json:"Type"`
	Temp   string `json:"Temp"`
	MediaErr string `json:"Med Err"`
	OtherErr string `json:"Other Err"`
	PredFail string `json:"Pred Fail"`
}

func NewMegaRAIDCollector(storCliPath string) *MegaRAIDCollector {
	return &MegaRAIDCollector{
		storCliPath: storCliPath,
		controllerStatus: prometheus.NewDesc(
			"megaraid_controller_status",
			"Status of MegaRAID controller (1=optimal, 0=not optimal)",
			[]string{"controller", "model", "serial"},
			nil,
		),
		controllerTemp: prometheus.NewDesc(
			"megaraid_controller_temperature_celsius",
			"Temperature of MegaRAID controller in Celsius",
			[]string{"controller", "model", "serial"},
			nil,
		),
		vdStatus: prometheus.NewDesc(
			"megaraid_vd_status",
			"Status of virtual drive (1=optimal, 0=not optimal)",
			[]string{"controller", "vd", "type", "access"},
			nil,
		),
		vdSize: prometheus.NewDesc(
			"megaraid_vd_size_bytes",
			"Size of virtual drive in bytes",
			[]string{"controller", "vd", "type"},
			nil,
		),
		pdStatus: prometheus.NewDesc(
			"megaraid_pd_status",
			"Status of physical drive (1=online, 0=not online)",
			[]string{"controller", "enclosure_slot", "model", "type"},
			nil,
		),
		pdTemp: prometheus.NewDesc(
			"megaraid_pd_temperature_celsius",
			"Temperature of physical drive in Celsius",
			[]string{"controller", "enclosure_slot", "model"},
			nil,
		),
		pdMediaErrors: prometheus.NewDesc(
			"megaraid_pd_media_errors_total",
			"Total media errors on physical drive",
			[]string{"controller", "enclosure_slot", "model"},
			nil,
		),
		pdOtherErrors: prometheus.NewDesc(
			"megaraid_pd_other_errors_total",
			"Total other errors on physical drive",
			[]string{"controller", "enclosure_slot", "model"},
			nil,
		),
		pdPredictiveFailures: prometheus.NewDesc(
			"megaraid_pd_predictive_failures_total",
			"Total predictive failures on physical drive",
			[]string{"controller", "enclosure_slot", "model"},
			nil,
		),
	}
}

func (c *MegaRAIDCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.controllerStatus
	ch <- c.controllerTemp
	ch <- c.vdStatus
	ch <- c.vdSize
	ch <- c.pdStatus
	ch <- c.pdTemp
	ch <- c.pdMediaErrors
	ch <- c.pdOtherErrors
	ch <- c.pdPredictiveFailures
}

func (c *MegaRAIDCollector) Collect(ch chan<- prometheus.Metric) {
	// Collect controller info
	controllers, err := c.getControllerInfo()
	if err != nil {
		log.Printf("Error collecting controller info: %v", err)
		return
	}

	for _, ctrl := range controllers {
		c.collectControllerMetrics(ch, ctrl)
	}

	// Collect VD info
	vds, err := c.getVirtualDrives()
	if err != nil {
		log.Printf("Error collecting VD info: %v", err)
	} else {
		for ctlId, vdList := range vds {
			c.collectVDMetrics(ch, ctlId, vdList)
		}
	}

	// Collect PD info
	pds, err := c.getPhysicalDrives()
	if err != nil {
		log.Printf("Error collecting PD info: %v", err)
	} else {
		for ctlId, pdList := range pds {
			c.collectPDMetrics(ch, ctlId, pdList)
		}
	}
}

func (c *MegaRAIDCollector) getControllerInfo() ([]SystemInfo, error) {
	cmd := exec.Command(c.storCliPath, "show", "all", "J")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute storcli: %v", err)
	}

	var response StorCliResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var controllers []SystemInfo
	for _, ctrl := range response.Controllers {
		controllers = append(controllers, ctrl.ResponseData.SystemOverview...)
	}

	return controllers, nil
}

func (c *MegaRAIDCollector) getVirtualDrives() (map[int][]VDInfo, error) {
	cmd := exec.Command(c.storCliPath, "/call", "show", "all", "J")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute storcli: %v", err)
	}

	var response StorCliResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	result := make(map[int][]VDInfo)
	for i, ctrl := range response.Controllers {
		result[i] = ctrl.ResponseData.VDList
	}

	return result, nil
}

func (c *MegaRAIDCollector) getPhysicalDrives() (map[int][]PDInfo, error) {
	cmd := exec.Command(c.storCliPath, "/call", "show", "all", "J")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute storcli: %v", err)
	}

	var response StorCliResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	result := make(map[int][]PDInfo)
	for i, ctrl := range response.Controllers {
		result[i] = ctrl.ResponseData.PDList
	}

	return result, nil
}

func (c *MegaRAIDCollector) collectControllerMetrics(ch chan<- prometheus.Metric, ctrl SystemInfo) {
	ctlStr := strconv.Itoa(ctrl.Ctl)
	
	// Controller status
	status := 0.0
	if strings.ToLower(ctrl.Status) == "optimal" {
		status = 1.0
	}
	ch <- prometheus.MustNewConstMetric(
		c.controllerStatus,
		prometheus.GaugeValue,
		status,
		ctlStr, ctrl.Model, ctrl.SerialNo,
	)

	// Controller temperature
	if temp := parseTemperature(ctrl.Temperature); temp > 0 {
		ch <- prometheus.MustNewConstMetric(
			c.controllerTemp,
			prometheus.GaugeValue,
			temp,
			ctlStr, ctrl.Model, ctrl.SerialNo,
		)
	}
}

func (c *MegaRAIDCollector) collectVDMetrics(ch chan<- prometheus.Metric, ctlId int, vds []VDInfo) {
	ctlStr := strconv.Itoa(ctlId)
	
	for _, vd := range vds {
		// VD status
		status := 0.0
		if strings.ToLower(vd.State) == "optl" {
			status = 1.0
		}
		ch <- prometheus.MustNewConstMetric(
			c.vdStatus,
			prometheus.GaugeValue,
			status,
			ctlStr, vd.DGVD, vd.Type, vd.Access,
		)

		// VD size
		if size := parseSize(vd.Size); size > 0 {
			ch <- prometheus.MustNewConstMetric(
				c.vdSize,
				prometheus.GaugeValue,
				size,
				ctlStr, vd.DGVD, vd.Type,
			)
		}
	}
}

func (c *MegaRAIDCollector) collectPDMetrics(ch chan<- prometheus.Metric, ctlId int, pds []PDInfo) {
	ctlStr := strconv.Itoa(ctlId)
	
	for _, pd := range pds {
		// PD status
		status := 0.0
		if strings.ToLower(pd.State) == "onln" {
			status = 1.0
		}
		ch <- prometheus.MustNewConstMetric(
			c.pdStatus,
			prometheus.GaugeValue,
			status,
			ctlStr, pd.EIDSlt, pd.Model, pd.Type,
		)

		// PD temperature
		if temp := parseTemperature(pd.Temp); temp > 0 {
			ch <- prometheus.MustNewConstMetric(
				c.pdTemp,
				prometheus.GaugeValue,
				temp,
				ctlStr, pd.EIDSlt, pd.Model,
			)
		}

		// Error counts
		if mediaErr := parseErrorCount(pd.MediaErr); mediaErr >= 0 {
			ch <- prometheus.MustNewConstMetric(
				c.pdMediaErrors,
				prometheus.CounterValue,
				mediaErr,
				ctlStr, pd.EIDSlt, pd.Model,
			)
		}

		if otherErr := parseErrorCount(pd.OtherErr); otherErr >= 0 {
			ch <- prometheus.MustNewConstMetric(
				c.pdOtherErrors,
				prometheus.CounterValue,
				otherErr,
				ctlStr, pd.EIDSlt, pd.Model,
			)
		}

		if predFail := parseErrorCount(pd.PredFail); predFail >= 0 {
			ch <- prometheus.MustNewConstMetric(
				c.pdPredictiveFailures,
				prometheus.CounterValue,
				predFail,
				ctlStr, pd.EIDSlt, pd.Model,
			)
		}
	}
}

func parseTemperature(tempStr string) float64 {
	if tempStr == "" || tempStr == "N/A" {
		return 0
	}
	tempStr = strings.TrimSuffix(tempStr, "C")
	if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
		return temp
	}
	return 0
}

func parseSize(sizeStr string) float64 {
	if sizeStr == "" || sizeStr == "N/A" {
		return 0
	}
	
	sizeStr = strings.TrimSpace(sizeStr)
	multiplier := 1.0
	
	if strings.HasSuffix(sizeStr, "TB") {
		multiplier = 1024 * 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "TB")
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "GB")
	}
	
	if size, err := strconv.ParseFloat(strings.TrimSpace(sizeStr), 64); err == nil {
		return size * multiplier
	}
	return 0
}

func parseErrorCount(errStr string) float64 {
	if errStr == "" || errStr == "N/A" || errStr == "-" {
		return 0
	}
	if count, err := strconv.ParseFloat(errStr, 64); err == nil {
		return count
	}
	return -1
}
