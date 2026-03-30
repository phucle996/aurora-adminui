import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, MapPinned, Plus } from 'lucide-react'
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
import { Textarea } from '@/components/ui/textarea'
import { createZone } from './api'

export function AddZonePage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const createZoneMutation = useMutation({
    mutationFn: createZone,
    onSuccess: (zone) => {
      toast.success(`Zone "${zone.name}" created`)
      queryClient.invalidateQueries({ queryKey: ['zones'] })
      navigate({ to: '/zones' })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to create zone'
      )
    },
  })

  function handleCreateZone() {
    if (!name.trim()) {
      toast.error('Zone name is required')
      return
    }
    createZoneMutation.mutate({
      name: name.trim(),
      description: description.trim(),
    })
  }

  return (
    <>
      <Header fixed>
        <div className='min-w-0'>
          <p className='subtle-kicker'>Topology</p>
          <h1 className='truncate text-lg font-semibold text-foreground'>
            Add Zone
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
            <p className='subtle-kicker'>Zone catalog</p>
            <h1 className='page-title'>Create a new zone</h1>
            <p className='page-copy'>
              Keep the zone model minimal. A zone only needs a name and a clear
              description for operators.
            </p>
          </div>
          <Button variant='outline' asChild>
            <Link to='/zones'>
              <ArrowLeft className='size-4' />
              Back to zones
            </Link>
          </Button>
        </section>

        <Card className='max-w-3xl'>
          <CardHeader>
            <div className='flex items-center gap-3'>
              <span className='flex size-11 items-center justify-center rounded-2xl bg-accent text-accent-foreground'>
                <MapPinned className='size-5' />
              </span>
              <div>
                <p className='subtle-kicker'>New zone</p>
                <CardTitle>Zone definition</CardTitle>
                <CardDescription>
                  Create a zone with the minimal model used by the current
                  topology flow.
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className='space-y-5'>
            <div className='space-y-2'>
              <label className='text-sm font-medium text-foreground'>
                Zone name
              </label>
              <Input
                value={name}
                onChange={(event) => setName(event.target.value)}
                placeholder='Example: Ho Chi Minh Primary'
              />
            </div>

            <div className='space-y-2'>
              <label className='text-sm font-medium text-foreground'>
                Description
              </label>
              <Textarea
                className='min-h-32'
                value={description}
                onChange={(event) => setDescription(event.target.value)}
                placeholder='Describe what this zone is used for, where it sits in the topology, and what operators should expect from it.'
              />
            </div>

            <div className='flex flex-wrap gap-3'>
              <Button
                onClick={handleCreateZone}
                disabled={createZoneMutation.isPending}
              >
                <Plus className='size-4' />
                {createZoneMutation.isPending ? 'Creating...' : 'Create zone'}
              </Button>
              <Button variant='outline' asChild>
                <Link to='/zones'>Cancel</Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      </Main>
    </>
  )
}
