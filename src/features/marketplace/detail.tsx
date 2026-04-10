import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { ArrowLeft, Boxes, ShieldAlert } from 'lucide-react'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { getMarketplaceApp } from './api'

export function MarketplaceDetailPage(props: { appID: string }) {
  const detailQuery = useQuery({
    queryKey: ['admin-marketplace-apps', props.appID],
    queryFn: async () => getMarketplaceApp(props.appID),
  })

  const app = detailQuery.data

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
            <h1 className='page-title'>{app?.name || 'Marketplace detail'}</h1>
            <p className='page-copy'>
              Review the marketplace entry, the linked resource model, and the
              linked template render it will use.
            </p>
          </div>
          <Button variant='outline' asChild>
            <Link to='/marketplace'>
              <ArrowLeft className='size-4' />
              Back to list
            </Link>
          </Button>
        </section>

        {detailQuery.isLoading ? (
          <div className='rounded-xl border border-border/80 bg-muted/50 px-4 py-8 text-sm text-muted-foreground'>
            Loading marketplace app...
          </div>
        ) : detailQuery.isError || !app ? (
          <div className='flex items-start gap-3 rounded-xl border border-warning/25 bg-warning-soft px-4 py-4 text-sm text-warning'>
            <ShieldAlert className='mt-0.5 size-4 shrink-0' />
            <span>
              {detailQuery.error instanceof Error
                ? detailQuery.error.message
                : 'Marketplace app not found'}
            </span>
          </div>
        ) : (
          <div className='grid gap-6 xl:grid-cols-[minmax(0,1.6fr)_360px]'>
            <div className='space-y-6'>
              <Card className='rounded-2xl border-border/80'>
                <CardHeader>
                  <CardTitle>App profile</CardTitle>
                  <CardDescription>
                    Identity and linked resource platform model for this app.
                  </CardDescription>
                </CardHeader>
                <CardContent className='grid gap-4 md:grid-cols-2'>
                  <InfoRow label='Name' value={app.name} />
                  <InfoRow label='Slug' value={app.slug} mono />
                  <InfoRow label='Resource type' value={app.resource_type} />
                  <InfoRow label='Resource model' value={app.resource_model} mono />
                  <InfoRow label='Template render' value={app.template_name} />
                </CardContent>
              </Card>

              <Card className='rounded-2xl border-border/80'>
                <CardHeader>
                  <CardTitle>Description</CardTitle>
                </CardHeader>
                <CardContent className='space-y-4 text-sm leading-7 text-muted-foreground'>
                  <p>{app.summary}</p>
                  <p>{app.description}</p>
                </CardContent>
              </Card>
            </div>

            <aside className='space-y-6 xl:sticky xl:top-24 xl:self-start'>
              <Card className='rounded-2xl border-border/80'>
                <CardHeader>
                  <CardTitle className='flex items-center gap-2'>
                    <Boxes className='size-4 text-muted-foreground' />
                    Model version
                  </CardTitle>
                  <CardDescription>
                    This marketplace entry uses the version that comes from the
                    configured app catalog.
                  </CardDescription>
                </CardHeader>
                <CardContent className='flex flex-wrap gap-2'>
                  {app.versions.map((version) => (
                    <Badge key={version} variant='default'>
                      {version}
                    </Badge>
                  ))}
                </CardContent>
              </Card>
            </aside>
          </div>
        )}
      </Main>
    </>
  )
}

function InfoRow(props: { label: string; value: string; mono?: boolean }) {
  return (
    <div className='grid gap-1 rounded-xl border border-border/70 bg-muted/20 px-4 py-3'>
      <span className='text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground'>
        {props.label}
      </span>
      <span className={props.mono ? 'break-all font-mono text-sm font-medium text-foreground' : 'break-all text-sm font-medium text-foreground'}>
        {props.value}
      </span>
    </div>
  )
}
