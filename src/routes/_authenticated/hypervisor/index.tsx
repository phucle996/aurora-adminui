import { createFileRoute } from '@tanstack/react-router'
import { HypervisorNodes } from '@/features/hypervisor'

export const Route = createFileRoute('/_authenticated/hypervisor/')({
  component: HypervisorNodes,
})
