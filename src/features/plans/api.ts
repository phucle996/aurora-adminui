import { z } from 'zod'

const planSchema = z.object({
  id: z.string(),
  resource_type: z.string(),
  resource_model: z.string(),
  code: z.string(),
  name: z.string(),
  description: z.string(),
  status: z.string(),
  vcpu: z.number(),
  ram_gb: z.number(),
  disk_gb: z.number(),
})

const listPlansResponseSchema = z.object({
  message: z.string(),
  data: z.object({
    items: z.array(planSchema),
  }),
})

export type Plan = z.infer<typeof planSchema>

export async function listPlans(): Promise<Plan[]> {
  const response = await fetch('/api/v1/admin/plans', {
    credentials: 'include',
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load plans')
  }

  return listPlansResponseSchema.parse(payload).data.items
}

export async function createPlan(input: {
  resourceType: string
  resourceModel: string
  code: string
  name: string
  description: string
  vcpu: number
  ramGb: number
  diskGb: number
}): Promise<Plan> {
  const response = await fetch('/api/v1/admin/plans', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to create plan')
  }

  return planSchema.parse(payload?.data)
}
