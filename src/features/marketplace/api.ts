import { z } from 'zod'

const marketplaceCatalogItemSchema = z.object({
  id: z.string(),
  name: z.string(),
  slug: z.string(),
  summary: z.string(),
  template_id: z.string(),
  template_name: z.string(),
  resource_type: z.string(),
  resource_model: z.string(),
  versions: z.array(z.string()),
})

const marketplaceModelOptionSchema = z.object({
  resource_definition_id: z.string(),
  resource_type: z.string(),
  resource_model: z.string(),
  versions: z.array(z.string()),
})

const marketplaceTemplateOptionSchema = z.object({
  id: z.string(),
  name: z.string(),
  resource_type: z.string(),
  resource_model: z.string(),
  version: z.string(),
})

const marketplaceDetailSchema = marketplaceCatalogItemSchema.extend({
  description: z.string(),
  resource_definition_id: z.string(),
})

const listMarketplaceResponseSchema = z.object({
  message: z.string(),
  data: z.object({
    items: z.array(marketplaceCatalogItemSchema),
  }),
})

const listMarketplaceModelOptionsResponseSchema = z.object({
  message: z.string(),
  data: z.object({
    items: z.array(marketplaceModelOptionSchema),
  }),
})

const listMarketplaceTemplateOptionsResponseSchema = z.object({
  message: z.string(),
  data: z.object({
    items: z.array(marketplaceTemplateOptionSchema),
  }),
})

const marketplaceItemResponseSchema = z.object({
  message: z.string(),
  data: marketplaceDetailSchema,
})

export type MarketplaceCatalogItem = z.infer<typeof marketplaceCatalogItemSchema>
export type MarketplaceApp = z.infer<typeof marketplaceDetailSchema>
export type MarketplaceModelOption = z.infer<typeof marketplaceModelOptionSchema>
export type MarketplaceTemplateOption = z.infer<typeof marketplaceTemplateOptionSchema>

export async function listMarketplaceModelOptions(): Promise<MarketplaceModelOption[]> {
  const response = await fetch('/api/v1/admin/marketplace/model-options', {
    credentials: 'include',
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load marketplace model options')
  }

  return listMarketplaceModelOptionsResponseSchema.parse(payload).data.items
}

export async function listMarketplaceTemplateOptions(): Promise<MarketplaceTemplateOption[]> {
  const response = await fetch('/api/v1/admin/marketplace/template-options', {
    credentials: 'include',
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load marketplace template options')
  }

  return listMarketplaceTemplateOptionsResponseSchema.parse(payload).data.items
}

export async function listMarketplaceApps(): Promise<MarketplaceCatalogItem[]> {
  const response = await fetch('/api/v1/admin/marketplace', {
    credentials: 'include',
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load marketplace apps')
  }

  return listMarketplaceResponseSchema.parse(payload).data.items
}

export async function getMarketplaceApp(id: string): Promise<MarketplaceApp> {
  const response = await fetch(`/api/v1/admin/marketplace/${id}`, {
    credentials: 'include',
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load marketplace app')
  }

  return marketplaceItemResponseSchema.parse(payload).data
}

export async function createMarketplaceApp(input: {
  name: string
  slug: string
  summary: string
  description: string
  resource_definition_id: string
  template_id: string
}): Promise<MarketplaceApp> {
  const response = await fetch('/api/v1/admin/marketplace', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to create marketplace app')
  }

  return marketplaceItemResponseSchema.parse(payload).data
}
