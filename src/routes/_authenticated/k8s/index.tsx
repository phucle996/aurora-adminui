import { createFileRoute } from '@tanstack/react-router'
import { K8sPlatformPage } from '@/features/k8s'

export const Route = createFileRoute('/_authenticated/k8s/')({
  component: K8sPlatformPage,
})
