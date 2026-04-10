import { type ReactNode, useEffect, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import { ArrowLeft } from 'lucide-react'
import { toast } from 'sonner'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  createMarketplaceApp,
  listMarketplaceModelOptions,
  listMarketplaceTemplateOptions,
} from './api'
import { marketplaceSlug, validateMarketplaceInput } from './utils'

export function NewMarketplacePage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [form, setForm] = useState({
    name: '',
    slug: '',
    summary: '',
    description: '',
    resource_definition_id: '',
    template_id: '',
  })

  const definitionsQuery = useQuery({
    queryKey: ['marketplace-model-options'],
    queryFn: listMarketplaceModelOptions,
  })

  const templatesQuery = useQuery({
    queryKey: ['marketplace-template-options'],
    queryFn: listMarketplaceTemplateOptions,
  })

  const resourceModels = useMemo(() => {
    const options = definitionsQuery.data || []
    return [...options].sort((left, right) =>
      left.resource_model.localeCompare(right.resource_model)
    )
  }, [definitionsQuery.data])

  const selectedDefinition = useMemo(
    () =>
      resourceModels.find(
        (item) => item.resource_definition_id === form.resource_definition_id
      ) || null,
    [form.resource_definition_id, resourceModels]
  )

  const templateOptions = useMemo(() => {
    const selectedModel = selectedDefinition?.resource_model || ''
    return (templatesQuery.data || []).filter(
      (item) => item.resource_model === selectedModel
    )
  }, [selectedDefinition?.resource_model, templatesQuery.data])

  const selectedTemplate = useMemo(
    () =>
      templateOptions.find((item) => item.id === form.template_id) || null,
    [form.template_id, templateOptions]
  )

  useEffect(() => {
    if (!form.resource_definition_id && resourceModels[0]) {
      setForm((current) => ({
        ...current,
        resource_definition_id: resourceModels[0].resource_definition_id,
      }))
    }
  }, [form.resource_definition_id, resourceModels])

  useEffect(() => {
    if (!templateOptions.length) {
      if (form.template_id) {
        setForm((current) => ({ ...current, template_id: '' }))
      }
      return
    }
    const selectedStillExists = templateOptions.some(
      (item) => item.id === form.template_id
    )
    if (!selectedStillExists) {
      setForm((current) => ({
        ...current,
        template_id: templateOptions[0].id,
      }))
    }
  }, [form.template_id, templateOptions])

  const createMutation = useMutation({
    mutationFn: async () => {
      const input = {
        name: form.name,
        slug: marketplaceSlug(form.slug || form.name),
        summary: form.summary,
        description: form.description,
        resource_definition_id: form.resource_definition_id,
        template_id: form.template_id,
      }

      const validationError = validateMarketplaceInput(input)
      if (validationError) {
        throw new Error(validationError)
      }
      return createMarketplaceApp(input)
    },
    onSuccess: (item) => {
      toast.success('Marketplace app created')
      queryClient.invalidateQueries({ queryKey: ['admin-marketplace-apps'] })
      navigate({
        to: '/marketplace/$appId',
        params: { appId: item.id },
      })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to create marketplace app'
      )
    },
  })

  return (
    <>
      <Header fixed>
        <Search />
        <div className='ms-auto flex items-center space-x-4'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main className='flex flex-1 flex-col gap-4 sm:gap-6'>
        <section className='page-header'>
          <div className='space-y-2'>
            <p className='subtle-kicker'>Marketplace</p>
            <h1 className='page-title'>Add marketplace app</h1>
            <p className='page-copy'>
              Define one marketplace app entry, link it to one resource model,
              choose one template render, and inherit the supported versions from
              that linked model.
            </p>
          </div>
          <Button variant='outline' asChild>
            <Link to='/marketplace'>
              <ArrowLeft className='size-4' />
              Back to list
            </Link>
          </Button>
        </section>

        <Card className='rounded-2xl border-border/80'>
          <CardHeader>
            <CardTitle>Marketplace identity</CardTitle>
            <CardDescription>
              The linked resource model decides the version set, and the linked
              template decides how the app will render at deploy time.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form
              className='grid gap-5 md:grid-cols-2'
              onSubmit={async (event) => {
                event.preventDefault()
                await createMutation.mutateAsync()
              }}
            >
              <Field label='App name'>
                <Input
                  value={form.name}
                  onChange={(event) =>
                    setForm((current) => ({
                      ...current,
                      name: event.target.value,
                      slug: current.slug || marketplaceSlug(event.target.value),
                    }))
                  }
                  placeholder='n8n'
                />
              </Field>

              <Field label='Slug'>
                <Input
                  value={form.slug}
                  onChange={(event) =>
                    setForm((current) => ({
                      ...current,
                      slug: marketplaceSlug(event.target.value),
                    }))
                  }
                  placeholder='n8n'
                />
              </Field>

              <Field label='Linked resource model'>
                <select
                  value={form.resource_definition_id}
                  onChange={(event) => {
                    const nextModel = event.target.value
                    const option = resourceModels.find(
                      (item) =>
                        item.resource_definition_id === nextModel
                    )
                    setForm((current) => ({
                      ...current,
                      resource_definition_id:
                        option?.resource_definition_id ||
                        current.resource_definition_id,
                    }))
                  }}
                  className='h-10 w-full rounded-md border border-input bg-background px-3 text-sm'
                >
                  {resourceModels.map((item) => (
                    <option
                      key={item.resource_definition_id}
                      value={item.resource_definition_id}
                    >
                      {item.resource_type} / {item.resource_model}
                    </option>
                  ))}
                </select>
              </Field>

              <Field label='Linked template render'>
                <select
                  value={form.template_id}
                  onChange={(event) =>
                    setForm((current) => ({
                      ...current,
                      template_id: event.target.value,
                    }))
                  }
                  className='h-10 w-full rounded-md border border-input bg-background px-3 text-sm'
                  disabled={!selectedDefinition || templatesQuery.isLoading}
                >
                  {templateOptions.map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.name} ({item.version})
                    </option>
                  ))}
                </select>
              </Field>

              <Field className='md:col-span-2' label='Summary'>
                <Input
                  value={form.summary}
                  onChange={(event) =>
                    setForm((current) => ({
                      ...current,
                      summary: event.target.value,
                    }))
                  }
                  placeholder='Workflow automation app for event-driven pipelines.'
                />
              </Field>

              <Field className='md:col-span-2' label='Description'>
                <textarea
                  rows={5}
                  value={form.description}
                  onChange={(event) =>
                    setForm((current) => ({
                      ...current,
                      description: event.target.value,
                    }))
                  }
                  placeholder='Describe what the marketplace entry deploys and how it should be used.'
                  className='w-full rounded-md border border-input bg-background px-3 py-2 text-sm'
                />
              </Field>

              <Field className='md:col-span-2' label='Inherited versions'>
                <div className='flex flex-wrap gap-2 rounded-md border border-input bg-background px-3 py-3'>
                  {selectedDefinition?.versions.length ? (
                    selectedDefinition.versions.map((version) => {
                      return (
                        <span
                          key={version}
                          className='inline-flex items-center gap-2 rounded-full border border-border/70 px-3 py-1.5 text-sm'
                        >
                          <span>{version}</span>
                        </span>
                      )
                    })
                  ) : (
                    <span className='text-sm text-muted-foreground'>
                      No versions available for this resource model.
                    </span>
                  )}
                </div>
              </Field>

              <Field className='md:col-span-2' label='Selected template'>
                <div className='rounded-md border border-input bg-background px-3 py-3 text-sm text-muted-foreground'>
                  {selectedTemplate ? (
                    <span>
                      {selectedTemplate.name} for {selectedTemplate.resource_model}{' '}
                      {selectedTemplate.version}
                    </span>
                  ) : (
                    <span>No template render is available for this resource model.</span>
                  )}
                </div>
              </Field>

              <div className='md:col-span-2 flex items-center gap-3'>
                <Button type='submit' disabled={createMutation.isPending}>
                  Create marketplace app
                </Button>
                <p className='text-sm text-muted-foreground'>
                  Resource, resource model, and deploy template are resolved from
                  the linked resource definition and template render.
                </p>
              </div>
            </form>
          </CardContent>
        </Card>
      </Main>
    </>
  )
}

function Field(props: {
  label: string
  className?: string
  children: ReactNode
}) {
  return (
    <div className={props.className ? `space-y-2 ${props.className}` : 'space-y-2'}>
      <p className='text-sm font-medium text-foreground'>{props.label}</p>
      {props.children}
    </div>
  )
}
