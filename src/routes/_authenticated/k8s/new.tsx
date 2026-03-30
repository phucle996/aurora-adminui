import { createFileRoute } from '@tanstack/react-router'
import { AddK8sClusterPage } from '@/features/k8s/new'

export const Route = createFileRoute('/_authenticated/k8s/new')({
  component: AddK8sClusterPage,
})
