export type Zone = {
  id: string
  name: string
  description: string
  resource_count: number
  can_delete: boolean
}

type ListZonesResponse = {
  items: Zone[]
}

type Envelope<T> = {
  message?: string
  data?: T
}

export async function listZones(): Promise<Zone[]> {
  const response = await fetch('/api/v1/admin/zones', {
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<ListZonesResponse>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load zones')
  }
  return payload?.data?.items || []
}

export async function createZone(input: {
  name: string
  description: string
}): Promise<Zone> {
  const response = await fetch('/api/v1/admin/zones', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<Zone>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to create zone')
  }
  if (!payload?.data) {
    throw new Error('Zone was created but no response data was returned')
  }
  return payload.data
}

export async function deleteZone(id: string): Promise<void> {
  const response = await fetch(`/api/v1/admin/zones/${id}`, {
    method: 'DELETE',
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<null>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to delete zone')
  }
}
