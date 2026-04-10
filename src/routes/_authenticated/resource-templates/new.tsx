import { createFileRoute } from '@tanstack/react-router'
import { NewTemplateRenderPage } from '@/features/resource-templates/new'

export const Route = createFileRoute('/_authenticated/resource-templates/new')({
  component: NewTemplateRenderPage,
})
