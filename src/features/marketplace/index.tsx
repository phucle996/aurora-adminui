import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { Eye, Rocket, ShieldAlert } from 'lucide-react'
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
import { listMarketplaceApps } from './api'

export function MarketplaceListPage() {
  const appsQuery = useQuery({
    queryKey: ['admin-marketplace-apps'],
    queryFn: listMarketplaceApps,
  })

  const apps = useMemo(() => appsQuery.data || [], [appsQuery.data])

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
            <h1 className='page-title'>Marketplace</h1>
            <p className='page-copy'>
              Curate application marketplace entries that map one app to one
              linked resource model and inherit its version.
            </p>
          </div>
          <Button asChild>
            <Link to='/marketplace/new'>
              <Rocket className='size-4' />
              Add new
            </Link>
          </Button>
        </section>

        {appsQuery.isError ? (
          <div className='flex items-start gap-3 rounded-xl border border-warning/25 bg-warning-soft px-4 py-4 text-sm text-warning'>
            <ShieldAlert className='mt-0.5 size-4 shrink-0' />
            <span>
              {appsQuery.error instanceof Error
                ? appsQuery.error.message
                : 'Failed to load marketplace apps'}
            </span>
          </div>
        ) : null}

        <div className='overflow-hidden rounded-md border bg-card'>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className='w-[220px]'>Name</TableHead>
                <TableHead className='w-[140px]'>Resource</TableHead>
                <TableHead className='w-[180px]'>Resource model</TableHead>
                <TableHead className='w-[220px]'>Versions</TableHead>
                <TableHead>Description</TableHead>
                <TableHead className='w-[110px]'>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {appsQuery.isLoading ? (
                <TableRow>
                  <TableCell colSpan={6} className='h-24 text-center text-sm text-muted-foreground'>
                    Loading marketplace catalog...
                  </TableCell>
                </TableRow>
              ) : apps.length > 0 ? (
                apps.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell>
                      <div className='space-y-1'>
                        <div className='font-medium text-foreground'>{item.name}</div>
                        <div className='font-mono text-xs text-muted-foreground'>
                          {item.slug}
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {item.resource_type}
                    </TableCell>
                    <TableCell className='font-mono text-sm text-muted-foreground'>
                      {item.resource_model}
                    </TableCell>
                    <TableCell className='max-w-[240px] whitespace-normal text-sm text-muted-foreground'>
                      {item.versions.join(', ') || '-'}
                    </TableCell>
                    <TableCell className='max-w-[480px] whitespace-normal text-sm leading-6 text-muted-foreground'>
                      {item.summary}
                    </TableCell>
                    <TableCell>
                      <Button variant='outline' size='sm' asChild>
                        <Link
                          to='/marketplace/$appId'
                          params={{ appId: item.id }}
                        >
                          <Eye className='size-4' />
                          Detail
                        </Link>
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={6} className='h-24 text-center text-sm text-muted-foreground'>
                    No marketplace apps have been defined yet.
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
