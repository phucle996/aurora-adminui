package repository

import (
	"context"
	"encoding/json"
	"strings"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	"aurora-adminui/internal/errorx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HypervisorRepoImple struct {
	db *pgxpool.Pool
}

func NewHypervisorRepo(db *pgxpool.Pool) domainrepo.HypervisorRepository {
	return &HypervisorRepoImple{db: db}
}

func (r *HypervisorRepoImple) ListNodes(ctx context.Context) ([]entity.HypervisorNode, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT n.node_id,
		        COALESCE(NULLIF(n.display_name, ''), NULLIF(n.hostname, ''), n.node_id) AS display_name,
		        COALESCE(z.id::text, '') AS zone_id,
		        COALESCE(z.name, NULLIF(n.zone, ''), '') AS zone_name,
		        n.status,
		        COALESCE(inv.cpu_model, ''),
		        COALESCE(inv.cpu_cores, 0),
		        COALESCE(inv.ram_total_bytes, 0) / 1048576,
		        COALESCE(jsonb_array_length(inv.disks), 0),
		        COALESCE(jsonb_array_length(inv.gpus), 0)
		 FROM hypervisor.nodes n
		 LEFT JOIN zone.zone_objects zo
		   ON zo.object_type = 'hypervisor_node'
		  AND zo.object_id = n.node_id
		 LEFT JOIN zone.zones z ON z.id = zo.zone_id
		 LEFT JOIN hypervisor.node_hardware_inventory inv ON inv.node_id = n.node_id
		 ORDER BY n.updated_at DESC, n.node_id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.HypervisorNode, 0)
	for rows.Next() {
		var item entity.HypervisorNode
		if err := rows.Scan(
			&item.NodeID,
			&item.Name,
			&item.ZoneID,
			&item.Zone,
			&item.Status,
			&item.CPUModel,
			&item.CPUCores,
			&item.RAMTotalMB,
			&item.DiskCount,
			&item.GPUCount,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *HypervisorRepoImple) GetNodeDetail(ctx context.Context, nodeID string) (*entity.HypervisorNodeDetail, error) {
	var item entity.HypervisorNodeDetail
	var disksRaw []byte
	var gpusRaw []byte
	err := r.db.QueryRow(
		ctx,
		`SELECT n.node_id,
		        COALESCE(NULLIF(n.display_name, ''), NULLIF(n.hostname, ''), n.node_id) AS display_name,
		        COALESCE(NULLIF(n.hostname, ''), n.node_id) AS hostname,
		        COALESCE(z.id::text, '') AS zone_id,
		        COALESCE(z.name, NULLIF(n.zone, ''), '') AS zone_name,
		        n.status,
		        COALESCE(inv.cpu_model, ''),
		        COALESCE(inv.cpu_cores, 0),
		        COALESCE(inv.ram_total_bytes, 0) / 1048576,
		        COALESCE(jsonb_array_length(inv.disks), 0),
		        COALESCE(jsonb_array_length(inv.gpus), 0),
		        COALESCE(inv.disks, '[]'::jsonb),
		        COALESCE(inv.gpus, '[]'::jsonb)
		   FROM hypervisor.nodes n
		   LEFT JOIN zone.zone_objects zo
		     ON zo.object_type = 'hypervisor_node'
		    AND zo.object_id = n.node_id
		   LEFT JOIN zone.zones z ON z.id = zo.zone_id
		   LEFT JOIN hypervisor.node_hardware_inventory inv ON inv.node_id = n.node_id
		  WHERE n.node_id = $1
		  LIMIT 1`,
		strings.TrimSpace(nodeID),
	).Scan(
		&item.NodeID,
		&item.Name,
		&item.Hostname,
		&item.ZoneID,
		&item.Zone,
		&item.Status,
		&item.CPUModel,
		&item.CPUCores,
		&item.RAMTotalMB,
		&item.DiskCount,
		&item.GPUCount,
		&disksRaw,
		&gpusRaw,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errorx.ErrHypervisorNodeNotFound
		}
		return nil, err
	}
	item.Disks = decodeDiskInventory(disksRaw)
	item.GPUs = decodeGPUInventory(gpusRaw)
	return &item, nil
}

func (r *HypervisorRepoImple) UpdateNodeName(ctx context.Context, nodeID, name string) error {
	tag, err := r.db.Exec(
		ctx,
		`UPDATE hypervisor.nodes
		    SET display_name = $2,
		        updated_at = NOW()
		  WHERE node_id = $1`,
		strings.TrimSpace(nodeID),
		strings.TrimSpace(name),
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errorx.ErrHypervisorNodeNotFound
	}
	return nil
}

func (r *HypervisorRepoImple) AssignNodeToZone(ctx context.Context, nodeID string, zoneID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var zoneName string
	if err := tx.QueryRow(ctx, `SELECT name FROM zone.zones WHERE id = $1`, zoneID).Scan(&zoneName); err != nil {
		if err == pgx.ErrNoRows {
			return errorx.ErrZoneNotFound
		}
		return err
	}

	tag, err := tx.Exec(
		ctx,
		`UPDATE hypervisor.nodes
		    SET zone = $2,
		        updated_at = NOW()
		  WHERE node_id = $1`,
		strings.TrimSpace(nodeID),
		zoneName,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errorx.ErrHypervisorNodeNotFound
	}

	if _, err := tx.Exec(
		ctx,
		`INSERT INTO zone.zone_objects (object_type, object_id, zone_id, created_at)
		 VALUES ('hypervisor_node', $1, $2, NOW())
		 ON CONFLICT (object_type, object_id) DO UPDATE
		 SET zone_id = EXCLUDED.zone_id`,
		strings.TrimSpace(nodeID),
		zoneID,
	); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

type diskInventoryRow struct {
	Name      string `json:"name"`
	Model     string `json:"model"`
	SizeGB    uint64 `json:"size_gb"`
	SizeBytes uint64 `json:"size_bytes"`
}

type gpuInventoryRow struct {
	Model            string `json:"model"`
	Vendor           string `json:"vendor"`
	PCIAddress       string `json:"pci_address"`
	DriverVersion    string `json:"driver_version"`
	MemoryTotalMB    uint64 `json:"memory_total_mb"`
	MemoryTotalBytes uint64 `json:"memory_total_bytes"`
}

func decodeDiskInventory(raw []byte) []entity.HypervisorDiskInventoryItem {
	if len(raw) == 0 {
		return nil
	}
	var rows []diskInventoryRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil
	}
	out := make([]entity.HypervisorDiskInventoryItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, entity.HypervisorDiskInventoryItem{
			Name:   strings.TrimSpace(row.Name),
			Model:  strings.TrimSpace(row.Model),
			SizeGB: coalesceDiskSizeGB(row.SizeGB, row.SizeBytes),
		})
	}
	return out
}

func coalesceDiskSizeGB(sizeGB, sizeBytes uint64) uint64 {
	if sizeGB > 0 {
		return sizeGB
	}
	if sizeBytes == 0 {
		return 0
	}
	return sizeBytes / 1024 / 1024 / 1024
}

func decodeGPUInventory(raw []byte) []entity.HypervisorGPUInventoryItem {
	if len(raw) == 0 {
		return nil
	}
	var rows []gpuInventoryRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil
	}
	out := make([]entity.HypervisorGPUInventoryItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, entity.HypervisorGPUInventoryItem{
			Model:         strings.TrimSpace(row.Model),
			Vendor:        strings.TrimSpace(row.Vendor),
			PCIAddress:    strings.TrimSpace(row.PCIAddress),
			DriverVersion: strings.TrimSpace(row.DriverVersion),
			MemoryTotalMB: coalesceMemoryTotalMB(row.MemoryTotalMB, row.MemoryTotalBytes),
		})
	}
	return out
}

func coalesceMemoryTotalMB(memoryTotalMB, memoryTotalBytes uint64) uint64 {
	if memoryTotalMB > 0 {
		return memoryTotalMB
	}
	if memoryTotalBytes == 0 {
		return 0
	}
	return memoryTotalBytes / 1024 / 1024
}
