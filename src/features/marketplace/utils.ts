import type { MarketplaceAppInput } from './types'

export function marketplaceSlug(value: string) {
  return value
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9-]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .replace(/--+/g, '-')
}

export function validateMarketplaceInput(input: MarketplaceAppInput) {
  if (!input.name.trim()) {
    return 'App name is required'
  }
  if (!marketplaceSlug(input.slug)) {
    return 'App slug is required'
  }
  if (!input.resource_definition_id.trim()) {
    return 'Linked resource definition is required'
  }
  if (!input.template_id.trim()) {
    return 'Linked template is required'
  }
  return null
}
