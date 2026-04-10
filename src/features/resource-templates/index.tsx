import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { Eye, FilePlus2, PencilLine, ShieldAlert, Trash2 } from 'lucide-react'
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
import { deleteTemplateRender, listTemplateRenderCatalog } from './api'
import { formatTemplateDateTime } from './utils'

export function TemplateRenderListPage() {
  const queryClient = useQueryClient()
  const templatesQuery = useQuery({
    queryKey: ['resource-template-renders'],
    queryFn: listTemplateRenderCatalog,
  })
  const deleteMutation = useMutation({
    mutationFn: deleteTemplateRender,
    onSuccess: () => {
      toast.success('Template render deleted')
      queryClient.invalidateQueries({ queryKey: ['resource-template-renders'] })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to delete template render'
      )
    },
  })

  const templates = templatesQuery.data || []

  async function handleDelete(templateID: string, templateName: string) {
    if (!window.confirm(`Delete template "${templateName}"?`)) {
      return
    }
    await deleteMutation.mutateAsync(templateID)
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
            <p className='subtle-kicker'>Resource platform</p>
            <h1 className='page-title'>Template Render</h1>
            <p className='page-copy'>
              Manage YAML render templates that dataplane workers consume after
              create database jobs are enqueued into Redis Streams.
            </p>
          </div>
          <Button asChild>
            <Link to='/resource-templates/new'>
              <FilePlus2 className='size-4' />
              New template
            </Link>
          </Button>
        </section>

        <div className='overflow-hidden rounded-md border bg-card'>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className='w-[260px]'>Name</TableHead>
                <TableHead className='w-[150px]'>Resource</TableHead>
                <TableHead className='w-[220px]'>Stream</TableHead>
                <TableHead className='w-[160px]'>YAML status</TableHead>
                <TableHead className='w-[180px]'>Updated</TableHead>
                <TableHead className='w-[160px] text-right'>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {templatesQuery.isLoading ? (
                <TableRow>
                  <TableCell colSpan={6} className='h-24 text-center'>
                    Loading template renders...
                  </TableCell>
                </TableRow>
              ) : templatesQuery.isError ? (
                <TableRow>
                  <TableCell colSpan={6} className='h-24'>
                    <div className='flex items-start justify-center gap-3 text-sm text-warning'>
                      <ShieldAlert className='mt-0.5 size-4 shrink-0' />
                      <span>
                        {templatesQuery.error instanceof Error
                          ? templatesQuery.error.message
                          : 'Failed to load template renders'}
                      </span>
                    </div>
                  </TableCell>
                </TableRow>
              ) : templates.length > 0 ? (
                templates.map((template) => {
                  return (
                    <TableRow key={template.id}>
                      <TableCell>
                        <div className='space-y-1'>
                          <p className='font-medium text-foreground'>
                            {template.name}
                          </p>
                          <p className='text-sm leading-6 text-muted-foreground'>
                            {template.description || '-'}
                          </p>
                        </div>
                      </TableCell>
                      <TableCell className='text-sm text-muted-foreground'>
                        <div>{template.resource_type}</div>
                        <div className='font-mono text-xs'>{template.resource_model}</div>
                      </TableCell>
                      <TableCell className='font-mono text-xs text-muted-foreground'>
                        <div>{template.stream_key}</div>
                        <div>{template.consumer_group}</div>
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant='outline'
                          className={
                            template.yaml_valid
                              ? 'border-emerald-200 bg-emerald-50 text-emerald-700'
                              : 'border-rose-200 bg-rose-50 text-rose-700'
                          }
                        >
                          {template.yaml_valid ? 'Valid' : 'Invalid'}
                        </Badge>
                      </TableCell>
                      <TableCell className='text-sm text-muted-foreground'>
                        {formatTemplateDateTime(template.updated_at)}
                      </TableCell>
                      <TableCell>
                        <div className='flex items-center justify-end gap-2'>
                          <Button variant='outline' size='icon' asChild>
                            <Link
                              to='/resource-templates/$templateId'
                              params={{ templateId: template.id }}
                            >
                              <Eye className='size-4' />
                            </Link>
                          </Button>
                          <Button variant='outline' size='icon' asChild>
                            <Link
                              to='/resource-templates/$templateId/edit'
                              params={{ templateId: template.id }}
                            >
                              <PencilLine className='size-4' />
                            </Link>
                          </Button>
                          <Button
                            variant='outline'
                            size='icon'
                            disabled={deleteMutation.isPending}
                            onClick={() => void handleDelete(template.id, template.name)}
                          >
                            <Trash2 className='size-4' />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  )
                })
              ) : (
                <TableRow>
                  <TableCell colSpan={6} className='h-24 text-center text-muted-foreground'>
                    No template renders have been created yet.
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
