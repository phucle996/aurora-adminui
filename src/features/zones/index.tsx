import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { Plus, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { deleteZone, listZones } from './api'

export function ZonesPage() {
  const queryClient = useQueryClient()
  const zonesQuery = useQuery({
    queryKey: ['zones'],
    queryFn: listZones,
  })
  const deleteZoneMutation = useMutation({
    mutationFn: deleteZone,
    onSuccess: () => {
      toast.success('Zone deleted')
      queryClient.invalidateQueries({ queryKey: ['zones'] })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to delete zone'
      )
    },
  })

  const zones = zonesQuery.data || []

  function handleDeleteZone(id: string, name: string, canDelete: boolean) {
    if (!canDelete) {
      toast.error('Zone still has attached nodes')
      return
    }
    if (!window.confirm(`Delete zone "${name}"?`)) {
      return
    }
    deleteZoneMutation.mutate(id)
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
        <div className='flex flex-wrap items-end justify-between gap-2'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>Zones</h2>
            <p className='text-muted-foreground'>
              Manage the list of placement zones with only a name and a short
              description.
            </p>
          </div>
          <Button asChild>
            <Link to='/zones/new'>
              <Plus className='size-4' />
              Add zone
            </Link>
          </Button>
        </div>

        <div className='overflow-hidden rounded-md border bg-card'>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className='w-[280px]'>Name</TableHead>
                <TableHead>Description</TableHead>
                <TableHead className='w-[120px]'>Nodes</TableHead>
                <TableHead className='w-[120px] text-right'>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {zonesQuery.isLoading ? (
                <TableRow>
                  <TableCell colSpan={4} className='h-24 text-center'>
                    Loading zones...
                  </TableCell>
                </TableRow>
              ) : zonesQuery.isError ? (
                <TableRow>
                  <TableCell colSpan={4} className='h-24 text-center'>
                    {zonesQuery.error instanceof Error
                      ? zonesQuery.error.message
                      : 'Failed to load zones'}
                  </TableCell>
                </TableRow>
              ) : zones.length > 0 ? (
                zones.map((zone) => (
                  <TableRow key={zone.id}>
                    <TableCell>
                      <div className='space-y-1'>
                        <p className='font-medium text-foreground'>
                          {zone.name}
                        </p>
                        <p className='font-mono text-xs text-muted-foreground'>
                          {zone.id}
                        </p>
                      </div>
                    </TableCell>
                    <TableCell className='text-sm leading-6 text-muted-foreground'>
                      {zone.description}
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {zone.resource_count}
                    </TableCell>
                    <TableCell className='text-right'>
                      <Button
                        variant='ghost'
                        size='icon'
                        disabled={
                          !zone.can_delete || deleteZoneMutation.isPending
                        }
                        onClick={() =>
                          handleDeleteZone(
                            zone.id,
                            zone.name,
                            zone.can_delete
                          )
                        }
                      >
                        <Trash2 className='size-4' />
                        <span className='sr-only'>Delete zone</span>
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={4} className='h-24 text-center'>
                    No zones defined.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </Main>
    </>
  )
}
