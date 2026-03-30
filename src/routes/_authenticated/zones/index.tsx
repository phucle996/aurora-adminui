import { createFileRoute } from '@tanstack/react-router'
import { ZonesPage } from '@/features/zones'

export const Route = createFileRoute('/_authenticated/zones/')({
  component: ZonesPage,
})
