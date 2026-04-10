import { createFileRoute } from '@tanstack/react-router'
import { EditTemplateRenderPage } from '@/features/resource-templates/edit'

export const Route = createFileRoute(
  '/_authenticated/resource-templates/$templateId/edit'
)({
  component: TemplateRenderEditRoute,
})

function TemplateRenderEditRoute() {
  const { templateId } = Route.useParams()
  return <EditTemplateRenderPage templateID={templateId} />
}
