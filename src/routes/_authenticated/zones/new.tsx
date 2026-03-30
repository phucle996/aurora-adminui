import { createFileRoute } from '@tanstack/react-router'
import { AddZonePage } from '@/features/zones/new'

export const Route = createFileRoute('/_authenticated/zones/new')({
  component: AddZonePage,
})
