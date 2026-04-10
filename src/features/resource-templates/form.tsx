import { useQuery } from '@tanstack/react-query'
import { useEffect, useMemo, useState } from 'react'
import { AlertCircle, CheckCircle2, PlayCircle } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { listTemplateResourceDefinitionOptions } from '@/features/resource-types/api'
import { validateTemplateYaml } from './utils'
import type { TemplateRenderInput, TemplateRenderRecord } from './types'

const EMPTY_TEMPLATE: TemplateRenderInput = {
  name: '',
  description: '',
  resource_definition_id: '',
  stream_key: '',
  consumer_group: '',
  yaml_template: `apiVersion: v1
kind: ConfigMap
metadata:
  name: "{{ spec.name }}"
  namespace: "{{ spec.namespace }}"
data:
  job-id: "{{ job_id }}"
  resource-id: "{{ resource_id }}"
`,
}

function FieldLabel(props: { label: string; hint?: string }) {
  return (
    <div className='space-y-1'>
      <p className='text-sm font-medium text-foreground'>{props.label}</p>
      {props.hint ? (
        <p className='text-sm text-muted-foreground'>{props.hint}</p>
      ) : null}
    </div>
  )
}

export function TemplateRenderForm(props: {
  initialValue?: TemplateRenderRecord
  submitLabel: string
  busy?: boolean
  onSubmit: (value: TemplateRenderInput) => Promise<void> | void
}) {
  const [form, setForm] = useState<TemplateRenderInput>(
    props.initialValue ? toInput(props.initialValue) : EMPTY_TEMPLATE
  )
  const definitionsQuery = useQuery({
    queryKey: ['template-resource-definition-options'],
    queryFn: listTemplateResourceDefinitionOptions,
  })

  useEffect(() => {
    if (!props.initialValue) return
    setForm(toInput(props.initialValue))
  }, [props.initialValue])

  const definitions = definitionsQuery.data || []
  // Keep edit mode stable even if the local form state hasn't rehydrated yet.
  const effectiveResourceDefinitionID =
    form.resource_definition_id || props.initialValue?.resource_definition_id || ''
  const currentDefinition =
    definitions.find((item) => item.id === effectiveResourceDefinitionID) || null
  const resourceTypeOptions = useMemo(
    () => Array.from(new Set(definitions.map((item) => item.resource_type))),
    [definitions]
  )
  const selectedResourceType = currentDefinition?.resource_type || ''
  const modelOptions = useMemo(
    () =>
      definitions.filter((item) => item.resource_type === selectedResourceType),
    [definitions, selectedResourceType]
  )
  const modelGroupOptions = useMemo(
    () => Array.from(new Set(modelOptions.map((item) => deriveModelGroup(item)).filter(Boolean))),
    [modelOptions]
  )
  const selectedModelName = deriveModelGroup(currentDefinition) || ''
  const versionOptions = useMemo(
    () => modelOptions.filter((item) => deriveModelGroup(item) === selectedModelName),
    [modelOptions, selectedModelName]
  )

  useEffect(() => {
    if (props.initialValue || definitions.length === 0) {
      return
    }
    const hasMatch = definitions.some(
      (item) => item.id === form.resource_definition_id
    )
    if (!hasMatch) {
      setForm((current) => ({
        ...current,
        resource_definition_id: definitions[0]?.id || '',
      }))
    }
  }, [definitions, form.resource_definition_id, props.initialValue])

  const yamlValidation = useMemo(
    () => validateTemplateYaml(form.yaml_template),
    [form.yaml_template]
  )

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!form.name.trim()) {
      return
    }
    if (!yamlValidation.valid) {
      return
    }
    await props.onSubmit({
      ...form,
      name: form.name.trim(),
      description: form.description.trim(),
      resource_definition_id: effectiveResourceDefinitionID.trim(),
      stream_key: '',
      consumer_group: '',
      yaml_template: form.yaml_template,
    })
  }

  const selectedVersion = currentDefinition?.resource_version || ''
  const selectedVersionValue = selectedVersion || '__none__'

  return (
    <form className='grid gap-6 xl:grid-cols-[minmax(0,1.7fr)_360px]' onSubmit={handleSubmit}>
      <div className='space-y-6'>
        <section className='rounded-2xl border border-border/80 bg-card p-6 shadow-sm'>
          <div className='grid gap-5 md:grid-cols-2'>
            <div className='space-y-2 md:col-span-2'>
              <FieldLabel
                label='Template name'
                hint='Operator-facing name shown in the resource platform catalog.'
              />
              <Input
                value={form.name}
                onChange={(event) =>
                  setForm((current) => ({ ...current, name: event.target.value }))
                }
                placeholder='CNPG PostgreSQL 17 Cluster'
              />
            </div>
            <div className='space-y-2 md:col-span-2'>
              <FieldLabel
                label='Description'
                hint='Describe what kind of operator CR or workload this template renders.'
              />
              <Textarea
                rows={3}
                value={form.description}
                onChange={(event) =>
                  setForm((current) => ({
                    ...current,
                    description: event.target.value,
                  }))
                }
                placeholder='Renders a database operator manifest from the resource job payload.'
              />
            </div>
            <div className='space-y-2'>
              <FieldLabel label='Resource type' />
              <Select
                value={selectedResourceType}
                onValueChange={(value) => {
                  const nextDefinition =
                    definitions.find((item) => item.resource_type === value) || null
                  setForm((current) => ({
                    ...current,
                    resource_definition_id: nextDefinition?.id || '',
                  }))
                }}
              >
                <SelectTrigger className='w-full'>
                  <SelectValue placeholder='Select resource type' />
                </SelectTrigger>
                <SelectContent>
                  {resourceTypeOptions.map((item) => (
                    <SelectItem key={item} value={item}>
                      {item}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className='space-y-2'>
              <FieldLabel label='Resource model' />
              <Select
                value={selectedModelName}
                onValueChange={(value) => {
                  const nextDefinition =
                    definitions.find(
                      (item) =>
                        item.resource_type === selectedResourceType &&
                        deriveModelGroup(item) === value
                    ) || null
                  setForm((current) => ({
                    ...current,
                    resource_definition_id: nextDefinition?.id || '',
                  }))
                }}
                disabled={modelOptions.length === 0}
              >
                <SelectTrigger className='w-full'>
                  <SelectValue placeholder='Select resource model' />
                </SelectTrigger>
                <SelectContent>
                  {modelGroupOptions.map((item) => (
                    <SelectItem key={item} value={item}>
                      {item}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className='space-y-2'>
              <FieldLabel label='Resource version' />
              <Select
                value={selectedVersionValue}
                onValueChange={(value) =>
                  setForm((current) => ({
                    ...current,
                    resource_definition_id:
                      definitions.find(
                        (item) =>
                          item.resource_type === selectedResourceType &&
                          deriveModelGroup(item) === selectedModelName &&
                          (item.resource_version || '__none__') === value
                      )?.id || '',
                  }))
                }
                disabled={versionOptions.length === 0}
              >
                <SelectTrigger className='w-full'>
                  <SelectValue placeholder='Select resource version' />
                </SelectTrigger>
                <SelectContent>
                  {versionOptions.map((item) => (
                    <SelectItem key={item.id} value={item.resource_version || '__none__'}>
                      {item.resource_version || '-'}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        </section>

        <section className='rounded-2xl border border-border/80 bg-card p-6 shadow-sm'>
          <div className='space-y-4'>
            <div className='flex flex-wrap items-start justify-between gap-3'>
              <div className='space-y-1'>
                <FieldLabel
                  label='Template YAML'
                  hint='Write the YAML manifest with placeholders that dataplane workers render from the Redis message.'
                />
              </div>
              <Badge
                variant='outline'
                className={
                  yamlValidation.valid
                    ? 'border-emerald-200 bg-emerald-50 text-emerald-700'
                    : 'border-rose-200 bg-rose-50 text-rose-700'
                }
              >
                {yamlValidation.valid ? (
                  <CheckCircle2 className='size-3.5' />
                ) : (
                  <AlertCircle className='size-3.5' />
                )}
                {yamlValidation.valid ? 'YAML valid' : 'YAML invalid'}
              </Badge>
            </div>

            <Textarea
              rows={22}
              className='font-mono text-xs leading-6'
              value={form.yaml_template}
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  yaml_template: event.target.value,
                }))
              }
              placeholder='apiVersion: ...'
            />

            <div
              className={`rounded-2xl border px-4 py-3 text-sm ${
                yamlValidation.valid
                  ? 'border-emerald-200 bg-emerald-50 text-emerald-700'
                  : 'border-rose-200 bg-rose-50 text-rose-700'
              }`}
            >
              {yamlValidation.message}
            </div>
          </div>
        </section>
      </div>

      <aside className='space-y-6 xl:sticky xl:top-24 xl:self-start'>
        <section className='rounded-2xl border border-border/80 bg-card p-6 shadow-sm'>
          <div className='space-y-4'>
            <div>
              <p className='subtle-kicker'>Template summary</p>
              <h2 className='text-lg font-semibold text-foreground'>
                Render contract
              </h2>
            </div>
            <SummaryRow label='Resource type' value={selectedResourceType || '-'} />
            <SummaryRow label='Resource model' value={selectedModelName || '-'} />
            <SummaryRow label='Resource version' value={selectedVersion || '-'} />
            <SummaryRow label='Runtime stream' value='Resolved by cluster when a resource job is created' />
            <SummaryRow
              label='YAML documents'
              value={String(yamlValidation.documentCount || 0)}
            />
            <Button type='submit' className='w-full' disabled={props.busy || !yamlValidation.valid}>
              <PlayCircle className='size-4' />
              {props.submitLabel}
            </Button>
          </div>
        </section>

      </aside>
    </form>
  )
}

function SummaryRow(props: { label: string; value: string }) {
  return (
    <div className='flex items-start justify-between gap-4 border-b border-dashed border-border/70 pb-3 text-sm last:border-b-0 last:pb-0'>
      <span className='text-muted-foreground'>{props.label}</span>
      <span className='text-right font-medium text-foreground break-all'>
        {props.value}
      </span>
    </div>
  )
}

function toInput(template: TemplateRenderRecord): TemplateRenderInput {
  return {
    name: template.name,
    description: template.description,
    resource_definition_id: template.resource_definition_id,
    stream_key: template.stream_key,
    consumer_group: template.consumer_group,
    yaml_template: template.yaml_template,
  }
}

function deriveModelGroup(
  item:
    | {
        resource_model: string
        resource_version: string
      }
    | null
    | undefined
) {
  if (!item) return ''
  const version = item.resource_version?.trim()
  const model = item.resource_model?.trim()
  if (!version || !model) return model || ''

  const candidates = [`-${version}`, `-${version.split('.').join('-')}`]
  for (const suffix of candidates) {
    if (model.endsWith(suffix) && model.length > suffix.length) {
      return model.slice(0, -suffix.length)
    }
  }
  return model
}
