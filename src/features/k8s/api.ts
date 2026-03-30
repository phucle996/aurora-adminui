type Envelope<T> = {
  message?: string
  data?: T
}

export type K8sClusterListItem = {
  id: string
  name: string
  description: string
  api_server_url: string
  kubernetes_version: string
  validation_status: 'pending' | 'valid' | 'invalid' | 'unreachable'
  last_validated_at?: string
  supports_dbaas: boolean
  supports_serverless: boolean
  supports_generic_workloads: boolean
  zone_name: string
}

export type K8sClusterDetail = {
  id: string
  name: string
  description: string
  import_mode: string
  api_server_url: string
  current_context: string
  kubernetes_version: string
  validation_status: 'pending' | 'valid' | 'invalid' | 'unreachable'
  last_validated_at?: string
  last_validation_error?: string
  supports_dbaas: boolean
  supports_serverless: boolean
  supports_generic_workloads: boolean
  created_at: string
  zone_id?: string
  zone_name: string
}

export async function listK8sClusters(): Promise<K8sClusterListItem[]> {
  const response = await fetch('/api/v1/admin/k8s/clusters', {
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<{ items: K8sClusterListItem[] }>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load kubernetes clusters')
  }
  return payload?.data?.items || []
}

export async function getK8sClusterDetail(
  id: string
): Promise<K8sClusterDetail> {
  const response = await fetch(`/api/v1/admin/k8s/clusters/${id}`, {
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<K8sClusterDetail>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load kubernetes cluster')
  }
  if (!payload?.data) {
    throw new Error('Cluster detail was not returned')
  }
  return payload.data
}

export async function createK8sCluster(input: {
  name: string
  description: string
  zoneId: string
  supportsDbaas: boolean
  supportsServerless: boolean
  supportsGenericWorkloads: boolean
  kubeconfig: File
}): Promise<K8sClusterDetail> {
  const formData = new FormData()
  formData.append('name', input.name)
  formData.append('description', input.description)
  formData.append('zone_id', input.zoneId)
  formData.append('supports_dbaas', String(input.supportsDbaas))
  formData.append('supports_serverless', String(input.supportsServerless))
  formData.append(
    'supports_generic_workloads',
    String(input.supportsGenericWorkloads)
  )
  formData.append('kubeconfig', input.kubeconfig)

  const response = await fetch('/api/v1/admin/k8s/clusters', {
    method: 'POST',
    credentials: 'include',
    body: formData,
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<K8sClusterDetail>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to import kubernetes cluster')
  }
  if (!payload?.data) {
    throw new Error('Cluster was created but no detail was returned')
  }
  return payload.data
}

export async function revalidateK8sCluster(
  id: string
): Promise<K8sClusterDetail> {
  const response = await fetch(`/api/v1/admin/k8s/clusters/${id}/revalidate`, {
    method: 'POST',
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<K8sClusterDetail>
    | null

  if (!response.ok) {
    throw new Error(
      payload?.message || 'Failed to revalidate kubernetes cluster'
    )
  }
  if (!payload?.data) {
    throw new Error('Cluster was revalidated but no detail was returned')
  }
  return payload.data
}

export async function deleteK8sCluster(id: string): Promise<void> {
  const response = await fetch(`/api/v1/admin/k8s/clusters/${id}`, {
    method: 'DELETE',
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<null>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to delete kubernetes cluster')
  }
}
