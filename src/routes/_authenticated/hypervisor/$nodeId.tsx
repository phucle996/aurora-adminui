import { createFileRoute } from '@tanstack/react-router'
import { HypervisorNodeDetailPage } from '@/features/hypervisor/detail'

export const Route = createFileRoute('/_authenticated/hypervisor/$nodeId')({
  component: HypervisorNodeDetailRoute,
})

function HypervisorNodeDetailRoute() {
  const { nodeId } = Route.useParams()
  return <HypervisorNodeDetailPage nodeID={nodeId} />
}
