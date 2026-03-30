export type HypervisorNode = {
  node_id: string
  name: string
  zone_id: string
  zone: string
  status: string
  cpu_model: string
  cpu_cores: number
  ram_total_mb: number
  disk_count: number
  gpu_count: number
}

export type HypervisorDiskInventoryItem = {
  name: string
  model: string
  size_gb: number
}

export type HypervisorGPUInventoryItem = {
  model: string
  vendor: string
  pci_address: string
  driver_version: string
  memory_total_mb: number
}

export type HypervisorMetricPoint = {
  timestamp: string
  value: number
}

export type HypervisorMetricSeries = {
  name: string
  label: string
  unit: string
  latest: number
  points: HypervisorMetricPoint[]
}

export type HypervisorNodeDetail = HypervisorNode & {
  hostname: string
  disks: HypervisorDiskInventoryItem[]
  gpus: HypervisorGPUInventoryItem[]
}

type HypervisorNodesResponse = {
  items: HypervisorNode[]
}

async function readJSON<T>(response: Response): Promise<T | null> {
  return (await response.json().catch(() => null)) as T | null
}

export async function listHypervisorNodes(): Promise<HypervisorNode[]> {
  const response = await fetch('/api/v1/admin/hypervisor/nodes', {
    credentials: 'include',
  })
  const payload = await readJSON<{
    message?: string
    data?: HypervisorNodesResponse
  }>(response)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load hypervisor nodes')
  }
  return payload?.data?.items || []
}

export async function getHypervisorNodeDetail(
  nodeID: string
): Promise<HypervisorNodeDetail> {
  const response = await fetch(
    `/api/v1/admin/hypervisor/nodes/${encodeURIComponent(nodeID)}`,
    {
      credentials: 'include',
    }
  )
  const payload = await readJSON<{
    message?: string
    data?: HypervisorNodeDetail
  }>(response)

  if (!response.ok || !payload?.data) {
    throw new Error(payload?.message || 'Failed to load hypervisor node')
  }
  return payload.data
}

export async function updateHypervisorNodeName(input: {
  nodeID: string
  name: string
}): Promise<void> {
  const response = await fetch(
    `/api/v1/admin/hypervisor/nodes/${encodeURIComponent(input.nodeID)}/name`,
    {
      method: 'PATCH',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name: input.name }),
    }
  )
  const payload = await readJSON<{ message?: string }>(response)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to update node name')
  }
}

export async function assignHypervisorNodeZone(input: {
  nodeID: string
  zoneID: string
}): Promise<void> {
  const response = await fetch(
    `/api/v1/admin/hypervisor/nodes/${encodeURIComponent(input.nodeID)}/zone`,
    {
      method: 'PATCH',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ zone_id: input.zoneID }),
    }
  )
  const payload = await readJSON<{ message?: string }>(response)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to assign node to zone')
  }
}

export function buildHypervisorNodeMetricsWebSocketURL(nodeID: string): string {
  if (typeof window === 'undefined') {
    return ''
  }

  const url = new URL(window.location.origin)
  url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:'
  url.pathname = '/api/v1/admin/hypervisor/metrics/ws'
  url.searchParams.set('id', nodeID)
  return url.toString()
}
