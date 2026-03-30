package entity

import "time"

type HypervisorNode struct {
	NodeID     string
	Name       string
	ZoneID     string
	Zone       string
	Status     string
	CPUModel   string
	CPUCores   int
	RAMTotalMB uint64
	DiskCount  int
	GPUCount   int
}

type HypervisorDiskInventoryItem struct {
	Name   string
	Model  string
	SizeGB uint64
}

type HypervisorGPUInventoryItem struct {
	Model         string
	Vendor        string
	PCIAddress    string
	DriverVersion string
	MemoryTotalMB uint64
}

type HypervisorMetricPoint struct {
	Timestamp time.Time
	Value     float64
}

type HypervisorMetricSeries struct {
	Name   string
	Label  string
	Unit   string
	Latest float64
	Points []HypervisorMetricPoint
}

type HypervisorNodeDetail struct {
	NodeID     string
	Name       string
	Hostname   string
	ZoneID     string
	Zone       string
	Status     string
	CPUModel   string
	CPUCores   int
	RAMTotalMB uint64
	DiskCount  int
	GPUCount   int
	Disks      []HypervisorDiskInventoryItem
	GPUs       []HypervisorGPUInventoryItem
}
