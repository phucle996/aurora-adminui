import { createFileRoute } from '@tanstack/react-router'
import { AddRolePage } from '@/features/roles/new'

export const Route = createFileRoute('/_authenticated/roles/new')({
  component: AddRolePage,
})
