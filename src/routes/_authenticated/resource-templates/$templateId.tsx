import { Outlet, createFileRoute, useRouterState } from '@tanstack/react-router'
import { TemplateRenderDetailPage } from '@/features/resource-templates/detail'

export const Route = createFileRoute(
  '/_authenticated/resource-templates/$templateId'
)({
  component: TemplateRenderDetailRoute,
})

function TemplateRenderDetailRoute() {
  const { templateId } = Route.useParams()
  const pathname = useRouterState({ select: (state) => state.location.pathname })
  if (pathname.endsWith('/edit')) {
    return <Outlet />
  }
  return <TemplateRenderDetailPage templateID={templateId} />
}
