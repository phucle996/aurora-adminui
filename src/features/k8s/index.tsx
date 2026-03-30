import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { Eye, Plus, RefreshCw, ShieldAlert, Trash2 } from 'lucide-react'
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  deleteK8sCluster,
  listK8sClusters,
  revalidateK8sCluster,
} from './api'

function formatRelativeDateTime(value?: string) {
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

function validationTone(status: string) {
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

function capabilityBadges(cluster: {
  supports_dbaas: boolean
  supports_serverless: boolean
  supports_generic_workloads: boolean
}) {
  return [
    cluster.supports_dbaas ? 'DBaaS' : '',
    cluster.supports_serverless ? 'Serverless' : '',
    cluster.supports_generic_workloads ? 'Generic workloads' : '',
  ].filter(Boolean)
}

export function K8sPlatformPage() {
  const queryClient = useQueryClient()
  const clustersQuery = useQuery({
    queryKey: ['k8s-clusters'],
    queryFn: listK8sClusters,
  })

  const revalidateMutation = useMutation({
    mutationFn: revalidateK8sCluster,
    onSuccess: (cluster) => {
      toast.success(`Revalidated ${cluster.name}`)
      queryClient.invalidateQueries({ queryKey: ['k8s-clusters'] })
      queryClient.invalidateQueries({
        queryKey: ['k8s-cluster', cluster.id],
      })
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
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to delete cluster'
      )
    },
  })

  const clusters = clustersQuery.data || []

  function handleDeleteCluster(id: string, name: string) {
    if (!window.confirm(`Delete cluster "${name}"?`)) {
      return
    }
    deleteMutation.mutate(id)
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
            <p className='subtle-kicker'>Resource substrate</p>
            <h1 className='page-title'>K8s Resource Platform</h1>
            <p className='page-copy'>
              Register Kubernetes clusters as generic execution substrates for
              future DBaaS, serverless, and operator-managed platform resources.
            </p>
          </div>
          <Button asChild>
            <Link to='/k8s/new'>
              <Plus className='size-4' />
              Add cluster
            </Link>
          </Button>
        </section>

        <div className='overflow-hidden rounded-md border bg-card'>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className='w-[250px]'>Name</TableHead>
                <TableHead className='w-[140px]'>Zone</TableHead>
                <TableHead className='w-[160px]'>Status</TableHead>
                <TableHead className='w-[130px]'>Version</TableHead>
                <TableHead className='w-[280px]'>API server</TableHead>
                <TableHead className='w-[240px]'>Capabilities</TableHead>
                <TableHead className='w-[190px]'>Last validated</TableHead>
                <TableHead className='w-[180px] text-right'>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {clustersQuery.isLoading ? (
                <TableRow>
                  <TableCell colSpan={8} className='h-24 text-center'>
                    Loading clusters...
                  </TableCell>
                </TableRow>
              ) : clustersQuery.isError ? (
                <TableRow>
                  <TableCell colSpan={8} className='h-24'>
                    <div className='flex items-start justify-center gap-3 text-sm text-warning'>
                      <ShieldAlert className='mt-0.5 size-4 shrink-0' />
                      <span>
                        {clustersQuery.error instanceof Error
                          ? clustersQuery.error.message
                          : 'Failed to load kubernetes clusters'}
                      </span>
                    </div>
                  </TableCell>
                </TableRow>
              ) : clusters.length > 0 ? (
                clusters.map((cluster) => (
                  <TableRow key={cluster.id}>
                    <TableCell>
                      <div className='space-y-1'>
                        <p className='font-medium text-foreground'>
                          {cluster.name}
                        </p>
                        <p className='text-sm leading-6 text-muted-foreground'>
                          {cluster.description || '-'}
                        </p>
                      </div>
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {cluster.zone_name || '-'}
                    </TableCell>
                    <TableCell>
                      <div className='flex items-center gap-2 text-sm text-muted-foreground'>
                        <span
                          className={`size-2.5 rounded-full ${validationTone(cluster.validation_status)}`}
                        />
                        <span>{cluster.validation_status}</span>
                      </div>
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {cluster.kubernetes_version || '-'}
                    </TableCell>
                    <TableCell className='font-mono text-xs text-muted-foreground'>
                      {cluster.api_server_url}
                    </TableCell>
                    <TableCell>
                      <div className='flex flex-wrap gap-2'>
                        {capabilityBadges(cluster).length > 0 ? (
                          capabilityBadges(cluster).map((capability) => (
                            <Badge
                              key={capability}
                              variant='secondary'
                              className='rounded-full px-3 py-1'
                            >
                              {capability}
                            </Badge>
                          ))
                        ) : (
                          <span className='text-sm text-muted-foreground'>-</span>
                        )}
                      </div>
                    </TableCell>
                    <TableCell className='text-sm text-muted-foreground'>
                      {formatRelativeDateTime(cluster.last_validated_at)}
                    </TableCell>
                    <TableCell className='text-right'>
                      <div className='flex justify-end gap-2'>
                        <Button variant='ghost' size='icon' asChild>
                          <Link to='/k8s/$clusterId' params={{ clusterId: cluster.id }}>
                            <Eye className='size-4' />
                          </Link>
                        </Button>
                        <Button
                          variant='ghost'
                          size='icon'
                          onClick={() => revalidateMutation.mutate(cluster.id)}
                          disabled={revalidateMutation.isPending}
                        >
                          <RefreshCw className='size-4' />
                        </Button>
                        <Button
                          variant='ghost'
                          size='icon'
                          onClick={() =>
                            handleDeleteCluster(cluster.id, cluster.name)
                          }
                          disabled={deleteMutation.isPending}
                        >
                          <Trash2 className='size-4 text-destructive' />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={8} className='h-24 text-center'>
                    No kubernetes clusters found.
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
