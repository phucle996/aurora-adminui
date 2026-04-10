type Envelope<T> = {
  message?: string
  data?: T
}

export type ResourceDefinitionRecord = {
  id: string
  resource_type: string
  resource_model: string
  resource_version: string
  display_name: string
  status: string
  resource_count: number
}

export type TemplateResourceDefinitionOptionRecord = {
  id: string
  resource_type: string
  resource_model: string
  resource_version: string
}

type ResourceDefinitionListResponse = {
  items: ResourceDefinitionRecord[]
}

type TemplateResourceDefinitionOptionsResponse = {
  items: TemplateResourceDefinitionOptionRecord[]
}

export async function listResourceDefinitions(): Promise<ResourceDefinitionRecord[]> {
  const response = await fetch('/api/v1/admin/resource-definitions', {
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<ResourceDefinitionListResponse>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load resource definitions')
  }
  return payload?.data?.items || []
}

export async function listTemplateResourceDefinitionOptions(): Promise<TemplateResourceDefinitionOptionRecord[]> {
  const response = await fetch('/api/v1/admin/resource-definitions/template-options', {
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<TemplateResourceDefinitionOptionsResponse>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load template resource definition options')
  }
  return payload?.data?.items || []
}

export async function createResourceDefinition(input: {
  resource_type: string
  resource_model: string
  resource_version: string
  display_name: string
}): Promise<void> {
  const response = await fetch('/api/v1/admin/resource-definitions', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  const payload = (await response.json().catch(() => null)) as Envelope<null> | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to create resource definition')
  }
}

export async function deleteResourceDefinition(definitionID: string): Promise<void> {
  const response = await fetch(`/api/v1/admin/resource-definitions/${encodeURIComponent(definitionID)}`, {
    method: 'DELETE',
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<null>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to delete resource definition')
  }
}

export async function updateResourceDefinitionStatus(
  definitionID: string,
  status: string
): Promise<ResourceDefinitionRecord> {
  const response = await fetch(`/api/v1/admin/resource-definitions/${encodeURIComponent(definitionID)}`, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ status }),
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<ResourceDefinitionRecord>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to update resource definition status')
  }
  if (!payload?.data) {
    throw new Error('Resource definition status was updated but no data was returned')
  }
  return payload.data
}
