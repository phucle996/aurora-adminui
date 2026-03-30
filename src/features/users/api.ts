import { userListSchema, type User } from './data/schema'

type UsersResponse = {
  items: User[]
}

export async function listUsers(): Promise<User[]> {
  const response = await fetch('/api/v1/admin/users', {
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | { message?: string; data?: UsersResponse }
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load users')
  }

  return userListSchema.parse(payload?.data?.items || [])
}
