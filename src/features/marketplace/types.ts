export type MarketplaceAppRecord = {
  id: string
  name: string
  slug: string
  summary: string
  description: string
  resource_definition_id: string
  template_id: string
  template_name: string
  resource_type: string
  resource_model: string
  versions: string[]
}

export type MarketplaceAppInput = {
  name: string
  slug: string
  summary: string
  description: string
  resource_definition_id: string
  template_id: string
}
