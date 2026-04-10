import { createFileRoute } from '@tanstack/react-router'
import { MarketplaceListPage } from '@/features/marketplace/index'

export const Route = createFileRoute('/_authenticated/marketplace/')({
  component: MarketplaceRoute,
})

function MarketplaceRoute() {
  return <MarketplaceListPage />
}
