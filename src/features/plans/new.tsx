import { useEffect, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, Boxes, Plus } from 'lucide-react'
import { toast } from 'sonner'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { listResourceDefinitions } from '@/features/resource-types/api'
import { createPlan } from './api'

export function AddPlanPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [resourceType, setResourceType] = useState('')
  const [resourceModel, setResourceModel] = useState('')
  const [code, setCode] = useState('')
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [vcpu, setVCPU] = useState('2')
  const [ramGb, setRamGb] = useState('4')
  const [diskGb, setDiskGb] = useState('50')

  const resourceDefinitionsQuery = useQuery({
    queryKey: ['resource-definitions'],
    queryFn: listResourceDefinitions,
  })

  const resourceTypes = useMemo(() => {
    const items = resourceDefinitionsQuery.data || []
    return [...new Set(items.map((item) => item.resource_type))].sort((left, right) =>
      left.localeCompare(right)
    )
  }, [resourceDefinitionsQuery.data])

  const resourceModels = useMemo(() => {
    const items = resourceDefinitionsQuery.data || []
    return items
      .filter((item) => item.resource_type === resourceType)
      .map((item) => item.resource_model)
      .filter((value, index, array) => array.indexOf(value) === index)
      .sort((left, right) => left.localeCompare(right))
  }, [resourceDefinitionsQuery.data, resourceType])

  useEffect(() => {
    if (!resourceType && resourceTypes[0]) {
      setResourceType(resourceTypes[0])
    }
  }, [resourceType, resourceTypes])

  useEffect(() => {
    if (resourceModels.length === 0) {
      if (resourceModel !== '') {
        setResourceModel('')
      }
      return
    }
    if (!resourceModels.includes(resourceModel)) {
      setResourceModel(resourceModels[0] || '')
    }
  }, [resourceModel, resourceModels])

  const createPlanMutation = useMutation({
    mutationFn: createPlan,
    onSuccess: (plan) => {
      toast.success(`Plan "${plan.name}" created`)
      queryClient.invalidateQueries({ queryKey: ['plans'] })
      navigate({ to: '/plans' })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to create plan'
      )
    },
  })

  function handleCreatePlan() {
    if (!resourceType || !resourceModel || !code.trim() || !name.trim()) {
      toast.error('Resource type, resource model, plan code, and name are required')
      return
    }

    createPlanMutation.mutate({
      resourceType,
      resourceModel,
      code: code.trim(),
      name: name.trim(),
      description: description.trim(),
      vcpu: Number(vcpu),
      ramGb: Number(ramGb),
      diskGb: Number(diskGb),
    })
  }

  return (
    <>
      <Header fixed>
        <div className='min-w-0'>
          <p className='subtle-kicker'>Plan catalog</p>
          <h1 className='truncate text-lg font-semibold text-foreground'>
            Add Plan
          </h1>
        </div>
        <div className='ms-auto flex items-center space-x-4'>
          <Search />
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main className='flex flex-col gap-6'>
        <section className='page-header'>
          <div className='space-y-2'>
            <p className='subtle-kicker'>Resource package</p>
            <h1 className='page-title'>Create a new plan</h1>
            <p className='page-copy'>
              Add a provisioning plan to the shared catalog so operators can
              use it for new resources.
            </p>
          </div>
          <Button variant='outline' asChild>
            <Link to='/plans'>
              <ArrowLeft className='size-4' />
              Back to plans
            </Link>
          </Button>
        </section>

        <div className='grid gap-6 xl:grid-cols-[minmax(0,420px)_minmax(0,1fr)]'>
          <Card>
            <CardHeader>
              <CardTitle>Plan definition</CardTitle>
            </CardHeader>
            <CardContent className='space-y-5'>
              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Resource type
                </label>
                <select
                  value={resourceType}
                  onChange={(event) => setResourceType(event.target.value)}
                  className='h-10 w-full rounded-md border border-input bg-background px-3 text-sm'
                  disabled={resourceDefinitionsQuery.isLoading}
                >
                  {resourceTypes.length > 0 ? (
                    resourceTypes.map((item) => (
                      <option key={item} value={item}>
                        {item}
                      </option>
                    ))
                  ) : (
                    <option value=''>
                      {resourceDefinitionsQuery.isLoading
                        ? 'Loading resource types...'
                        : 'No resource types found'}
                    </option>
                  )}
                </select>
              </div>

              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Resource model
                </label>
                <select
                  value={resourceModel}
                  onChange={(event) => setResourceModel(event.target.value)}
                  className='h-10 w-full rounded-md border border-input bg-background px-3 text-sm'
                  disabled={resourceDefinitionsQuery.isLoading || resourceModels.length === 0}
                >
                  {resourceModels.length > 0 ? (
                    resourceModels.map((item) => (
                      <option key={item} value={item}>
                        {item}
                      </option>
                    ))
                  ) : (
                    <option value=''>
                      {resourceDefinitionsQuery.isLoading
                        ? 'Loading resource models...'
                        : 'No resource models found'}
                    </option>
                  )}
                </select>
              </div>

              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Plan code
                </label>
                <Input
                  value={code}
                  onChange={(event) => setCode(event.target.value)}
                  placeholder='Example: vps-standard-2c4g'
                />
              </div>

              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Plan name
                </label>
                <Input
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder='Example: Standard 2C / 4G'
                />
              </div>

              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Description
                </label>
                <Textarea
                  className='min-h-28'
                  value={description}
                  onChange={(event) => setDescription(event.target.value)}
                  placeholder='Describe the intended workload or profile for this plan.'
                />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className='space-y-3'>
              <div className='flex items-center gap-3'>
                <span className='flex size-11 items-center justify-center rounded-2xl bg-accent text-accent-foreground'>
                  <Boxes className='size-5' />
                </span>
                <div>
                  <p className='subtle-kicker'>Resource sizing</p>
                  <CardTitle>VPS compute profile</CardTitle>
                </div>
              </div>
            </CardHeader>
            <CardContent className='space-y-5'>
              <div className='grid gap-4 md:grid-cols-3'>
                <div className='space-y-2'>
                  <label className='text-sm font-medium text-foreground'>
                    vCPU
                  </label>
                  <Input
                    type='number'
                    min={1}
                    value={vcpu}
                    onChange={(event) => setVCPU(event.target.value)}
                  />
                </div>

                <div className='space-y-2'>
                  <label className='text-sm font-medium text-foreground'>
                    RAM (GB)
                  </label>
                  <Input
                    type='number'
                    min={1}
                    value={ramGb}
                    onChange={(event) => setRamGb(event.target.value)}
                  />
                </div>

                <div className='space-y-2'>
                  <label className='text-sm font-medium text-foreground'>
                    Disk (GB)
                  </label>
                  <Input
                    type='number'
                    min={1}
                    value={diskGb}
                    onChange={(event) => setDiskGb(event.target.value)}
                  />
                </div>
              </div>

              <div className='rounded-xl border border-border/80 bg-muted/40 px-4 py-3'>
                <p className='text-sm font-medium text-foreground'>
                  Active resource shape
                </p>
                <p className='mt-1 text-sm text-muted-foreground'>
                  {vcpu || '0'} vCPU · {ramGb || '0'} GB RAM · {diskGb || '0'} GB
                  disk
                </p>
              </div>

              <div className='flex flex-wrap gap-3'>
                <Button
                  onClick={handleCreatePlan}
                  disabled={createPlanMutation.isPending}
                >
                  <Plus className='size-4' />
                  {createPlanMutation.isPending ? 'Creating...' : 'Create plan'}
                </Button>
                <Button variant='outline' asChild>
                  <Link to='/plans'>Cancel</Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </Main>
    </>
  )
}
