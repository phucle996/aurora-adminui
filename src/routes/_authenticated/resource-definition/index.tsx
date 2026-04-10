import { createFileRoute } from '@tanstack/react-router'
import { ResourceTypesPage } from '@/features/resource-types'

export const Route = createFileRoute('/_authenticated/resource-definition/')({
  component: ResourceDefinitionRoute,
})

function ResourceDefinitionRoute() {
  return <ResourceTypesPage />
}
