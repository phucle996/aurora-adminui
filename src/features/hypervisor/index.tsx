import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import {
  Check,
  Cpu,
  Eye,
  HardDrive,
  MemoryStick,
  MoreHorizontal,
  Pencil,
  Server,
  ShieldAlert,
  X,
} from 'lucide-react'
import { toast } from 'sonner'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Input } from '@/components/ui/input'
import {
  assignHypervisorNodeZone,
  listHypervisorNodes,
  type HypervisorNode,
  updateHypervisorNodeName,
} from './api'
import { listZones, type Zone } from '@/features/zones/api'

function formatMegabytes(megabytes: number) {
  if (!megabytes || megabytes <= 0) return '-'
  if (megabytes >= 1024) {
    const gigabytes = megabytes / 1024
    const fractionDigits = gigabytes >= 10 ? 0 : 1
    return `${gigabytes.toFixed(fractionDigits)} GB`
  }
  return `${Math.round(megabytes)} MB`
}

function statusDotTone(status: string) {
  switch (status.toLowerCase()) {
    case 'online':
      return 'bg-emerald-500'
    case 'offline':
      return 'bg-rose-500'
    default:
      return 'bg-amber-500'
  }
}

export function HypervisorNodes() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [editingNodeID, setEditingNodeID] = useState<string | null>(null)
  const [draftName, setDraftName] = useState('')
  const nodesQuery = useQuery({
    queryKey: ['hypervisor-nodes'],
    queryFn: listHypervisorNodes,
  })
  const zonesQuery = useQuery({
    queryKey: ['zones'],
    queryFn: listZones,
  })
  const renameMutation = useMutation({
    mutationFn: updateHypervisorNodeName,
    onSuccess: () => {
      toast.success('Node name updated')
      setEditingNodeID(null)
      setDraftName('')
      queryClient.invalidateQueries({ queryKey: ['hypervisor-nodes'] })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to update node name'
      )
    },
  })
  const assignZoneMutation = useMutation({
    mutationFn: assignHypervisorNodeZone,
    onSuccess: () => {
      toast.success('Node assigned to zone')
      queryClient.invalidateQueries({ queryKey: ['hypervisor-nodes'] })
      queryClient.invalidateQueries({ queryKey: ['zones'] })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error
          ? error.message
          : 'Failed to assign node to zone'
      )
    },
  })

  const nodes = nodesQuery.data || []
  const zones = zonesQuery.data || []

  function beginEdit(node: HypervisorNode) {
    setEditingNodeID(node.node_id)
    setDraftName(node.name || node.node_id)
  }

  function cancelEdit() {
    setEditingNodeID(null)
    setDraftName('')
  }

  function saveEdit(nodeID: string) {
    const nextName = draftName.trim()
    if (!nextName) {
      toast.error('Name is required')
      return
    }
    renameMutation.mutate({ nodeID, name: nextName })
  }

  function handleViewNode(node: HypervisorNode) {
    navigate({
      to: '/hypervisor/$nodeId',
      params: { nodeId: node.node_id },
    })
  }

  function handleAssignZone(node: HypervisorNode, zone: Zone) {
    assignZoneMutation.mutate({
      nodeID: node.node_id,
      zoneID: zone.id,
    })
  }

  return (
    <>
      <Header fixed>
        <div className='min-w-0'>
          <p className='subtle-kicker'>Infrastructure</p>
          <h1 className='truncate text-lg font-semibold text-foreground'>
            KVM Nodes
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
            <p className='subtle-kicker'>Node hardware</p>
            <h1 className='page-title'>KVM node configuration</h1>
            <p className='page-copy'>
              Review the registered KVM fleet by hardware profile and placement
              zone.
            </p>
          </div>
        </section>

        <section className='space-y-4'>
          <div className='flex flex-row items-center justify-between gap-4'>
            <div>
              <p className='subtle-kicker'>Node fleet</p>
              <h2 className='text-xl font-semibold tracking-tight text-foreground'>
                KVM hardware inventory
              </h2>
            </div>
            <span className='rounded-full border border-border/80 px-3 py-1 text-xs font-medium text-muted-foreground'>
              {nodes.length} nodes
            </span>
          </div>
          {nodesQuery.isLoading ? (
            <div className='rounded-xl border border-border/80 bg-muted/50 px-4 py-8 text-sm text-muted-foreground'>
              Loading hypervisor nodes...
            </div>
          ) : nodesQuery.isError ? (
            <div className='flex items-start gap-3 rounded-xl border border-warning/25 bg-warning-soft px-4 py-4 text-sm text-warning'>
              <ShieldAlert className='mt-0.5 size-4 shrink-0' />
              <span>
                {nodesQuery.error instanceof Error
                  ? nodesQuery.error.message
                  : 'Failed to load hypervisor nodes'}
              </span>
            </div>
          ) : nodes.length === 0 ? (
            <div className='rounded-xl border border-border/80 bg-muted/50 px-4 py-8 text-sm text-muted-foreground'>
              No KVM nodes found in `hypervisor.nodes`.
            </div>
          ) : (
            <div className='overflow-hidden rounded-2xl border border-border/80'>
              <div className='overflow-x-auto'>
                <table className='w-full min-w-[980px] text-left text-sm'>
                  <thead className='bg-muted/50 text-muted-foreground'>
                    <tr>
                      <th className='px-4 py-3 font-medium'>Name</th>
                      <th className='px-4 py-3 font-medium'>Zone</th>
                      <th className='px-4 py-3 font-medium'>CPU</th>
                      <th className='px-4 py-3 font-medium'>RAM</th>
                      <th className='px-4 py-3 font-medium'>Disk</th>
                      <th className='px-4 py-3 font-medium'>GPU</th>
                      <th className='px-4 py-3 text-right font-medium'>
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {nodes.map((node) => (
                      <tr
                        key={node.node_id}
                        className='border-t border-border/80 bg-background align-top'
                      >
                        <td className='px-4 py-3'>
                          <div className='group flex items-start gap-3'>
                            <span
                              className={`mt-3 size-2.5 shrink-0 rounded-full ${statusDotTone(node.status)}`}
                            />
                            <span className='mt-0.5 flex size-9 items-center justify-center rounded-xl bg-secondary text-muted-foreground'>
                              <Server className='size-4' />
                            </span>
                            <div className='min-w-0 space-y-1'>
                              {editingNodeID === node.node_id ? (
                                <div className='space-y-2'>
                                  <Input
                                    value={draftName}
                                    onChange={(event) =>
                                      setDraftName(event.target.value)
                                    }
                                    className='h-8'
                                    autoFocus
                                  />
                                  <div className='flex items-center gap-2'>
                                    <Button
                                      size='icon'
                                      className='size-8'
                                      onClick={() => saveEdit(node.node_id)}
                                      disabled={renameMutation.isPending}
                                    >
                                      <Check className='size-4' />
                                    </Button>
                                    <Button
                                      size='icon'
                                      variant='outline'
                                      className='size-8'
                                      onClick={cancelEdit}
                                      disabled={renameMutation.isPending}
                                    >
                                      <X className='size-4' />
                                    </Button>
                                  </div>
                                </div>
                              ) : (
                                <>
                                  <div className='flex items-center gap-2'>
                                    <p className='truncate font-medium text-foreground'>
                                      {node.name || node.node_id}
                                    </p>
                                    <button
                                      type='button'
                                      onClick={() => beginEdit(node)}
                                      className='opacity-0 transition-opacity group-hover:opacity-100'
                                    >
                                      <Pencil className='size-3.5 text-muted-foreground' />
                                      <span className='sr-only'>
                                        Edit node name
                                      </span>
                                    </button>
                                  </div>
                                  <p className='font-mono text-xs text-muted-foreground'>
                                    {node.node_id}
                                  </p>
                                </>
                              )}
                            </div>
                          </div>
                        </td>
                        <td className='px-4 py-3 text-muted-foreground'>
                          {node.zone || '-'}
                        </td>
                        <td className='px-4 py-3 text-muted-foreground'>
                          <div className='flex items-center gap-2'>
                            <Cpu className='size-4 text-muted-foreground/80' />
                            <span>
                              {node.cpu_cores > 0
                                ? `${node.cpu_cores} cores`
                                : '-'}
                            </span>
                          </div>
                        </td>
                        <td className='px-4 py-3 text-muted-foreground'>
                          <div className='flex items-center gap-2'>
                            <MemoryStick className='size-4 text-muted-foreground/80' />
                            <span>{formatMegabytes(node.ram_total_mb)}</span>
                          </div>
                        </td>
                        <td className='px-4 py-3 text-muted-foreground'>
                          <div className='flex items-center gap-2'>
                            <HardDrive className='size-4 text-muted-foreground/80' />
                            <span>
                              {node.disk_count > 0
                                ? `${node.disk_count} disk${node.disk_count > 1 ? 's' : ''}`
                                : '-'}
                            </span>
                          </div>
                        </td>
                        <td className='px-4 py-3 text-muted-foreground'>
                          <div className='flex items-center gap-2'>
                            <Server className='size-4 text-muted-foreground/80' />
                            <span>
                              {node.gpu_count > 0
                                ? `${node.gpu_count} GPU${node.gpu_count > 1 ? 's' : ''}`
                                : '-'}
                            </span>
                          </div>
                        </td>
                        <td className='px-4 py-3'>
                          <div className='flex items-center justify-end gap-2'>
                            <button
                              type='button'
                              onClick={() => handleViewNode(node)}
                              className='inline-flex h-9 w-9 items-center justify-center rounded-xl border border-border bg-background text-muted-foreground transition hover:border-primary/40 hover:text-foreground'
                              aria-label={`View ${node.name || node.node_id}`}
                            >
                              <Eye className='size-4' />
                            </button>
                            <DropdownMenu modal={false}>
                              <DropdownMenuTrigger asChild>
                                <button
                                  type='button'
                                  className='inline-flex h-9 w-9 items-center justify-center rounded-xl border border-border bg-background text-muted-foreground transition hover:border-primary/40 hover:text-foreground'
                                  aria-label={`Open actions for ${node.name || node.node_id}`}
                                >
                                  <MoreHorizontal className='size-4' />
                                </button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align='end' className='w-56'>
                                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                <DropdownMenuItem
                                  onClick={() => beginEdit(node)}
                                >
                                  Rename
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuLabel>
                                  Assign to zone
                                </DropdownMenuLabel>
                                {zones.length > 0 ? (
                                  <DropdownMenuRadioGroup
                                    value={node.zone_id || ''}
                                  >
                                    {zones.map((zone) => (
                                      <DropdownMenuRadioItem
                                        key={zone.id}
                                        value={zone.id}
                                        disabled={assignZoneMutation.isPending}
                                        onSelect={(event) => {
                                          event.preventDefault()
                                          handleAssignZone(node, zone)
                                        }}
                                      >
                                        {zone.name}
                                      </DropdownMenuRadioItem>
                                    ))}
                                  </DropdownMenuRadioGroup>
                                ) : (
                                  <DropdownMenuItem disabled>
                                    {zonesQuery.isLoading
                                      ? 'Loading zones...'
                                      : 'No zones available'}
                                  </DropdownMenuItem>
                                )}
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </section>
      </Main>
    </>
  )
}
