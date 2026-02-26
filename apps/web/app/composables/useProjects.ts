export function useProjects() {
  const config = useRuntimeConfig()
  const baseUrl = config.public.apiBaseUrl

  function fetchProjects() {
    return useFetch<Project[]>(`${baseUrl}/projects`)
  }

  async function createProject(data: ProjectInput): Promise<Project> {
    return $fetch<Project>(`${baseUrl}/projects`, {
      method: 'POST',
      body: data,
    })
  }

  async function updateProject(id: number, data: ProjectInput): Promise<Project> {
    return $fetch<Project>(`${baseUrl}/projects/${id}`, {
      method: 'PATCH',
      body: data,
    })
  }

  async function deleteProject(id: number): Promise<void> {
    await $fetch(`${baseUrl}/projects/${id}`, {
      method: 'DELETE',
    })
  }

  return { fetchProjects, createProject, updateProject, deleteProject }
}
