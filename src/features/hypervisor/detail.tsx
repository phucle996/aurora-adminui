import { useEffect, useRef, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import {
  ArrowLeft,
  CircleDot,
  Cpu,
  HardDrive,
  MemoryStick,
  Monitor,
  Server,
  ShieldAlert,
} from 'lucide-react'
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
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
  buildHypervisorNodeMetricsWebSocketURL,
  getHypervisorNodeDetail,
  type HypervisorMetricSeries,
} from './api'

function formatBytes(bytes: number) {
  if (!bytes || bytes <= 0) return '-'
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  let value = bytes
  let index = 0
  while (value >= 1024 && index < units.length - 1) {
    value /= 1024
    index += 1
  }
  const fractionDigits = value >= 10 || index === 0 ? 0 : 1
  return `${value.toFixed(fractionDigits)} ${units[index]}`
}

function formatMegabytes(megabytes: number) {
  if (!megabytes || megabytes <= 0) return '-'
  if (megabytes >= 1024) {
    const gigabytes = megabytes / 1024
    const fractionDigits = gigabytes >= 10 ? 0 : 1
    return `${gigabytes.toFixed(fractionDigits)} GB`
  }
  return `${Math.round(megabytes)} MB`
}

function formatGigabytes(gigabytes: number) {
  if (!gigabytes || gigabytes <= 0) return '-'
  return `${Math.round(gigabytes)} GB`
}

function formatRate(bytesPerSecond: number) {
  if (!bytesPerSecond || bytesPerSecond <= 0) return '0 B/s'
  return `${formatBytes(bytesPerSecond)}/s`
}

function formatPercent(value: number) {
  return `${value.toFixed(1)}%`
}

function formatMetricValue(value: number, formatter?: (value: number) => string) {
  if (formatter) {
    return formatter(value)
  }
  return String(Number(value.toFixed(2)))
}

function compactTimeLabel(timestamp: string) {
  const date = new Date(timestamp)
  if (Number.isNaN(date.getTime())) {
    return '--:--'
  }
  return date.toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

function statusTone(status: string) {
  switch (status.toLowerCase()) {
    case 'online':
      return 'bg-emerald-500'
    case 'offline':
      return 'bg-rose-500'
    default:
      return 'bg-amber-500'
  }
}

function InfoRow(props: { label: string; value: string }) {
  return (
    <div className='grid gap-1 rounded-xl border border-border/70 bg-muted/20 px-4 py-3'>
      <span className='text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground'>
        {props.label}
      </span>
      <span className='text-sm font-medium text-foreground'>{props.value}</span>
    </div>
  )
}

function InventoryRow(props: { title: string; meta: string; trailing: string }) {
  return (
    <div className='flex items-center justify-between gap-4 rounded-xl border border-border/70 bg-muted/20 px-4 py-3'>
      <div className='min-w-0'>
        <p className='truncate text-sm font-medium text-foreground'>{props.title}</p>
        <p className='truncate text-xs text-muted-foreground'>{props.meta}</p>
      </div>
      <span className='shrink-0 text-xs font-medium text-muted-foreground'>
        {props.trailing}
      </span>
    </div>
  )
}

function MetricChartCard(props: {
  title: string
  description: string
  value: string
  color: string
  points: Array<{ timestamp: string; value: number }>
  formatAxis?: (value: number) => string
}) {
  const chartData = props.points.map((point) => ({
    label: compactTimeLabel(point.timestamp),
    value: Number(point.value.toFixed(2)),
    tooltipLabel: new Date(point.timestamp).toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    }),
  }))
  const gradientID = `fill-${props.title.replace(/\s+/g, '-').toLowerCase()}`

  return (
    <Card>
      <CardHeader>
        <div className='flex flex-wrap items-center justify-between gap-3'>
          <div className='space-y-1'>
            <CardTitle>{props.title}</CardTitle>
            <CardDescription>{props.description}</CardDescription>
          </div>
          <p className='text-xl font-semibold text-foreground'>{props.value}</p>
        </div>
      </CardHeader>
      <CardContent>
        {chartData.length === 0 ? (
          <div className='rounded-xl border border-dashed border-border/80 bg-muted/20 px-4 py-10 text-sm text-muted-foreground'>
            Waiting for live telemetry stream...
          </div>
        ) : (
          <div
            className='rounded-xl border border-border/70 px-2 py-3'
            style={{
              backgroundColor: 'hsl(var(--background) / 0.72)',
              backgroundImage: [
                'linear-gradient(to right, hsl(var(--border) / 0.18) 1px, transparent 1px)',
                'linear-gradient(to bottom, hsl(var(--border) / 0.18) 1px, transparent 1px)',
              ].join(','),
              backgroundSize: '32px 32px',
            }}
          >
            <ResponsiveContainer width='100%' height={260}>
              <AreaChart data={chartData}>
                <defs>
                  <linearGradient id={gradientID} x1='0' y1='0' x2='0' y2='1'>
                    <stop offset='5%' stopColor={props.color} stopOpacity={0.35} />
                    <stop offset='95%' stopColor={props.color} stopOpacity={0.02} />
                  </linearGradient>
                </defs>
                <CartesianGrid
                  stroke='hsl(var(--border) / 0.24)'
                  strokeDasharray='3 3'
                />
                <XAxis
                  dataKey='label'
                  tickLine={false}
                  axisLine={false}
                  tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
                />
                <YAxis
                  tickLine={false}
                  axisLine={false}
                  width={56}
                  tickFormatter={props.formatAxis}
                  tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
                />
                <Tooltip
                  labelFormatter={() => ''}
                  separator=''
                  formatter={(value) => [
                    formatMetricValue(Number(value), props.formatAxis),
                    '',
                  ]}
                  cursor={{ stroke: props.color, strokeOpacity: 0.25 }}
                  contentStyle={{
                    borderRadius: 16,
                    border: '1px solid hsl(var(--border))',
                    background: 'hsl(var(--card))',
                  }}
                />
                <Area
                  type='monotone'
                  dataKey='value'
                  stroke={props.color}
                  fill={`url(#${gradientID})`}
                  strokeWidth={2}
                  isAnimationActive={false}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

function latestMetricValue(
  series: Record<string, HypervisorMetricSeries>,
  name: string
) {
  return series[name]?.latest ?? 0
}

type HypervisorMetricsStreamPayload = {
  id: string
  generated_at: string
  window_sec: number
  step_sec: number
  series: HypervisorMetricSeries[]
}

export function HypervisorNodeDetailPage(props: { nodeID: string }) {
  const nodeQuery = useQuery({
    queryKey: ['hypervisor-node-detail', props.nodeID],
    queryFn: () => getHypervisorNodeDetail(props.nodeID),
  })
  const [liveMetrics, setLiveMetrics] = useState<HypervisorMetricSeries[]>([])
  const [chartMetrics, setChartMetrics] = useState<HypervisorMetricSeries[]>([])
  const lastChartCommitAtRef = useRef(0)

  useEffect(() => {
    if (!props.nodeID) {
      setLiveMetrics([])
      setChartMetrics([])
      lastChartCommitAtRef.current = 0
      return
    }

    let socket: WebSocket | null = null
    let reconnectTimer: number | null = null
    let cancelled = false

    const connect = () => {
      if (cancelled) {
        return
      }

      socket = new WebSocket(buildHypervisorNodeMetricsWebSocketURL(props.nodeID))
      socket.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data) as
            | HypervisorMetricsStreamPayload
            | { type?: string; message?: string }
          if ('series' in payload && Array.isArray(payload.series)) {
            setLiveMetrics(payload.series)
            const generatedAt = new Date(payload.generated_at).getTime()
            const now = Number.isNaN(generatedAt) ? Date.now() : generatedAt
            if (
              lastChartCommitAtRef.current === 0 ||
              now-lastChartCommitAtRef.current >= 15_000
            ) {
              lastChartCommitAtRef.current = now
              setChartMetrics(payload.series)
            }
          }
        } catch {
          // Ignore malformed frames and keep the last good snapshot.
        }
      }
      socket.onclose = () => {
        if (!cancelled) {
          reconnectTimer = window.setTimeout(connect, 3000)
        }
      }
      socket.onerror = () => {
        socket?.close()
      }
    }

    connect()

    return () => {
      cancelled = true
      if (reconnectTimer !== null) {
        window.clearTimeout(reconnectTimer)
      }
      socket?.close()
    }
  }, [props.nodeID])

  const node = nodeQuery.data
  const chartSeries = chartMetrics.reduce<Record<string, HypervisorMetricSeries>>(
    (acc, item) => {
      if (!item.label && !acc[item.name]) {
        acc[item.name] = item
      }
      return acc
    },
    {}
  )
  const liveSeries = liveMetrics.reduce<Record<string, HypervisorMetricSeries>>(
    (acc, item) => {
      if (!item.label && !acc[item.name]) {
        acc[item.name] = item
      }
      return acc
    },
    {}
  )
  const perGPUSeries = chartMetrics.filter(
    (item) => item.name === 'gpu_usage_percent_device'
  )
  const cpuUsage = latestMetricValue(liveSeries, 'cpu_usage_percent')
  const ramUsedBytes = latestMetricValue(liveSeries, 'ram_used_bytes')
  const diskRate = latestMetricValue(liveSeries, 'disk_io_bytes_per_sec')
  const networkRate = latestMetricValue(liveSeries, 'network_bytes_per_sec')

  return (
    <>
      <Header fixed>
        <div className='min-w-0'>
          <p className='subtle-kicker'>Infrastructure</p>
          <h1 className='truncate text-lg font-semibold text-foreground'>
            Node Detail
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
            <p className='subtle-kicker'>Hypervisor node</p>
            <h1 className='page-title'>{node?.name || props.nodeID}</h1>
            <p className='page-copy'>
              Review identity, live streamed telemetry, and the full hardware
              profile for this KVM node.
            </p>
          </div>
          <Button variant='outline' asChild>
            <Link to='/hypervisor'>
              <ArrowLeft className='size-4' />
              Back to nodes
            </Link>
          </Button>
        </section>

        {nodeQuery.isLoading ? (
          <div className='rounded-xl border border-border/80 bg-muted/50 px-4 py-8 text-sm text-muted-foreground'>
            Loading node detail...
          </div>
        ) : nodeQuery.isError ? (
          <div className='flex items-start gap-3 rounded-xl border border-warning/25 bg-warning-soft px-4 py-4 text-sm text-warning'>
            <ShieldAlert className='mt-0.5 size-4 shrink-0' />
            <span>
              {nodeQuery.error instanceof Error
                ? nodeQuery.error.message
                : 'Failed to load node detail'}
            </span>
          </div>
        ) : node ? (
          <>
            <Card>
              <CardHeader>
                <div className='flex flex-wrap items-center justify-between gap-3'>
                  <div className='flex items-center gap-3'>
                    <span
                      className={`size-3 rounded-full ${statusTone(node.status)}`}
                    />
                    <div>
                      <CardTitle>{node.name}</CardTitle>
                      <CardDescription>{node.node_id}</CardDescription>
                    </div>
                  </div>
                  <Badge
                    variant='secondary'
                    className='gap-2 rounded-full px-3 py-1'
                  >
                    <CircleDot className='size-3.5' />
                    {node.status}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className='grid gap-4 md:grid-cols-2 xl:grid-cols-4'>
                <InfoRow label='Node ID' value={node.node_id} />
                <InfoRow label='Hostname' value={node.hostname || '-'} />
                <InfoRow label='Zone' value={node.zone || 'Unassigned'} />
                <InfoRow label='Status' value={node.status} />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <div className='flex items-center gap-3'>
                  <span className='flex size-11 items-center justify-center rounded-2xl bg-accent text-accent-foreground'>
                    <Server className='size-5' />
                  </span>
                  <div>
                    <p className='subtle-kicker'>Hardware profile</p>
                    <CardTitle>Compute and inventory</CardTitle>
                    <CardDescription>
                      Static hardware shape plus detailed disk and GPU model
                      inventory.
                    </CardDescription>
                  </div>
                </div>
              </CardHeader>
              <CardContent className='space-y-6'>
                <div className='grid gap-4 md:grid-cols-2 xl:grid-cols-4'>
                  <div className='rounded-2xl border border-border/70 bg-muted/20 p-4'>
                    <div className='mb-3 flex items-center gap-2 text-sm font-medium text-foreground'>
                      <Cpu className='size-4 text-muted-foreground' />
                      CPU
                    </div>
                    <p className='text-2xl font-semibold text-foreground'>
                      {node.cpu_cores} cores
                    </p>
                    <p className='mt-1 text-sm text-muted-foreground'>
                      {node.cpu_model || 'CPU model unavailable'}
                    </p>
                  </div>

                  <div className='rounded-2xl border border-border/70 bg-muted/20 p-4'>
                    <div className='mb-3 flex items-center gap-2 text-sm font-medium text-foreground'>
                      <MemoryStick className='size-4 text-muted-foreground' />
                      Memory
                    </div>
                    <p className='text-2xl font-semibold text-foreground'>
                      {formatMegabytes(node.ram_total_mb)}
                    </p>
                    <p className='mt-1 text-sm text-muted-foreground'>
                      Total system RAM
                    </p>
                  </div>

                  <div className='rounded-2xl border border-border/70 bg-muted/20 p-4'>
                    <div className='mb-3 flex items-center gap-2 text-sm font-medium text-foreground'>
                      <HardDrive className='size-4 text-muted-foreground' />
                      Disk devices
                    </div>
                    <p className='text-2xl font-semibold text-foreground'>
                      {node.disk_count}
                    </p>
                    <p className='mt-1 text-sm text-muted-foreground'>
                      Attached storage devices
                    </p>
                  </div>

                  <div className='rounded-2xl border border-border/70 bg-muted/20 p-4'>
                    <div className='mb-3 flex items-center gap-2 text-sm font-medium text-foreground'>
                      <Monitor className='size-4 text-muted-foreground' />
                      GPUs
                    </div>
                    <p className='text-2xl font-semibold text-foreground'>
                      {node.gpu_count}
                    </p>
                    <p className='mt-1 text-sm text-muted-foreground'>
                      Visible accelerators
                    </p>
                  </div>
                </div>

                <div className='grid gap-6 xl:grid-cols-2'>
                  <div className='space-y-3'>
                    <div className='flex items-center justify-between gap-3'>
                      <div>
                        <p className='text-sm font-semibold text-foreground'>
                          Disk inventory
                        </p>
                        <p className='text-xs text-muted-foreground'>
                          One row per disk model and block device.
                        </p>
                      </div>
                      <Badge variant='secondary'>{node.disk_count} disks</Badge>
                    </div>
                    {node.disks.length === 0 ? (
                      <div className='rounded-xl border border-dashed border-border/80 bg-muted/20 px-4 py-10 text-sm text-muted-foreground'>
                        No disk inventory reported for this node.
                      </div>
                    ) : (
                      node.disks.map((disk, index) => (
                        <InventoryRow
                          key={`${disk.name}-${index}`}
                          title={disk.model || disk.name || 'Unknown disk'}
                          meta={disk.name || 'Unnamed block device'}
                          trailing={formatGigabytes(disk.size_gb)}
                        />
                      ))
                    )}
                  </div>

                  <div className='space-y-3'>
                    <div className='flex items-center justify-between gap-3'>
                      <div>
                        <p className='text-sm font-semibold text-foreground'>
                          GPU inventory
                        </p>
                        <p className='text-xs text-muted-foreground'>
                          One row per GPU model with vendor and PCI reference.
                        </p>
                      </div>
                      <Badge variant='secondary'>{node.gpu_count} gpus</Badge>
                    </div>
                    {node.gpus.length === 0 ? (
                      <div className='rounded-xl border border-dashed border-border/80 bg-muted/20 px-4 py-10 text-sm text-muted-foreground'>
                        No GPU inventory reported for this node.
                      </div>
                    ) : (
                      node.gpus.map((gpu, index) => (
                        <InventoryRow
                          key={`${gpu.pci_address}-${index}`}
                          title={gpu.model || 'Unknown GPU'}
                          meta={[
                            gpu.vendor || 'unknown vendor',
                            gpu.pci_address || 'no pci address',
                            gpu.driver_version || 'driver unknown',
                          ].join(' · ')}
                          trailing={
                            gpu.memory_total_mb > 0
                              ? formatMegabytes(gpu.memory_total_mb)
                              : 'shared'
                          }
                        />
                      ))
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>

            <section className='grid gap-4 xl:grid-cols-2'>
              <MetricChartCard
                title='CPU usage trend'
                description='Live streamed host CPU usage over the last 15 minutes.'
                value={formatPercent(cpuUsage)}
                color='#10b981'
                points={chartSeries.cpu_usage_percent?.points || []}
                formatAxis={(value) => `${value}%`}
              />
              <MetricChartCard
                title='Memory usage trend'
                description='Live host RAM usage sampled from the telemetry stream.'
                value={formatBytes(ramUsedBytes)}
                color='#3b82f6'
                points={chartSeries.ram_used_bytes?.points || []}
                formatAxis={(value) => formatBytes(value)}
              />
              <MetricChartCard
                title='Disk throughput trend'
                description='Aggregated disk read/write throughput from the node.'
                value={formatRate(diskRate)}
                color='#f59e0b'
                points={chartSeries.disk_io_bytes_per_sec?.points || []}
                formatAxis={(value) => formatBytes(value)}
              />
              <MetricChartCard
                title='Network throughput trend'
                description='Aggregated network RX/TX throughput across interfaces.'
                value={formatRate(networkRate)}
                color='#0ea5e9'
                points={chartSeries.network_bytes_per_sec?.points || []}
                formatAxis={(value) => formatBytes(value)}
              />
            </section>

            {perGPUSeries.length > 0 ? (
              <section className='grid gap-4 xl:grid-cols-2'>
                {perGPUSeries.map((series) => (
                  <MetricChartCard
                    key={`${series.name}-${series.label}`}
                    title={`GPU · ${series.label}`}
                    description='Per-GPU utilization streamed from Victoria labels.'
                    value={formatPercent(series.latest)}
                    color='#8b5cf6'
                    points={series.points || []}
                    formatAxis={(value) => `${value}%`}
                  />
                ))}
              </section>
            ) : null}
          </>
        ) : null}
      </Main>
    </>
  )
}
