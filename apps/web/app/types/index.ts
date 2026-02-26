export interface Project {
  id: number
  name: string
  description: string
  status: 'Active' | 'Archived'
  key_count: number
  created_at: string
  updated_at: string
}

export interface ProjectInput {
  name?: string
  description?: string
  status?: 'Active' | 'Archived'
}

export interface ApiKeyMasked {
  id: number
  key_value_masked: string
  name: string
  is_enabled: boolean
  last_used_at: string | null
  project_id: number
  created_at: string
}

export interface ApiKeyFull {
  id: number
  key_value: string
  name: string
  is_enabled: boolean
  last_used_at: string | null
  project_id: number
  created_at: string
}

export interface ApiKeyRevealed {
  id: number
  key_value: string
}

export interface ApiError {
  error: string
}
