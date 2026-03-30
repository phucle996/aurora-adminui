import { createFileRoute } from '@tanstack/react-router'
import { AddPlanPage } from '@/features/plans/new'

export const Route = createFileRoute('/_authenticated/plans/new')({
  component: AddPlanPage,
})
