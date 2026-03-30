import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, CheckSquare, Plus, ShieldAlert } from 'lucide-react'
import { toast } from 'sonner'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { createRole, listPermissions } from './api'

export function AddRolePage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [selectedPermissionIDs, setSelectedPermissionIDs] = useState<string[]>(
    []
  )

  const permissionsQuery = useQuery({
    queryKey: ['role-permissions'],
    queryFn: listPermissions,
  })

  const createRoleMutation = useMutation({
    mutationFn: createRole,
    onSuccess: (role) => {
      toast.success(`Role "${role.name}" created`)
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      navigate({ to: '/roles' })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to create role'
      )
    },
  })

  const selectedCount = selectedPermissionIDs.length
  const permissions = permissionsQuery.data || []
  const groupedPermissions = useMemo(() => {
    const groups = new Map<string, typeof permissions>()
    for (const permission of permissions) {
      const prefix = permission.name.split('.')[0] || 'other'
      const existing = groups.get(prefix) || []
      existing.push(permission)
      groups.set(prefix, existing)
    }
    return Array.from(groups.entries()).sort(([left], [right]) =>
      left.localeCompare(right)
    )
  }, [permissions])

  function togglePermission(permissionID: string, checked: boolean) {
    setSelectedPermissionIDs((current) => {
      if (checked) {
        return current.includes(permissionID)
          ? current
          : [...current, permissionID]
      }
      return current.filter((item) => item !== permissionID)
    })
  }

  function handleCreateRole() {
    if (!name.trim()) {
      toast.error('Role name is required')
      return
    }
    createRoleMutation.mutate({
      name: name.trim(),
      description: description.trim(),
      permissionIds: selectedPermissionIDs,
    })
  }

  return (
    <>
      <Header fixed>
        <div className='min-w-0'>
          <p className='subtle-kicker'>IAM</p>
          <h1 className='truncate text-lg font-semibold text-foreground'>
            Add Role
          </h1>
        </div>
        <div className='ms-auto flex items-center space-x-4'>
          <Search />
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main className='flex flex-col gap-6'>
        <section className='page-header'>
          <div className='space-y-2'>
            <p className='subtle-kicker'>Access model</p>
            <h1 className='page-title'>Create a new role</h1>
            <p className='page-copy'>
              Define a global IAM role, then activate the permissions that
              should belong to it.
            </p>
          </div>
          <Button variant='outline' asChild>
            <Link to='/roles'>
              <ArrowLeft className='size-4' />
              Back to roles
            </Link>
          </Button>
        </section>

        <div className='grid gap-6 xl:grid-cols-[minmax(0,420px)_minmax(0,1fr)]'>
          <Card>
            <CardHeader>
              <CardTitle>Role definition</CardTitle>
            </CardHeader>
            <CardContent className='space-y-5'>
              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Role name
                </label>
                <Input
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder='Example: smtp-operator'
                />
              </div>

              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Scope
                </label>
                <Input value='global' disabled />
              </div>

              <div className='space-y-2'>
                <label className='text-sm font-medium text-foreground'>
                  Description
                </label>
                <Textarea
                  className='min-h-32'
                  value={description}
                  onChange={(event) => setDescription(event.target.value)}
                  placeholder='Describe what this role is allowed to operate and who should use it.'
                />
              </div>

              <div className='rounded-xl border border-border/80 bg-muted/40 px-4 py-3'>
                <p className='text-sm font-medium text-foreground'>
                  Selected permissions
                </p>
                <p className='mt-1 text-sm text-muted-foreground'>
                  {selectedCount} active permission
                  {selectedCount === 1 ? '' : 's'}
                </p>
              </div>

              <div className='flex flex-wrap gap-3'>
                <Button
                  onClick={handleCreateRole}
                  disabled={
                    createRoleMutation.isPending || permissionsQuery.isLoading
                  }
                >
                  <Plus className='size-4' />
                  {createRoleMutation.isPending ? 'Creating...' : 'Create role'}
                </Button>
                <Button variant='outline' asChild>
                  <Link to='/roles'>Cancel</Link>
                </Button>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className='space-y-3'>
              <div className='flex items-center gap-3'>
                <span className='flex size-11 items-center justify-center rounded-2xl bg-accent text-accent-foreground'>
                  <CheckSquare className='size-5' />
                </span>
                <div>
                  <p className='subtle-kicker'>Permission catalog</p>
                  <CardTitle>Activate permissions for this role</CardTitle>
                </div>
              </div>
            </CardHeader>
            <CardContent className='space-y-5'>
              {permissionsQuery.isLoading ? (
                <div className='rounded-xl border border-border/80 bg-muted/50 px-4 py-8 text-sm text-muted-foreground'>
                  Loading permissions...
                </div>
              ) : permissionsQuery.isError ? (
                <div className='flex items-start gap-3 rounded-xl border border-warning/25 bg-warning-soft px-4 py-4 text-sm text-warning'>
                  <ShieldAlert className='mt-0.5 size-4 shrink-0' />
                  <span>
                    {permissionsQuery.error instanceof Error
                      ? permissionsQuery.error.message
                      : 'Failed to load permissions'}
                  </span>
                </div>
              ) : groupedPermissions.length > 0 ? (
                groupedPermissions.map(([groupName, groupPermissions]) => (
                  <div key={groupName} className='space-y-3'>
                    <div>
                      <p className='text-sm font-semibold uppercase tracking-[0.18em] text-muted-foreground'>
                        {groupName}
                      </p>
                    </div>
                    <div className='grid gap-3 md:grid-cols-2'>
                      {groupPermissions.map((permission) => {
                        const checked = selectedPermissionIDs.includes(
                          permission.id
                        )
                        return (
                          <label
                            key={permission.id}
                            className='flex cursor-pointer items-start gap-3 rounded-xl border border-border/80 bg-muted/30 px-4 py-4 transition-colors hover:border-primary/40 hover:bg-muted/50'
                          >
                            <Checkbox
                              checked={checked}
                              onCheckedChange={(value) =>
                                togglePermission(permission.id, value === true)
                              }
                              className='mt-0.5'
                            />
                            <div className='min-w-0 space-y-1'>
                              <p className='text-sm font-medium text-foreground'>
                                {permission.name}
                              </p>
                              <p className='text-sm leading-6 text-muted-foreground'>
                                {permission.description || 'No description'}
                              </p>
                            </div>
                          </label>
                        )
                      })}
                    </div>
                  </div>
                ))
              ) : (
                <div className='rounded-xl border border-border/80 bg-muted/50 px-4 py-8 text-sm text-muted-foreground'>
                  No permissions found.
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </Main>
    </>
  )
}
