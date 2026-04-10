import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, ShieldAlert } from 'lucide-react'
import { toast } from 'sonner'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { getTemplateRender, updateTemplateRender } from './api'
import { TemplateRenderForm } from './form'
import type { TemplateRenderInput } from './types'

export function EditTemplateRenderPage(props: { templateID: string }) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const detailQuery = useQuery({
    queryKey: ['resource-template-render', props.templateID],
    queryFn: () => getTemplateRender(props.templateID),
  })

  const updateMutation = useMutation({
    mutationFn: (input: TemplateRenderInput) =>
      updateTemplateRender(props.templateID, input),
    onSuccess: (template) => {
      toast.success('Template render updated')
      queryClient.invalidateQueries({ queryKey: ['resource-template-renders'] })
      queryClient.setQueryData(
        ['resource-template-render', template.id],
        template
      )
      navigate({
        to: '/resource-templates/$templateId',
        params: { templateId: template.id },
      })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error
          ? error.message
          : 'Failed to update template render'
      )
    },
  })

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
            <h1 className='page-title'>Edit YAML template</h1>
            <p className='page-copy'>
              Refine the YAML render contract without changing the surrounding
              resource platform workflow.
            </p>
          </div>
          <Button variant='outline' asChild>
            <Link
              to='/resource-templates/$templateId'
              params={{ templateId: props.templateID }}
            >
              <ArrowLeft className='size-4' />
              Back to detail
            </Link>
          </Button>
        </section>

        {detailQuery.isLoading ? (
          <div className='rounded-xl border border-border/80 bg-muted/50 px-4 py-8 text-sm text-muted-foreground'>
            Loading template render...
          </div>
        ) : detailQuery.isError || !detailQuery.data ? (
          <div className='flex items-start gap-3 rounded-xl border border-warning/25 bg-warning-soft px-4 py-4 text-sm text-warning'>
            <ShieldAlert className='mt-0.5 size-4 shrink-0' />
            <span>
              {detailQuery.error instanceof Error
                ? detailQuery.error.message
                : 'Template render not found'}
            </span>
          </div>
        ) : (
          <TemplateRenderForm
            initialValue={detailQuery.data}
            submitLabel='Save template'
            busy={updateMutation.isPending}
            onSubmit={async (value) => {
              await updateMutation.mutateAsync(value)
            }}
          />
        )}
      </Main>
    </>
  )
}
