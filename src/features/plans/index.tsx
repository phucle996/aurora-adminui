import { useMemo, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { ArrowDown, ArrowUp, ChevronsUpDown, Plus, ShieldAlert } from 'lucide-react'
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
import { listPlans } from './api'

type SortKey = 'resource_type' | 'status' | 'vcpu' | 'ram_gb' | 'disk_gb'
type SortDirection = 'desc' | 'asc'

export function PlansPage() {
  const [sortState, setSortState] = useState<{
    key: SortKey
    direction: SortDirection
  } | null>(null)

  const plansQuery = useQuery({
    queryKey: ['plans'],
    queryFn: listPlans,
  })

  const plans = useMemo(() => {
    const items = [...(plansQuery.data || [])]
    if (!sortState) {
      return items
    }

    return items.sort((left, right) => {
      const leftValue = left[sortState.key]
      const rightValue = right[sortState.key]

      if (typeof leftValue === 'number' && typeof rightValue === 'number') {
        return sortState.direction === 'desc'
          ? rightValue - leftValue
          : leftValue - rightValue
      }

      const compared = String(leftValue).localeCompare(String(rightValue))
      return sortState.direction === 'desc' ? -compared : compared
    })
  }, [plansQuery.data, sortState])

  function toggleSort(key: SortKey) {
    setSortState((current) => {
      if (!current || current.key !== key) {
        return { key, direction: 'desc' }
      }
      if (current.direction === 'desc') {
        return { key, direction: 'asc' }
      }
      return null
    })
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
            <h2 className='text-2xl font-bold tracking-tight'>Plan List</h2>
            <p className='text-muted-foreground'>
              Review resource plans from the platform catalog with only the
              fields used by operators during provisioning.
            </p>
          </div>
          <Button asChild>
            <Link to='/plans/new'>
              <Plus className='size-4' />
              Add plan
            </Link>
          </Button>
        </div>

        <div className='overflow-hidden rounded-md border bg-card'>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className='w-[220px]'>Name</TableHead>
                <TableHead className='w-[160px]'>Code</TableHead>
                <SortableHead
                  className='w-[140px]'
                  label='Resource'
                  active={sortState?.key === 'resource_type' ? sortState.direction : null}
                  onClick={() => toggleSort('resource_type')}
                />
                <SortableHead
                  className='w-[120px]'
                  label='Status'
                  active={sortState?.key === 'status' ? sortState.direction : null}
                  onClick={() => toggleSort('status')}
                />
                <SortableHead
                  className='w-[100px]'
                  label='vCPU'
                  active={sortState?.key === 'vcpu' ? sortState.direction : null}
                  onClick={() => toggleSort('vcpu')}
                />
                <SortableHead
                  className='w-[100px]'
                  label='RAM'
                  active={sortState?.key === 'ram_gb' ? sortState.direction : null}
                  onClick={() => toggleSort('ram_gb')}
                />
                <SortableHead
                  className='w-[100px]'
                  label='Disk'
                  active={sortState?.key === 'disk_gb' ? sortState.direction : null}
                  onClick={() => toggleSort('disk_gb')}
                />
                <TableHead>Description</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {plansQuery.isLoading ? (
                <TableRow>
                  <TableCell colSpan={8} className='h-24 text-center'>
                    Loading plans...
                  </TableCell>
                </TableRow>
              ) : plansQuery.isError ? (
                <TableRow>
                  <TableCell colSpan={8} className='h-24'>
                    <div className='flex items-start justify-center gap-3 text-sm text-warning'>
                      <ShieldAlert className='mt-0.5 size-4 shrink-0' />
                      <span>
                        {plansQuery.error instanceof Error
                          ? plansQuery.error.message
                          : 'Failed to load plans'}
                      </span>
                    </div>
                  </TableCell>
                </TableRow>
              ) : plans.length > 0 ? (
                plans.map((plan) => (
                  <TableRow key={plan.id}>
                    <TableCell className='font-medium text-foreground'>
                      {plan.name}
                    </TableCell>
                    <TableCell className='font-mono text-xs text-muted-foreground'>
                      {plan.code}
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {plan.resource_type}
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {plan.status}
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {plan.vcpu}
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {plan.ram_gb} GB
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {plan.disk_gb} GB
                    </TableCell>
                    <TableCell className='text-sm leading-6 text-muted-foreground'>
                      {plan.description || '-'}
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={8} className='h-24 text-center'>
                    No plans found.
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

function SortableHead(props: {
  label: string
  className?: string
  active: SortDirection | null
  onClick: () => void
}) {
  return (
    <TableHead className={props.className}>
      <button
        type='button'
        onClick={props.onClick}
        className='inline-flex items-center gap-1 text-left text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground transition hover:text-foreground'
      >
        <span>{props.label}</span>
        {props.active === 'desc' ? (
          <ArrowDown className='size-3.5' />
        ) : props.active === 'asc' ? (
          <ArrowUp className='size-3.5' />
        ) : (
          <ChevronsUpDown className='size-3.5' />
        )}
      </button>
    </TableHead>
  )
}
