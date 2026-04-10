import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  createResourceDefinition,
  deleteResourceDefinition,
  listResourceDefinitions,
  type ResourceDefinitionRecord,
  updateResourceDefinitionStatus,
} from './api'

type ResourceDefinitionForm = {
  resource_type: string
  resource_model: string
  resource_version: string
  display_name: string
}

type ResourceDefinitionField = keyof ResourceDefinitionForm

const requiredFields: Array<{
  key: ResourceDefinitionField
  label: string
}> = [
  { key: 'resource_type', label: 'Resource type' },
  { key: 'resource_model', label: 'Resource model' },
  { key: 'resource_version', label: 'Resource version' },
  { key: 'display_name', label: 'Display name' },
]

function sortRecords(items: ResourceDefinitionRecord[]) {
  return [...items].sort((left, right) => {
    if (left.resource_type !== right.resource_type) {
      return left.resource_type.localeCompare(right.resource_type)
    }
    if (left.resource_model !== right.resource_model) {
      return left.resource_model.localeCompare(right.resource_model)
    }
    if (left.resource_version !== right.resource_version) {
      return right.resource_version.localeCompare(left.resource_version)
    }
    return left.display_name.localeCompare(right.display_name)
  })
}

export function ResourceTypesPage() {
  const queryClient = useQueryClient()
  const [form, setForm] = useState<ResourceDefinitionForm>({
    resource_type: '',
    resource_model: '',
    resource_version: '',
    display_name: '',
  })
  const [fieldErrors, setFieldErrors] = useState<
    Partial<Record<ResourceDefinitionField, string>>
  >({})

  const listQuery = useQuery({
    queryKey: ['resource-definitions'],
    queryFn: listResourceDefinitions,
  })

  const createMutation = useMutation({
    mutationFn: createResourceDefinition,
    onSuccess: () => {
      toast.success('Resource definition created')
      queryClient.invalidateQueries({ queryKey: ['resource-definitions'] })
      setForm({
        resource_type: '',
        resource_model: '',
        resource_version: '',
        display_name: '',
      })
      setFieldErrors({})
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to create resource type'
      )
    },
  })

  const deleteMutation = useMutation({
    mutationFn: deleteResourceDefinition,
    onSuccess: () => {
      toast.success('Resource definition deleted')
      queryClient.invalidateQueries({ queryKey: ['resource-definitions'] })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to delete resource definition'
      )
    },
  })

  const statusMutation = useMutation({
    mutationFn: ({ definitionID, status }: { definitionID: string; status: string }) =>
      updateResourceDefinitionStatus(definitionID, status),
    onSuccess: (_, variables) => {
      toast.success(`Resource definition marked as ${variables.status}`)
      queryClient.invalidateQueries({ queryKey: ['resource-definitions'] })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to update resource definition status'
      )
    },
  })

  const rows = useMemo(
    () => sortRecords(listQuery.data || []),
    [listQuery.data]
  )

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const normalized = {
      resource_type: form.resource_type.trim(),
      resource_model: form.resource_model.trim(),
      resource_version: form.resource_version.trim(),
      display_name: form.display_name.trim(),
    }
    const nextErrors: Partial<Record<ResourceDefinitionField, string>> = {}
    requiredFields.forEach((field) => {
      if (!normalized[field.key]) {
        nextErrors[field.key] = `${field.label} is required`
      }
    })
    setFieldErrors(nextErrors)
    if (Object.keys(nextErrors).length > 0) {
      toast.error('Please fill the highlighted fields')
      return
    }
    await createMutation.mutateAsync(normalized)
  }

  function updateFormField(field: ResourceDefinitionField, value: string) {
    setForm((current) => ({
      ...current,
      [field]: value,
    }))
    setFieldErrors((current) => {
      if (!current[field]) {
        return current
      }
      const next = { ...current }
      delete next[field]
      return next
    })
  }

  function inputClassName(field: ResourceDefinitionField) {
    return fieldErrors[field]
      ? 'border-destructive focus-visible:ring-destructive'
      : undefined
  }

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
            <p className='subtle-kicker'>Resource platform</p>
            <h1 className='page-title'>Resource Definitions</h1>
            <p className='page-copy'>
              Define the catalog of resource definitions that template
              rendering and cluster routing can target.
            </p>
          </div>
        </section>

        <Card>
          <CardHeader>
            <CardTitle>New resource definition</CardTitle>
            <CardDescription>
              Example RD: type `database`, model `mysql`, version `8.4`.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form
              className='grid gap-5 md:grid-cols-2 xl:grid-cols-3'
              onSubmit={handleSubmit}
            >
              <div className='space-y-2'>
                <p className='text-sm font-medium text-foreground'>
                  Resource type
                </p>
                <Input
                  value={form.resource_type}
                  onChange={(event) =>
                    updateFormField('resource_type', event.target.value)
                  }
                  placeholder='database'
                  className={inputClassName('resource_type')}
                  aria-invalid={Boolean(fieldErrors.resource_type)}
                />
                {fieldErrors.resource_type ? (
                  <p className='text-xs font-medium text-destructive'>
                    {fieldErrors.resource_type}
                  </p>
                ) : null}
              </div>

              <div className='space-y-2'>
                <p className='text-sm font-medium text-foreground'>
                  Resource model
                </p>
                <Input
                  value={form.resource_model}
                  onChange={(event) =>
                    updateFormField('resource_model', event.target.value)
                  }
                  placeholder='mysql'
                  className={inputClassName('resource_model')}
                  aria-invalid={Boolean(fieldErrors.resource_model)}
                />
                {fieldErrors.resource_model ? (
                  <p className='text-xs font-medium text-destructive'>
                    {fieldErrors.resource_model}
                  </p>
                ) : null}
              </div>

              <div className='space-y-2'>
                <p className='text-sm font-medium text-foreground'>
                  Resource version
                </p>
                <Input
                  value={form.resource_version}
                  onChange={(event) =>
                    updateFormField('resource_version', event.target.value)
                  }
                  placeholder='7.0, 8.4'
                  className={inputClassName('resource_version')}
                  aria-invalid={Boolean(fieldErrors.resource_version)}
                />
                {fieldErrors.resource_version ? (
                  <p className='text-xs font-medium text-destructive'>
                    {fieldErrors.resource_version}
                  </p>
                ) : null}
              </div>

              <div className='space-y-2'>
                <p className='text-sm font-medium text-foreground'>
                  Display name
                </p>
                <Input
                  value={form.display_name}
                  onChange={(event) =>
                    updateFormField('display_name', event.target.value)
                  }
                  placeholder='MySQL 8.4'
                  className={inputClassName('display_name')}
                  aria-invalid={Boolean(fieldErrors.display_name)}
                />
                {fieldErrors.display_name ? (
                  <p className='text-xs font-medium text-destructive'>
                    {fieldErrors.display_name}
                  </p>
                ) : null}
              </div>

              <div className='md:col-span-2 xl:col-span-3'>
                <Button type='submit' disabled={createMutation.isPending}>
                  Create resource definition
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Catalog</CardTitle>
            <CardDescription>
              Each row is one exact RD record: one type, one model, one version.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {listQuery.isLoading ? (
              <div className='rounded-xl border border-border/70 bg-muted/20 px-4 py-10 text-center text-sm text-muted-foreground'>
                Loading resource definitions...
              </div>
            ) : rows.length > 0 ? (
              <div className='overflow-hidden rounded-md border bg-card'>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className='w-[140px]'>Type</TableHead>
                      <TableHead className='w-[180px]'>Model</TableHead>
                      <TableHead className='w-[120px]'>Version</TableHead>
                      <TableHead>Display name</TableHead>
                      <TableHead className='w-[140px]'>Status</TableHead>
                      <TableHead className='w-[110px]'>Resources</TableHead>
                      <TableHead className='w-[110px] text-right'>Delete</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {rows.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell className='text-sm font-medium text-foreground'>
                          {item.resource_type}
                        </TableCell>
                        <TableCell className='font-mono text-sm text-muted-foreground'>
                          {item.resource_model}
                        </TableCell>
                        <TableCell className='text-sm text-muted-foreground'>
                          {item.resource_version}
                        </TableCell>
                        <TableCell className='text-sm font-medium text-foreground'>
                          {item.display_name}
                        </TableCell>
                        <TableCell>
                          <select
                            value={item.status}
                            onChange={(event) =>
                              statusMutation.mutate({
                                definitionID: item.id,
                                status: event.target.value,
                              })
                            }
                            className='h-9 w-full rounded-md border border-input bg-background px-3 text-sm'
                            disabled={statusMutation.isPending}
                          >
                            <option value='draft'>draft</option>
                            <option value='ready'>ready</option>
                            <option value='maintain'>maintain</option>
                          </select>
                        </TableCell>
                        <TableCell className='text-sm text-muted-foreground'>
                          {item.resource_count}
                        </TableCell>
                        <TableCell className='text-right'>
                          <Button
                            variant='outline'
                            disabled={item.resource_count > 0 || deleteMutation.isPending}
                            onClick={() => deleteMutation.mutate(item.id)}
                          >
                            Delete
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            ) : (
              <div className='rounded-xl border border-border/70 bg-muted/20 px-4 py-10 text-center text-sm text-muted-foreground'>
                No resource definitions have been defined yet.
              </div>
            )}
          </CardContent>
        </Card>
      </Main>
    </>
  )
}
