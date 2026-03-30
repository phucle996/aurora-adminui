package response

import (
	"aurora-adminui/internal/domain/entity"
	"time"
)

type HypervisorNode struct {
	NodeID     string `json:"node_id"`
	Name       string `json:"name"`
	ZoneID     string `json:"zone_id"`
	Zone       string `json:"zone"`
	Status     string `json:"status"`
	CPUModel   string `json:"cpu_model"`
	CPUCores   int    `json:"cpu_cores"`
	RAMTotalMB uint64 `json:"ram_total_mb"`
	DiskCount  int    `json:"disk_count"`
	GPUCount   int    `json:"gpu_count"`
}

type ListHypervisorNodesResponse struct {
	Items []HypervisorNode `json:"items"`
}

type HypervisorNodeDetail struct {
	NodeID     string                        `json:"node_id"`
	Name       string                        `json:"name"`
	Hostname   string                        `json:"hostname"`
	ZoneID     string                        `json:"zone_id"`
	Zone       string                        `json:"zone"`
	Status     string                        `json:"status"`
	CPUModel   string                        `json:"cpu_model"`
	CPUCores   int                           `json:"cpu_cores"`
	RAMTotalMB uint64                        `json:"ram_total_mb"`
	DiskCount  int                           `json:"disk_count"`
	GPUCount   int                           `json:"gpu_count"`
	Disks      []HypervisorDiskInventoryItem `json:"disks"`
	GPUs       []HypervisorGPUInventoryItem  `json:"gpus"`
}

type HypervisorDiskInventoryItem struct {
	Name   string `json:"name"`
	Model  string `json:"model"`
	SizeGB uint64 `json:"size_gb"`
}

type HypervisorGPUInventoryItem struct {
	Model         string `json:"model"`
	Vendor        string `json:"vendor"`
	PCIAddress    string `json:"pci_address"`
	DriverVersion string `json:"driver_version"`
	MemoryTotalMB uint64 `json:"memory_total_mb"`
}

type HypervisorMetricPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type HypervisorMetricSeries struct {
	Name   string                  `json:"name"`
	Label  string                  `json:"label"`
	Unit   string                  `json:"unit"`
	Latest float64                 `json:"latest"`
	Points []HypervisorMetricPoint `json:"points"`
}

func NewListHypervisorNodesResponse(items []entity.HypervisorNode) ListHypervisorNodesResponse {
	out := make([]HypervisorNode, 0, len(items))
	for _, item := range items {
		out = append(out, HypervisorNode{
			NodeID:     item.NodeID,
			Name:       item.Name,
			ZoneID:     item.ZoneID,
			Zone:       item.Zone,
			Status:     item.Status,
			CPUModel:   item.CPUModel,
			CPUCores:   item.CPUCores,
			RAMTotalMB: item.RAMTotalMB,
			DiskCount:  item.DiskCount,
			GPUCount:   item.GPUCount,
		})
	}
	return ListHypervisorNodesResponse{Items: out}
}

func NewHypervisorNodeDetailResponse(item *entity.HypervisorNodeDetail) *HypervisorNodeDetail {
	if item == nil {
		return nil
	}
	disks := make([]HypervisorDiskInventoryItem, 0, len(item.Disks))
	for _, disk := range item.Disks {
		disks = append(disks, HypervisorDiskInventoryItem{
			Name:   disk.Name,
			Model:  disk.Model,
			SizeGB: disk.SizeGB,
		})
	}
	gpus := make([]HypervisorGPUInventoryItem, 0, len(item.GPUs))
	for _, gpu := range item.GPUs {
		gpus = append(gpus, HypervisorGPUInventoryItem{
			Model:         gpu.Model,
			Vendor:        gpu.Vendor,
			PCIAddress:    gpu.PCIAddress,
			DriverVersion: gpu.DriverVersion,
			MemoryTotalMB: gpu.MemoryTotalMB,
		})
	}
	return &HypervisorNodeDetail{
		NodeID:     item.NodeID,
		Name:       item.Name,
		Hostname:   item.Hostname,
		ZoneID:     item.ZoneID,
		Zone:       item.Zone,
		Status:     item.Status,
		CPUModel:   item.CPUModel,
		CPUCores:   item.CPUCores,
		RAMTotalMB: item.RAMTotalMB,
		DiskCount:  item.DiskCount,
		GPUCount:   item.GPUCount,
		Disks:      disks,
		GPUs:       gpus,
	}
}
