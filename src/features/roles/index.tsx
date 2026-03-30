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
import { listRoles } from './api'

export function RolesPage() {
  const rolesQuery = useQuery({
    queryKey: ['roles'],
    queryFn: listRoles,
  })

  const roles = rolesQuery.data || []

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
            <h2 className='text-2xl font-bold tracking-tight'>Role List</h2>
            <p className='text-muted-foreground'>
              Review IAM roles with only their scope, assignment count, and
              permission footprint.
            </p>
          </div>
          <Button asChild>
            <Link to='/roles/new'>
              <Plus className='size-4' />
              Add role
            </Link>
          </Button>
        </div>

        <div className='overflow-hidden rounded-md border bg-card'>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className='w-[240px]'>Name</TableHead>
                <TableHead className='w-[160px]'>Scope</TableHead>
                <TableHead className='w-[140px]'>Users</TableHead>
                <TableHead className='w-[160px]'>Permissions</TableHead>
                <TableHead>Description</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {rolesQuery.isLoading ? (
                <TableRow>
                  <TableCell colSpan={5} className='h-24 text-center'>
                    Loading roles...
                  </TableCell>
                </TableRow>
              ) : rolesQuery.isError ? (
                <TableRow>
                  <TableCell colSpan={5} className='h-24'>
                    <div className='flex items-start justify-center gap-3 text-sm text-warning'>
                      <ShieldAlert className='mt-0.5 size-4 shrink-0' />
                      <span>
                        {rolesQuery.error instanceof Error
                          ? rolesQuery.error.message
                          : 'Failed to load roles'}
                      </span>
                    </div>
                  </TableCell>
                </TableRow>
              ) : roles.length > 0 ? (
                roles.map((role) => (
                  <TableRow key={role.id}>
                    <TableCell>
                      <div className='space-y-1'>
                        <p className='font-medium text-foreground'>
                          {role.name}
                        </p>
                        <p className='font-mono text-xs text-muted-foreground'>
                          {role.id}
                        </p>
                      </div>
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {role.scope}
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {role.userCount}
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {role.permissionCount}
                    </TableCell>
                    <TableCell className='text-sm leading-6 text-muted-foreground'>
                      {role.description || '-'}
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={5} className='h-24 text-center'>
                    No roles found.
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
