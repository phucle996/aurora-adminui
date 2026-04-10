import type { TemplateRenderInput, TemplateRenderRecord } from './types'

type Envelope<T> = {
  message?: string
  data?: T
}

export type TemplateRenderCatalogRecord = {
  id: string
  name: string
  description: string
  resource_type: string
  resource_model: string
  stream_key: string
  consumer_group: string
  yaml_valid: boolean
  updated_at: string
}

type TemplateRenderCatalogResponse = {
  items: TemplateRenderCatalogRecord[]
}

export async function listTemplateRenderCatalog(): Promise<TemplateRenderCatalogRecord[]> {
  const response = await fetch('/api/v1/admin/resource-templates/catalog', {
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<TemplateRenderCatalogResponse>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load template render catalog')
  }
  return payload?.data?.items || []
}

export async function getTemplateRender(
  id: string
): Promise<TemplateRenderRecord> {
  const response = await fetch(`/api/v1/admin/resource-templates/${id}`, {
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<TemplateRenderRecord>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load template render')
  }
  if (!payload?.data) {
    throw new Error('Template render not found')
  }
  return payload.data
}

export async function createTemplateRender(
  input: TemplateRenderInput
): Promise<TemplateRenderRecord> {
  const response = await fetch('/api/v1/admin/resource-templates', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<TemplateRenderRecord>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to create template render')
  }
  if (!payload?.data) {
    throw new Error('Template render was created but no data was returned')
  }
  return payload.data
}

export async function updateTemplateRender(
  id: string,
  input: TemplateRenderInput
): Promise<TemplateRenderRecord> {
  const response = await fetch(`/api/v1/admin/resource-templates/${id}`, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<TemplateRenderRecord>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to update template render')
  }
  if (!payload?.data) {
    throw new Error('Template render was updated but no data was returned')
  }
  return payload.data
}

export async function deleteTemplateRender(id: string): Promise<void> {
  const response = await fetch(`/api/v1/admin/resource-templates/${id}`, {
    method: 'DELETE',
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<Record<string, never>>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to delete template render')
  }
}
