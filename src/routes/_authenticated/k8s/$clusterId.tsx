import { createFileRoute } from '@tanstack/react-router'
import { K8sClusterDetailPage } from '@/features/k8s/detail'

export const Route = createFileRoute('/_authenticated/k8s/$clusterId')({
  component: K8sClusterDetailRoute,
})

function K8sClusterDetailRoute() {
  const { clusterId } = Route.useParams()
  return <K8sClusterDetailPage clusterID={clusterId} />
}
