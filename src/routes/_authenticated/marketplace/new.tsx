import { createFileRoute } from '@tanstack/react-router'
import { NewMarketplacePage } from '@/features/marketplace/new'

export const Route = createFileRoute('/_authenticated/marketplace/new')({
  component: NewMarketplacePage,
})
