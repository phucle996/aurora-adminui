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
  zone_name: string
}

export type K8sClusterDetail = {
  id: string
  name: string
  description: string
  api_server_url: string
  current_context: string
  kubernetes_version: string
  validation_status: 'pending' | 'valid' | 'invalid' | 'unreachable'
  last_validated_at?: string
  created_at: string
  zone_id?: string
  zone_name: string
  nodes: K8sClusterNode[]
}

export type K8sClusterDetailPageData = {
  id: string
  name: string
  description: string
  api_server_url: string
  current_context: string
  kubernetes_version: string
  validation_status: 'pending' | 'valid' | 'invalid' | 'unreachable'
  last_validated_at?: string
  created_at: string
  zone_id?: string
  zone_name: string
  nodes: K8sClusterNode[]
  zone_options: {
    id: string
    name: string
  }[]
}

export type K8sClusterNode = {
  name: string
  roles: string[]
  kubelet_version: string
  container_runtime: string
  os_image: string
  kernel_version: string
  ready: boolean
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

export async function getK8sClusterDetailPageData(
  id: string
): Promise<K8sClusterDetailPageData> {
  const response = await fetch(`/api/v1/admin/k8s/clusters/${id}/page-data`, {
    credentials: 'include',
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<K8sClusterDetailPageData>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to load kubernetes cluster page')
  }
  if (!payload?.data) {
    throw new Error('Cluster page data was not returned')
  }
  return payload.data
}

export async function createK8sCluster(input: {
  name: string
  description: string
  zoneId: string
  kubeconfig: File
}): Promise<K8sClusterDetail> {
  const formData = new FormData()
  formData.append('name', input.name)
  formData.append('description', input.description)
  formData.append('zone_id', input.zoneId)
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

export async function updateK8sCluster(
  id: string,
  input: {
    zone_id: string
  }
): Promise<K8sClusterDetail> {
  const response = await fetch(`/api/v1/admin/k8s/clusters/${id}`, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  const payload = (await response.json().catch(() => null)) as
    | Envelope<K8sClusterDetail>
    | null

  if (!response.ok) {
    throw new Error(payload?.message || 'Failed to update kubernetes cluster')
  }
  if (!payload?.data) {
    throw new Error('Cluster was updated but no detail was returned')
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
