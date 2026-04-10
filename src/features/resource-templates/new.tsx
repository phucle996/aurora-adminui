import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import { ArrowLeft } from 'lucide-react'
import { toast } from 'sonner'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { createTemplateRender } from './api'
import { TemplateRenderForm } from './form'
import type { TemplateRenderInput } from './types'

export function NewTemplateRenderPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const createMutation = useMutation({
    mutationFn: (input: TemplateRenderInput) => createTemplateRender(input),
    onSuccess: (template) => {
      toast.success('Template render created')
      queryClient.invalidateQueries({ queryKey: ['resource-template-renders'] })
      navigate({
        to: '/resource-templates/$templateId',
        params: { templateId: template.id },
      })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error
          ? error.message
          : 'Failed to create template render'
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
            <h1 className='page-title'>New YAML template</h1>
            <p className='page-copy'>
              Define the Redis payload contract and YAML manifest that the
              dataplane worker renders at execution time.
            </p>
          </div>
          <Button variant='outline' asChild>
            <Link to='/resource-templates'>
              <ArrowLeft className='size-4' />
              Back to list
            </Link>
          </Button>
        </section>

        <TemplateRenderForm
          submitLabel='Create template'
          busy={createMutation.isPending}
          onSubmit={async (value) => {
            await createMutation.mutateAsync(value)
          }}
        />
      </Main>
    </>
  )
}
