import { useEffect, useState } from 'react'
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  deleteK8sCluster,
  getK8sClusterDetailPageData,
  revalidateK8sCluster,
  updateK8sCluster,
} from './api'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

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
      <span className='break-all text-sm font-medium text-foreground'>
        {props.value}
      </span>
    </div>
  )
}

export function K8sClusterDetailPage(props: { clusterID: string }) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [zoneID, setZoneID] = useState('')

  const pageDataQuery = useQuery({
    queryKey: ['k8s-cluster-page', props.clusterID],
    queryFn: () => getK8sClusterDetailPageData(props.clusterID),
  })

  const revalidateMutation = useMutation({
    mutationFn: revalidateK8sCluster,
    onSuccess: (cluster) => {
      toast.success(`Revalidated ${cluster.name}`)
      queryClient.invalidateQueries({ queryKey: ['k8s-clusters'] })
      queryClient.invalidateQueries({ queryKey: ['k8s-cluster-page', cluster.id] })
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

  const updateMutation = useMutation({
    mutationFn: (input: { zone_id: string }) =>
      updateK8sCluster(props.clusterID, input),
    onSuccess: (cluster) => {
      toast.success(`Updated ${cluster.name}`)
      queryClient.invalidateQueries({ queryKey: ['k8s-clusters'] })
      queryClient.invalidateQueries({ queryKey: ['k8s-cluster-page', cluster.id] })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to update cluster'
      )
    },
  })

  const cluster = pageDataQuery.data
  const pageErrorMessage =
    pageDataQuery.error instanceof Error
      ? pageDataQuery.error.message
      : 'Failed to load kubernetes cluster'

  useEffect(() => {
    if (!cluster) return
    setZoneID(cluster.zone_id || '')
  }, [cluster])

  function handleDelete() {
    if (!cluster) return
    if (!window.confirm(`Delete cluster "${cluster.name}"?`)) {
      return
    }
    deleteMutation.mutate(cluster.id)
  }

  function handleSaveSettings() {
    if (!zoneID) {
      toast.error('Please choose a zone for this cluster')
      return
    }
    updateMutation.mutate({
      zone_id: zoneID,
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
        <section className='page-header'>
          <div className='space-y-2'>
            <p className='subtle-kicker'>Kubernetes substrate</p>
            <h1 className='page-title'>
              {cluster?.name || 'Kubernetes Cluster Detail'}
            </h1>
            <p className='page-copy'>
              Review cluster identity, live node inventory, and zone placement
              for this execution substrate.
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
              disabled={revalidateMutation.isPending || pageDataQuery.isFetching}
            >
              <RefreshCw className='size-4' />
              Revalidate
            </Button>
            <Button
              variant='destructive'
              onClick={handleDelete}
              disabled={deleteMutation.isPending || pageDataQuery.isLoading}
            >
              <Trash2 className='size-4' />
              Delete
            </Button>
          </div>
        </section>

        {pageDataQuery.isLoading ? (
          <div className='rounded-md border bg-card p-8 text-sm text-muted-foreground'>
            Loading cluster detail...
          </div>
        ) : pageDataQuery.isError || !cluster ? (
          <div className='rounded-md border bg-card p-8'>
            <div className='flex items-start gap-3 text-sm text-warning'>
              <ShieldAlert className='mt-0.5 size-4 shrink-0' />
              <span>
                {pageErrorMessage}
              </span>
            </div>
          </div>
        ) : (
          <div className='grid gap-6 xl:grid-cols-[minmax(0,1.05fr)_minmax(0,0.95fr)]'>
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
                <InfoRow
                  label='Created at'
                  value={formatDateTime(cluster.created_at)}
                />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Cluster zone</CardTitle>
                <CardDescription>
                  Assign this cluster to exactly one zone. Resource definition
                  support is configured on the zone, not on the cluster.
                </CardDescription>
              </CardHeader>
              <CardContent className='grid gap-4'>
                <div className='grid gap-2'>
                  <label className='text-sm font-medium text-foreground'>
                    Assign to zone
                  </label>
                  <Select
                    value={zoneID || '__none__'}
                    onValueChange={(value) =>
                      setZoneID(value === '__none__' ? '' : value)
                    }
                  >
                    <SelectTrigger className='w-full'>
                      <SelectValue placeholder='Choose a zone' />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value='__none__' disabled>
                        Choose a zone
                      </SelectItem>
                      {(cluster.zone_options || []).map((zone) => (
                        <SelectItem key={zone.id} value={zone.id}>
                          {zone.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className='flex flex-wrap gap-3'>
                  <Button
                    onClick={handleSaveSettings}
                    disabled={
                      updateMutation.isPending || pageDataQuery.isLoading
                    }
                  >
                    Save cluster zone
                  </Button>
                  <Badge variant='secondary' className='rounded-full px-3 py-1'>
                    Zone: {cluster.zone_name || 'Unbound'}
                  </Badge>
                </div>
              </CardContent>
            </Card>

            <Card className='xl:col-span-2'>
              <CardHeader>
                <CardTitle>Cluster nodes</CardTitle>
                <CardDescription>
                  Live node inventory resolved through controlplane using the
                  stored kubeconfig for this cluster.
                </CardDescription>
              </CardHeader>
              <CardContent>
                {pageDataQuery.isLoading ? (
                  <div className='rounded-xl border border-dashed border-border/70 bg-muted/20 px-4 py-4 text-sm text-muted-foreground'>
                    Loading cluster nodes...
                  </div>
                ) : pageDataQuery.isError ? (
                  <div className='rounded-xl border border-dashed border-border/70 bg-muted/20 px-4 py-4 text-sm text-warning'>
                    {pageErrorMessage}
                  </div>
                ) : (cluster.nodes || []).length === 0 ? (
                  <div className='rounded-xl border border-dashed border-border/70 bg-muted/20 px-4 py-4 text-sm text-muted-foreground'>
                    No nodes were returned by the cluster API.
                  </div>
                ) : (
                  <div className='rounded-xl border border-border/70 bg-muted/20 p-2'>
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>Name</TableHead>
                          <TableHead>Status</TableHead>
                          <TableHead>Roles</TableHead>
                          <TableHead>Kubelet</TableHead>
                          <TableHead>Runtime</TableHead>
                          <TableHead>OS image</TableHead>
                          <TableHead>Kernel</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {(cluster.nodes || []).map((node) => (
                          <TableRow key={node.name}>
                            <TableCell className='font-medium text-foreground'>
                              {node.name}
                            </TableCell>
                            <TableCell>
                              <Badge
                                variant={node.ready ? 'default' : 'secondary'}
                                className='rounded-full px-2.5 py-0.5 text-[11px]'
                              >
                                {node.ready ? 'Ready' : 'Not ready'}
                              </Badge>
                            </TableCell>
                            <TableCell className='max-w-56 whitespace-normal text-muted-foreground'>
                              {node.roles.join(', ') || 'worker'}
                            </TableCell>
                            <TableCell>{node.kubelet_version || '-'}</TableCell>
                            <TableCell className='max-w-56 whitespace-normal text-muted-foreground'>
                              {node.container_runtime || '-'}
                            </TableCell>
                            <TableCell className='max-w-72 whitespace-normal text-muted-foreground'>
                              {node.os_image || '-'}
                            </TableCell>
                            <TableCell>{node.kernel_version || '-'}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        )}
      </Main>
    </>
  )
}
