import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/resource-types/')({
  beforeLoad: () => {
    throw redirect({ to: '/resource-definition' })
  },
})
