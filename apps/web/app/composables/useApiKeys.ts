export function useApiKeys() {
  const config = useRuntimeConfig()
  const baseUrl = config.public.apiBaseUrl

  function fetchKeys(projectId: number | Ref<number>) {
    const id = toValue(projectId)
    return useFetch<ApiKeyMasked[]>(`${baseUrl}/projects/${id}/keys`)
  }

  async function createKey(projectId: number, name: string): Promise<ApiKeyFull> {
    return $fetch<ApiKeyFull>(`${baseUrl}/projects/${projectId}/keys`, {
      method: 'POST',
      body: { name },
    })
  }

  async function toggleKey(id: number, isEnabled: boolean): Promise<ApiKeyMasked> {
    return $fetch<ApiKeyMasked>(`${baseUrl}/keys/${id}`, {
      method: 'PATCH',
      body: { is_enabled: isEnabled },
    })
  }

  async function deleteKey(id: number): Promise<void> {
    await $fetch(`${baseUrl}/keys/${id}`, {
      method: 'DELETE',
    })
  }

  async function revealKey(id: number): Promise<ApiKeyRevealed> {
    return $fetch<ApiKeyRevealed>(`${baseUrl}/keys/${id}/reveal`)
  }

  return { fetchKeys, createKey, toggleKey, deleteKey, revealKey }
}
