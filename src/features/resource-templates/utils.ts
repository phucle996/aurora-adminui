import { parseAllDocuments } from 'yaml'
import type { TemplateRenderValidation } from './types'

export function validateTemplateYaml(raw: string): TemplateRenderValidation {
  if (!raw.trim()) {
    return {
      valid: false,
      message: 'Template YAML is required.',
      documentCount: 0,
    }
  }

  const documents = parseAllDocuments(raw)
  const firstError = documents.flatMap((document) => document.errors)[0]
  if (firstError) {
    return {
      valid: false,
      message: firstError.message,
      documentCount: documents.length,
    }
  }

  return {
    valid: true,
    message: 'YAML syntax looks valid.',
    documentCount: documents.length,
  }
}

export function extractTemplatePlaceholders(raw: string) {
  const matches = raw.matchAll(/{{\s*([^{}]+?)\s*}}/g)
  return Array.from(new Set(Array.from(matches, (match) => match[1].trim())))
}

export function formatTemplateDateTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleString([], {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}
