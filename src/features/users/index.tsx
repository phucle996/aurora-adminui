import { useQuery } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { ShieldAlert } from 'lucide-react'
import { UsersDialogs } from './components/users-dialogs'
import { UsersPrimaryButtons } from './components/users-primary-buttons'
import { UsersProvider } from './components/users-provider'
import { UsersTable } from './components/users-table'
import { listUsers } from './api'

const route = getRouteApi('/_authenticated/users/')

export function Users() {
  const search = route.useSearch()
  const navigate = route.useNavigate()
  const usersQuery = useQuery({
    queryKey: ['users'],
    queryFn: listUsers,
  })

  return (
    <UsersProvider>
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
            <h2 className='text-2xl font-bold tracking-tight'>User List</h2>
            <p className='text-muted-foreground'>
              Manage your users and their roles here.
            </p>
          </div>
          <UsersPrimaryButtons />
        </div>
        {usersQuery.isLoading ? (
          <div className='rounded-xl border border-border/80 bg-muted/50 px-4 py-8 text-sm text-muted-foreground'>
            Loading users...
          </div>
        ) : usersQuery.isError ? (
          <div className='flex items-start gap-3 rounded-xl border border-warning/25 bg-warning-soft px-4 py-4 text-sm text-warning'>
            <ShieldAlert className='mt-0.5 size-4 shrink-0' />
            <span>
              {usersQuery.error instanceof Error
                ? usersQuery.error.message
                : 'Failed to load users'}
            </span>
          </div>
        ) : (
          <UsersTable
            data={usersQuery.data || []}
            search={search}
            navigate={navigate}
          />
        )}
      </Main>

      <UsersDialogs />
    </UsersProvider>
  )
}
