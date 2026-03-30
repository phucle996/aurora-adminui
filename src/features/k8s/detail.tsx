import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, RefreshCw, ShieldAlert, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
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
import {
  deleteK8sCluster,
  getK8sClusterDetail,
  revalidateK8sCluster,
} from './api'

function formatDateTime(value?: string) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleString([], {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

function statusTone(status: string) {
  switch (status) {
    case 'valid':
      return 'bg-emerald-500'
    case 'invalid':
      return 'bg-rose-500'
    case 'unreachable':
      return 'bg-amber-500'
    default:
      return 'bg-slate-400'
  }
}

function InfoRow(props: { label: string; value: string }) {
  return (
    <div className='grid gap-1 rounded-xl border border-border/70 bg-muted/20 px-4 py-3'>
      <span className='text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground'>
        {props.label}
      </span>
      <span className='text-sm font-medium text-foreground break-all'>
        {props.value}
      </span>
    </div>
  )
}

export function K8sClusterDetailPage(props: { clusterID: string }) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const detailQuery = useQuery({
    queryKey: ['k8s-cluster', props.clusterID],
    queryFn: () => getK8sClusterDetail(props.clusterID),
  })

  const revalidateMutation = useMutation({
    mutationFn: revalidateK8sCluster,
    onSuccess: (cluster) => {
      toast.success(`Revalidated ${cluster.name}`)
      queryClient.invalidateQueries({ queryKey: ['k8s-clusters'] })
      queryClient.setQueryData(['k8s-cluster', cluster.id], cluster)
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to revalidate cluster'
      )
    },
  })

  const deleteMutation = useMutation({
    mutationFn: deleteK8sCluster,
    onSuccess: () => {
      toast.success('Cluster deleted')
      queryClient.invalidateQueries({ queryKey: ['k8s-clusters'] })
      navigate({ to: '/k8s' })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to delete cluster'
      )
    },
  })

  const cluster = detailQuery.data

  function handleDelete() {
    if (!cluster) return
    if (!window.confirm(`Delete cluster "${cluster.name}"?`)) {
      return
    }
    deleteMutation.mutate(cluster.id)
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
        <section className='page-header'>
          <div className='space-y-2'>
            <p className='subtle-kicker'>Kubernetes substrate</p>
            <h1 className='page-title'>
              {cluster?.name || 'Kubernetes Cluster Detail'}
            </h1>
            <p className='page-copy'>
              Review cluster identity, access metadata, capability flags, and
              the latest validation result stored by the platform.
            </p>
          </div>
          <div className='flex flex-wrap gap-3'>
            <Button variant='outline' asChild>
              <Link to='/k8s'>
                <ArrowLeft className='size-4' />
                Back to clusters
              </Link>
            </Button>
            <Button
              variant='outline'
              onClick={() => revalidateMutation.mutate(props.clusterID)}
              disabled={revalidateMutation.isPending}
            >
              <RefreshCw className='size-4' />
              Revalidate
            </Button>
            <Button
              variant='destructive'
              onClick={handleDelete}
              disabled={deleteMutation.isPending || detailQuery.isLoading}
            >
              <Trash2 className='size-4' />
              Delete
            </Button>
          </div>
        </section>

        {detailQuery.isLoading ? (
          <div className='rounded-md border bg-card p-8 text-sm text-muted-foreground'>
            Loading cluster detail...
          </div>
        ) : detailQuery.isError || !cluster ? (
          <div className='rounded-md border bg-card p-8'>
            <div className='flex items-start gap-3 text-sm text-warning'>
              <ShieldAlert className='mt-0.5 size-4 shrink-0' />
              <span>
                {detailQuery.error instanceof Error
                  ? detailQuery.error.message
                  : 'Failed to load kubernetes cluster'}
              </span>
            </div>
          </div>
        ) : (
          <div className='grid gap-6 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)]'>
            <Card>
              <CardHeader className='space-y-3'>
                <div className='flex items-center gap-3'>
                  <span
                    className={`size-2.5 rounded-full ${statusTone(cluster.validation_status)}`}
                  />
                  <CardTitle>{cluster.name}</CardTitle>
                  <Badge variant='secondary' className='rounded-full px-3 py-1'>
                    {cluster.validation_status}
                  </Badge>
                </div>
                <CardDescription>
                  {cluster.description || 'No description provided.'}
                </CardDescription>
              </CardHeader>
              <CardContent className='grid gap-3 md:grid-cols-2'>
                <InfoRow label='Zone' value={cluster.zone_name || '-'} />
                <InfoRow label='Import mode' value={cluster.import_mode} />
                <InfoRow
                  label='API server'
                  value={cluster.api_server_url || '-'}
                />
                <InfoRow
                  label='Current context'
                  value={cluster.current_context || '-'}
                />
                <InfoRow
                  label='Kubernetes version'
                  value={cluster.kubernetes_version || '-'}
                />
                <InfoRow
                  label='Last validated'
                  value={formatDateTime(cluster.last_validated_at)}
                />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Capability flags</CardTitle>
                <CardDescription>
                  These flags decide which future platform products can target
                  this cluster.
                </CardDescription>
              </CardHeader>
              <CardContent className='space-y-3'>
                <div className='flex flex-wrap gap-2'>
                  <Badge
                    variant={cluster.supports_dbaas ? 'default' : 'secondary'}
                    className='rounded-full px-3 py-1'
                  >
                    DBaaS
                  </Badge>
                  <Badge
                    variant={
                      cluster.supports_serverless ? 'default' : 'secondary'
                    }
                    className='rounded-full px-3 py-1'
                  >
                    Serverless
                  </Badge>
                  <Badge
                    variant={
                      cluster.supports_generic_workloads
                        ? 'default'
                        : 'secondary'
                    }
                    className='rounded-full px-3 py-1'
                  >
                    Generic workloads
                  </Badge>
                </div>
                <div className='rounded-xl border border-border/70 bg-muted/20 px-4 py-4 text-sm leading-7 text-muted-foreground'>
                  Use `revalidate` any time to check whether the stored kubeconfig
                  can still reach the Kubernetes API server and perform the
                  lightweight namespace read used by the platform validator.
                </div>
              </CardContent>
            </Card>

            <Card className='xl:col-span-2'>
              <CardHeader>
                <CardTitle>Validation summary</CardTitle>
                <CardDescription>
                  The latest validation run determines whether this substrate is
                  usable, misconfigured, or temporarily unreachable.
                </CardDescription>
              </CardHeader>
              <CardContent className='grid gap-3 md:grid-cols-3'>
                <InfoRow label='Status' value={cluster.validation_status} />
                <InfoRow
                  label='Created at'
                  value={formatDateTime(cluster.created_at)}
                />
                <InfoRow
                  label='Zone binding'
                  value={cluster.zone_name || 'Unbound'}
                />
                <div className='md:col-span-3 rounded-xl border border-border/70 bg-muted/20 px-4 py-4'>
                  <p className='text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground'>
                    Last validation error
                  </p>
                  <p className='mt-2 text-sm leading-7 text-foreground'>
                    {cluster.last_validation_error || 'No validation error.'}
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </Main>
    </>
  )
}
