export type TemplateRenderRecord = {
  id: string
  name: string
  description: string
  resource_definition_id: string
  resource_type: string
  resource_model: string
  stream_key: string
  consumer_group: string
  yaml_template: string
  created_at: string
  updated_at: string
}

export type TemplateRenderInput = Omit<
  TemplateRenderRecord,
  'id' | 'created_at' | 'updated_at' | 'resource_type' | 'resource_model'
>

export type TemplateRenderValidation = {
  valid: boolean
  message: string
  documentCount: number
}
