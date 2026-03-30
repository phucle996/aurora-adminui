import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, CloudUpload, FileJson2, ShipWheel } from 'lucide-react'
import { toast } from 'sonner'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { listZones } from '@/features/zones/api'
import { createK8sCluster } from './api'

export function AddK8sClusterPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const zonesQuery = useQuery({
    queryKey: ['zones'],
    queryFn: listZones,
  })

  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [zoneId, setZoneId] = useState('')
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [supportsDbaas, setSupportsDbaas] = useState(true)
  const [supportsServerless, setSupportsServerless] = useState(false)
  const [supportsGenericWorkloads, setSupportsGenericWorkloads] = useState(true)

  const createClusterMutation = useMutation({
    mutationFn: createK8sCluster,
    onSuccess: (cluster) => {
      toast.success(`Cluster "${cluster.name}" imported`)
      queryClient.invalidateQueries({ queryKey: ['k8s-clusters'] })
      navigate({ to: '/k8s/$clusterId', params: { clusterId: cluster.id } })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to import cluster'
      )
    },
  })

  const selectedFileMeta = useMemo(() => {
    if (!selectedFile) return 'No kubeconfig selected yet'
    return `${selectedFile.name} · ${Math.ceil(selectedFile.size / 1024)} KB`
  }, [selectedFile])

  function handleSubmit() {
    if (!name.trim()) {
      toast.error('Cluster name is required')
      return
    }
    if (!selectedFile) {
      toast.error('Please attach a kubeconfig file')
      return
    }

    createClusterMutation.mutate({
      name: name.trim(),
      description: description.trim(),
      zoneId,
      supportsDbaas,
      supportsServerless,
      supportsGenericWorkloads,
      kubeconfig: selectedFile,
    })
  }

  return (
    <>
      <Header fixed>
        <div className='min-w-0'>
          <p className='subtle-kicker'>Kubernetes substrate</p>
          <h1 className='truncate text-lg font-semibold text-foreground'>
            Add K8s Cluster
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
            <p className='subtle-kicker'>Cluster import</p>
            <h1 className='page-title'>Register a Kubernetes cluster</h1>
            <p className='page-copy'>
              Import a kubeconfig-backed cluster, bind it to an optional zone,
              and declare which higher-level platform capabilities it can host.
            </p>
          </div>
          <Button variant='outline' asChild>
            <Link to='/k8s'>
              <ArrowLeft className='size-4' />
              Back to clusters
            </Link>
          </Button>
        </section>

        <div className='grid gap-6 xl:grid-cols-[minmax(0,420px)_minmax(0,1fr)]'>
          <Card>
            <CardHeader>
              <CardTitle>Cluster definition</CardTitle>
            </CardHeader>
            <CardContent className='space-y-5'>
              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Cluster name
                </label>
                <Input
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder='Example: aurora-platform-prod'
                />
              </div>

              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Zone
                </label>
                <Select
                  value={zoneId || '__none__'}
                  onValueChange={(value) =>
                    setZoneId(value === '__none__' ? '' : value)
                  }
                >
                  <SelectTrigger className='w-full'>
                    <SelectValue placeholder='No zone binding' />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='__none__'>No zone binding</SelectItem>
                    {(zonesQuery.data || []).map((zone) => (
                      <SelectItem key={zone.id} value={zone.id}>
                        {zone.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Description
                </label>
                <Textarea
                  className='min-h-28'
                  value={description}
                  onChange={(event) => setDescription(event.target.value)}
                  placeholder='Describe what this cluster should host in the platform.'
                />
              </div>

              <div className='space-y-3 rounded-xl border border-border/80 bg-muted/30 px-4 py-4'>
                <p className='text-sm font-medium text-foreground'>
                  Capability flags
                </p>
                <label className='flex items-start gap-3'>
                  <Checkbox
                    checked={supportsDbaas}
                    onCheckedChange={(checked) =>
                      setSupportsDbaas(checked === true)
                    }
                  />
                  <span className='space-y-1'>
                    <span className='block text-sm font-medium text-foreground'>
                      DBaaS
                    </span>
                    <span className='block text-sm text-muted-foreground'>
                      Allow this cluster to host future database control planes.
                    </span>
                  </span>
                </label>
                <label className='flex items-start gap-3'>
                  <Checkbox
                    checked={supportsServerless}
                    onCheckedChange={(checked) =>
                      setSupportsServerless(checked === true)
                    }
                  />
                  <span className='space-y-1'>
                    <span className='block text-sm font-medium text-foreground'>
                      Serverless
                    </span>
                    <span className='block text-sm text-muted-foreground'>
                      Mark this cluster as eligible for future serverless runtimes.
                    </span>
                  </span>
                </label>
                <label className='flex items-start gap-3'>
                  <Checkbox
                    checked={supportsGenericWorkloads}
                    onCheckedChange={(checked) =>
                      setSupportsGenericWorkloads(checked === true)
                    }
                  />
                  <span className='space-y-1'>
                    <span className='block text-sm font-medium text-foreground'>
                      Generic workloads
                    </span>
                    <span className='block text-sm text-muted-foreground'>
                      Keep this cluster available as a reusable substrate for
                      upcoming resource products.
                    </span>
                  </span>
                </label>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className='space-y-3'>
              <div className='flex items-center gap-3'>
                <span className='flex size-11 items-center justify-center rounded-2xl bg-accent text-accent-foreground'>
                  <ShipWheel className='size-5' />
                </span>
                <div>
                  <p className='subtle-kicker'>Kubeconfig import</p>
                  <CardTitle>Attach cluster access config</CardTitle>
                </div>
              </div>
            </CardHeader>
            <CardContent className='space-y-5'>
              <label className='flex cursor-pointer flex-col items-center justify-center gap-3 rounded-2xl border border-dashed border-border/80 bg-muted/20 px-6 py-10 text-center transition-colors hover:border-primary/40 hover:bg-muted/40'>
                <span className='flex size-12 items-center justify-center rounded-2xl bg-accent text-accent-foreground'>
                  <CloudUpload className='size-5' />
                </span>
                <div className='space-y-1'>
                  <p className='text-sm font-medium text-foreground'>
                    Upload kubeconfig
                  </p>
                  <p className='text-sm text-muted-foreground'>
                    The file will be encrypted in Postgres, then used to validate
                    connectivity and seed the cluster catalog.
                  </p>
                </div>
                <input
                  type='file'
                  accept=".yaml,.yml,.json,.conf,.config"
                  className='hidden'
                  onChange={(event) =>
                    setSelectedFile(event.target.files?.[0] || null)
                  }
                />
              </label>

              <div className='rounded-xl border border-border/80 bg-card px-4 py-4'>
                <div className='flex items-start justify-between gap-3'>
                  <div className='min-w-0 space-y-1'>
                    <div className='flex items-center gap-2'>
                      <FileJson2 className='size-4 text-muted-foreground' />
                      <p className='text-sm font-medium text-foreground'>
                        Selected file
                      </p>
                    </div>
                    <p className='text-sm text-muted-foreground'>
                      {selectedFileMeta}
                    </p>
                  </div>
                  <Badge variant='secondary' className='rounded-full px-3 py-1'>
                    kubeconfig
                  </Badge>
                </div>
              </div>

              <div className='rounded-xl border border-border/80 bg-muted/40 px-4 py-4'>
                <p className='text-sm font-medium text-foreground'>
                  Validation flow
                </p>
                <p className='mt-2 text-sm leading-7 text-muted-foreground'>
                  The platform will parse the current context, resolve the API
                  server, call the cluster for version discovery, and run a light
                  namespace read before saving the final validation status.
                </p>
              </div>

              <div className='flex flex-wrap gap-3'>
                <Button
                  onClick={handleSubmit}
                  disabled={createClusterMutation.isPending}
                >
                  <CloudUpload className='size-4' />
                  {createClusterMutation.isPending
                    ? 'Importing...'
                    : 'Import cluster'}
                </Button>
                <Button variant='outline' asChild>
                  <Link to='/k8s'>Cancel</Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </Main>
    </>
  )
}
