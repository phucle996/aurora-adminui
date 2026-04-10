import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import {
  ArrowLeft,
  Cpu,
  FileText,
  PencilLine,
  ShieldAlert,
  Trash2,
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
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { deleteTemplateRender, getTemplateRender } from './api'
import { formatTemplateDateTime } from './utils'

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

export function TemplateRenderDetailPage(props: { templateID: string }) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const detailQuery = useQuery({
    queryKey: ['resource-template-render', props.templateID],
    queryFn: () => getTemplateRender(props.templateID),
  })
  const deleteMutation = useMutation({
    mutationFn: deleteTemplateRender,
    onSuccess: async () => {
      toast.success('Template render deleted')
      await queryClient.invalidateQueries({ queryKey: ['resource-template-renders'] })
      navigate({ to: '/resource-templates' })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to delete template render'
      )
    },
  })

  const template = detailQuery.data

  async function handleDelete() {
    if (!template) {
      return
    }
    if (!window.confirm(`Delete template "${template.name}"?`)) {
      return
    }
    await deleteMutation.mutateAsync(template.id)
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
            <p className='subtle-kicker'>Template render</p>
            <h1 className='page-title'>
              {template?.name || 'Template render detail'}
            </h1>
            <p className='page-copy'>
              Review the Redis payload contract and the YAML that dataplane
              workers render from resource jobs.
            </p>
          </div>
          <div className='flex items-center gap-3'>
            <Button variant='outline' asChild>
              <Link to='/resource-templates'>
                <ArrowLeft className='size-4' />
                Back
              </Link>
            </Button>
            {template ? (
              <>
                <Button asChild>
                  <Link
                    to='/resource-templates/$templateId/edit'
                    params={{ templateId: template.id }}
                  >
                    <PencilLine className='size-4' />
                    Edit template
                  </Link>
                </Button>
                <Button
                  variant='destructive'
                  onClick={() => void handleDelete()}
                  disabled={deleteMutation.isPending}
                >
                  <Trash2 className='size-4' />
                  Delete
                </Button>
              </>
            ) : null}
          </div>
        </section>

        {detailQuery.isLoading ? (
          <div className='rounded-xl border border-border/80 bg-muted/50 px-4 py-8 text-sm text-muted-foreground'>
            Loading template render...
          </div>
        ) : detailQuery.isError || !template ? (
          <div className='flex items-start gap-3 rounded-xl border border-warning/25 bg-warning-soft px-4 py-4 text-sm text-warning'>
            <ShieldAlert className='mt-0.5 size-4 shrink-0' />
            <span>
              {detailQuery.error instanceof Error
                ? detailQuery.error.message
                : 'Template render not found'}
            </span>
          </div>
        ) : (
          <div className='grid gap-6 xl:grid-cols-[minmax(0,1.7fr)_360px]'>
            <div className='space-y-6'>
              <Card className='rounded-2xl border-border/80'>
                <CardHeader>
                  <CardTitle>Template identity</CardTitle>
                  <CardDescription>
                    Core metadata used by the resource platform catalog.
                  </CardDescription>
                </CardHeader>
                <CardContent className='grid gap-4 md:grid-cols-2'>
                  <InfoRow label='Resource type' value={template.resource_type} />
                  <InfoRow label='Resource model' value={template.resource_model} />
                  <InfoRow label='Stream key' value={template.stream_key} />
                  <InfoRow label='Consumer group' value={template.consumer_group} />
                  <InfoRow
                    label='Created'
                    value={formatTemplateDateTime(template.created_at)}
                  />
                  <InfoRow
                    label='Updated'
                    value={formatTemplateDateTime(template.updated_at)}
                  />
                </CardContent>
              </Card>

              <Card className='rounded-2xl border-border/80'>
                <CardHeader>
                  <CardTitle className='flex items-center gap-2'>
                    <FileText className='size-4 text-muted-foreground' />
                    YAML template
                  </CardTitle>
                  <CardDescription>
                    Manifest body that dataplane workers render from the queued
                    resource job payload.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <pre className='overflow-x-auto rounded-2xl border border-border/70 bg-muted/20 p-4 text-xs leading-6 text-foreground'>
                    {template.yaml_template}
                  </pre>
                </CardContent>
              </Card>
            </div>

            <aside className='space-y-6 xl:sticky xl:top-24 xl:self-start'>
              <Card className='rounded-2xl border-border/80'>
                <CardHeader>
                  <CardTitle className='flex items-center gap-2'>
                    <Cpu className='size-4 text-muted-foreground' />
                    Render pipeline
                  </CardTitle>
                  <CardDescription>
                    Runtime stream settings used before dataplane renders and
                    applies the YAML manifest.
                  </CardDescription>
                </CardHeader>
                <CardContent className='grid gap-3'>
                  <InfoRow label='Stream key' value={template.stream_key} />
                  <InfoRow
                    label='Consumer group'
                    value={template.consumer_group}
                  />
                </CardContent>
              </Card>
            </aside>
          </div>
        )}
      </Main>
    </>
  )
}
