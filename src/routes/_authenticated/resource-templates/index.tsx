import { createFileRoute } from '@tanstack/react-router'
import { TemplateRenderListPage } from '@/features/resource-templates'

export const Route = createFileRoute('/_authenticated/resource-templates/')({
  component: TemplateRenderListPage,
})
