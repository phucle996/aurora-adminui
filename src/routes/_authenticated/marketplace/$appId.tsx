import { createFileRoute } from '@tanstack/react-router'
import { MarketplaceDetailPage } from '@/features/marketplace/detail'

export const Route = createFileRoute('/_authenticated/marketplace/$appId')({
  component: MarketplaceDetailRoute,
})

function MarketplaceDetailRoute() {
  const { appId } = Route.useParams()
  return <MarketplaceDetailPage appID={appId} />
}
