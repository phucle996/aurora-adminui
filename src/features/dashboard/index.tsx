import {
  Activity,
  ArrowRight,
  Boxes,
  ShieldCheck,
  TriangleAlert,
  Users,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'

export function Dashboard() {
  return (
    <>
      <Header fixed>
        <div className='min-w-0'>
          <p className='subtle-kicker'>Administration</p>
          <h1 className='truncate text-lg font-semibold text-foreground'>
            Aurora Platform Overview
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
            <p className='subtle-kicker'>Operations Console</p>
            <h1 className='page-title'>Control the platform from one surface</h1>
            <p className='page-copy'>
              This workspace is tuned for operators managing tenants,
              permissions, plans, and resource availability. The layout follows
              the same infrastructure-first visual language as the
              controlplane.
            </p>
          </div>
          <div className='flex flex-wrap items-center gap-3'>
            <Button>
              Open control queue
              <ArrowRight className='size-4' />
            </Button>
            <Button variant='outline'>Review access policy</Button>
          </div>
        </section>

        <section className='admin-card-grid sm:grid-cols-2 xl:grid-cols-4'>
          {summaryCards.map((card) => (
            <Card key={card.title}>
              <CardHeader className='flex flex-row items-start justify-between gap-3 space-y-0'>
                <div className='space-y-2'>
                  <p className='subtle-kicker'>{card.kicker}</p>
                  <CardTitle className='text-sm'>{card.title}</CardTitle>
                </div>
                <span className='flex size-10 items-center justify-center rounded-2xl bg-secondary text-muted-foreground'>
                  <card.icon className='size-5' />
                </span>
              </CardHeader>
              <CardContent className='space-y-3'>
                <div className='stat-value'>{card.value}</div>
                <div className='flex items-center justify-between gap-2'>
                  <Badge
                    variant={card.tone === 'healthy' ? 'secondary' : 'outline'}
                    className={
                      card.tone === 'healthy'
                        ? 'border-0 bg-success-soft text-success'
                        : card.tone === 'warning'
                          ? 'border-0 bg-warning-soft text-warning'
                          : 'border-0 bg-info-soft text-info'
                    }
                  >
                    {card.delta}
                  </Badge>
                  <p className='text-xs text-muted-foreground'>{card.note}</p>
                </div>
              </CardContent>
            </Card>
          ))}
        </section>

        <section className='grid gap-5 xl:grid-cols-[1.5fr_1fr]'>
          <Card>
            <CardHeader>
              <p className='subtle-kicker'>Platform posture</p>
              <CardTitle>Live control checks</CardTitle>
            </CardHeader>
            <CardContent className='space-y-4'>
              {platformChecks.map((item) => (
                <div
                  key={item.title}
                  className='flex items-start justify-between gap-4 rounded-xl border border-border/80 bg-muted/60 px-4 py-3'
                >
                  <div className='space-y-1'>
                    <div className='flex items-center gap-2'>
                      <item.icon className='size-4 text-muted-foreground' />
                      <p className='font-medium text-foreground'>{item.title}</p>
                    </div>
                    <p className='text-sm leading-6 text-muted-foreground'>
                      {item.description}
                    </p>
                  </div>
                  <Badge
                    className={
                      item.status === 'Healthy'
                        ? 'border-0 bg-success-soft text-success'
                        : item.status === 'Warning'
                          ? 'border-0 bg-warning-soft text-warning'
                          : 'border-0 bg-info-soft text-info'
                    }
                  >
                    {item.status}
                  </Badge>
                </div>
              ))}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <p className='subtle-kicker'>Next actions</p>
              <CardTitle>Operator focus</CardTitle>
            </CardHeader>
            <CardContent className='space-y-3'>
              {operatorFocus.map((item) => (
                <div
                  key={item.title}
                  className='rounded-xl border border-border/80 px-4 py-3'
                >
                  <div className='flex items-center justify-between gap-3'>
                    <p className='font-medium text-foreground'>{item.title}</p>
                    <Badge variant='outline'>{item.owner}</Badge>
                  </div>
                  <p className='mt-2 text-sm leading-6 text-muted-foreground'>
                    {item.description}
                  </p>
                </div>
              ))}
            </CardContent>
          </Card>
        </section>

        <section className='grid gap-5 xl:grid-cols-[1.2fr_1fr]'>
          <Card>
            <CardHeader>
              <p className='subtle-kicker'>Execution backlog</p>
              <CardTitle>Queued change windows</CardTitle>
            </CardHeader>
            <CardContent className='space-y-3'>
              {changeWindows.map((window) => (
                <div
                  key={window.name}
                  className='flex flex-col gap-3 rounded-xl border border-border/80 px-4 py-4 sm:flex-row sm:items-center sm:justify-between'
                >
                  <div className='space-y-1'>
                    <p className='font-medium text-foreground'>{window.name}</p>
                    <p className='text-sm text-muted-foreground'>
                      {window.scope}
                    </p>
                  </div>
                  <div className='flex items-center gap-3'>
                    <Badge
                      className={
                        window.state === 'Ready'
                          ? 'border-0 bg-success-soft text-success'
                          : window.state === 'Blocked'
                            ? 'border-0 bg-warning-soft text-warning'
                            : 'border-0 bg-info-soft text-info'
                      }
                    >
                      {window.state}
                    </Badge>
                    <span className='text-sm text-muted-foreground'>
                      {window.window}
                    </span>
                  </div>
                </div>
              ))}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <p className='subtle-kicker'>Coverage</p>
              <CardTitle>Administrative surface map</CardTitle>
            </CardHeader>
            <CardContent className='grid gap-3 sm:grid-cols-2 xl:grid-cols-1'>
              {surfaceMap.map((surface) => (
                <div
                  key={surface.label}
                  className='rounded-xl border border-border/80 bg-muted/60 px-4 py-4'
                >
                  <p className='text-sm font-medium text-foreground'>
                    {surface.label}
                  </p>
                  <p className='mt-2 text-3xl font-semibold tracking-tight text-foreground'>
                    {surface.value}
                  </p>
                  <p className='mt-1 text-xs text-muted-foreground'>
                    {surface.note}
                  </p>
                </div>
              ))}
            </CardContent>
          </Card>
        </section>
      </Main>
    </>
  )
}

const summaryCards = [
  {
    kicker: 'Identity',
    title: 'Privileged administrators',
    value: '128',
    delta: 'Stable',
    note: '6 pending role reviews',
    tone: 'info',
    icon: Users,
  },
  {
    kicker: 'Catalog',
    title: 'Active resource packages',
    value: '24',
    delta: '+3 this quarter',
    note: 'Aligned to production SKUs',
    tone: 'healthy',
    icon: Boxes,
  },
  {
    kicker: 'Security',
    title: 'Policy exceptions',
    value: '7',
    delta: 'Needs review',
    note: '2 expiring this week',
    tone: 'warning',
    icon: ShieldCheck,
  },
  {
    kicker: 'Runtime',
    title: 'Protected control surfaces',
    value: '11',
    delta: 'All online',
    note: 'No degraded admin paths',
    tone: 'healthy',
    icon: Activity,
  },
]

const platformChecks = [
  {
    title: 'Access boundary validation',
    description:
      'Role mappings, session guards, and elevated scopes are synchronized across the workspace.',
    status: 'Healthy',
    icon: ShieldCheck,
  },
  {
    title: 'Plan catalog hygiene',
    description:
      'Immutable packages are versioned cleanly and retired catalog entries are not exposed for provisioning.',
    status: 'Healthy',
    icon: Boxes,
  },
  {
    title: 'Operator review backlog',
    description:
      'Two changes are still waiting for final confirmation before they can be promoted to the next maintenance window.',
    status: 'Warning',
    icon: TriangleAlert,
  },
]

const operatorFocus = [
  {
    title: 'Retire legacy tenant bootstrap flow',
    description:
      'Finalize the new package-bound provisioning path before the next controlplane rollout window.',
    owner: 'Platform',
  },
  {
    title: 'Consolidate security notifications',
    description:
      'Route admin-facing warnings into one queue so approvals and follow-up actions are easier to triage.',
    owner: 'IAM',
  },
  {
    title: 'Promote catalog guardrails',
    description:
      'Keep non-production package drafts out of operator create flows until validation completes.',
    owner: 'Plan',
  },
]

const changeWindows = [
  {
    name: 'Controlplane configuration refresh',
    scope: 'Identity, package catalog, and provisioning policies',
    state: 'Ready',
    window: 'Today · 22:30 ICT',
  },
  {
    name: 'Admin workspace permissions review',
    scope: 'Operator groups and scoped access clean-up',
    state: 'Monitoring',
    window: 'Tomorrow · 09:00 ICT',
  },
  {
    name: 'Tenant onboarding hardening',
    scope: 'Input validation and protected flows',
    state: 'Blocked',
    window: 'Needs sign-off',
  },
]

const surfaceMap = [
  {
    label: 'Identity domains',
    value: '05',
    note: 'Centralized sign-in, roles, and approvals',
  },
  {
    label: 'Catalog families',
    value: '03',
    note: 'Packages prepared for infrastructure rollout',
  },
  {
    label: 'Protected routes',
    value: '18',
    note: 'Authenticated operator experiences',
  },
  {
    label: 'Data stores',
    value: '04',
    note: 'Backoffice state, audit, cache, and control metadata',
  },
]
