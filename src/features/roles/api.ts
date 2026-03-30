import { z } from 'zod'

const roleSchema = z.object({
  id: z.string(),
  name: z.string(),
  scope: z.string(),
  description: z.string(),
  userCount: z.number(),
  permissionCount: z.number(),
})

const roleListResponseSchema = z.object({
  message: z.string(),
  data: z.object({
    items: z.array(roleSchema),
  }),
})

const permissionSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
})

const permissionListResponseSchema = z.object({
  message: z.string(),
  data: z.object({
    items: z.array(permissionSchema),
  }),
})

export type Role = z.infer<typeof roleSchema>
export type Permission = z.infer<typeof permissionSchema>

export async function listRoles(): Promise<Role[]> {
  const response = await fetch('/api/v1/admin/roles', {
    credentials: 'include',
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load roles')
  }

  return roleListResponseSchema.parse(payload).data.items
}

export async function listPermissions(): Promise<Permission[]> {
  const response = await fetch('/api/v1/admin/roles/permissions', {
    credentials: 'include',
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load permissions')
  }

  return permissionListResponseSchema.parse(payload).data.items
}

export async function createRole(input: {
  name: string
  description: string
  permissionIds: string[]
}): Promise<Role> {
  const response = await fetch('/api/v1/admin/roles', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  const payload = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to create role')
  }

  return roleSchema.parse(payload?.data)
}
