import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { Plus, ShieldAlert } from 'lucide-react'
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

export function PlansPage() {
  const plansQuery = useQuery({
    queryKey: ['plans'],
    queryFn: listPlans,
  })

  const plans = plansQuery.data || []

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
                <TableHead className='w-[140px]'>Resource</TableHead>
                <TableHead className='w-[120px]'>Status</TableHead>
                <TableHead className='w-[100px]'>vCPU</TableHead>
                <TableHead className='w-[100px]'>RAM</TableHead>
                <TableHead className='w-[100px]'>Disk</TableHead>
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
