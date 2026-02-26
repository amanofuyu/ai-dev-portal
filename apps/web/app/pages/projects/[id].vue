<script setup lang="ts">
const route = useRoute()
const toast = useToastContext()
const config = useRuntimeConfig()
const baseUrl = config.public.apiBaseUrl

const projectId = computed(() => Number(route.params.id))

const { data: project, status: projectStatus, error: projectError } = useFetch<Project>(
  () => `${baseUrl}/projects/${projectId.value}`,
)

const { fetchKeys, deleteKey } = useApiKeys()
const { data: keys, status: keysStatus, error: keysError, refresh: refreshKeys } = fetchKeys(projectId.value)

const showCreateModal = ref(false)
const showEditForm = ref(false)
const showDeleteConfirm = ref(false)
const deletingKey = ref<ApiKeyMasked | null>(null)

function openDeleteConfirm(key: ApiKeyMasked) {
  deletingKey.value = key
  showDeleteConfirm.value = true
}

async function handleDeleteKey() {
  if (!deletingKey.value)
    return
  try {
    await deleteKey(deletingKey.value.id)
    toast?.success('API key deleted')
    showDeleteConfirm.value = false
    deletingKey.value = null
    await refreshKeys()
  }
  catch (err: any) {
    toast?.error(err?.data?.error || 'Failed to delete key')
  }
}

function onKeyCreated() {
  showCreateModal.value = false
  refreshKeys()
}

function onProjectSaved() {
  showEditForm.value = false
  refreshNuxtData()
}
</script>

<template>
  <div>
    <!-- Back navigation -->
    <NuxtLink to="/" class="btn btn-ghost btn-sm mb-6 gap-1">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
      </svg>
      Back to Projects
    </NuxtLink>

    <!-- Loading -->
    <div v-if="projectStatus === 'pending'" class="flex flex-col gap-6">
      <div class="skeleton h-32 rounded-xl" />
      <div class="skeleton h-64 rounded-xl" />
    </div>

    <!-- Error -->
    <div v-else-if="projectError" class="alert alert-error">
      <span>Failed to load project: {{ projectError.message }}</span>
    </div>

    <!-- Content -->
    <template v-else-if="project">
      <!-- Project Info Card -->
      <div class="card mb-8 bg-base-200">
        <div class="card-body">
          <div class="flex items-start justify-between">
            <div>
              <h1 class="text-2xl font-bold">
                {{ project.name }}
              </h1>
              <p v-if="project.description" class="mt-1 text-base-content/60">
                {{ project.description }}
              </p>
            </div>
            <div class="flex items-center gap-2">
              <span
                class="badge badge-sm"
                :class="project.status === 'Active' ? 'badge-success' : 'badge-ghost'"
              >
                {{ project.status }}
              </span>
              <button class="btn btn-ghost btn-sm" title="Edit" @click="showEditForm = true">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- API Keys Section -->
      <div class="flex flex-col gap-4">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">
            API Keys
          </h2>
          <button class="btn btn-primary btn-sm" @click="showCreateModal = true">
            + Generate New Key
          </button>
        </div>

        <!-- Keys Loading -->
        <div v-if="keysStatus === 'pending'" class="skeleton h-48 rounded-xl" />

        <!-- Keys Error -->
        <div v-else-if="keysError" class="alert alert-error">
          <span>Failed to load keys: {{ keysError.message }}</span>
          <button class="btn btn-ghost btn-sm" @click="refreshKeys()">
            Retry
          </button>
        </div>

        <!-- Keys Table -->
        <KeyTable
          v-else-if="keys"
          :keys="keys"
          @create="showCreateModal = true"
          @delete="openDeleteConfirm"
        />
      </div>

      <!-- Edit Project Modal -->
      <ProjectForm
        v-if="showEditForm"
        :project="project"
        @close="showEditForm = false"
        @saved="onProjectSaved"
      />

      <!-- Create Key Modal -->
      <KeyCreateModal
        v-if="showCreateModal"
        :project-id="projectId"
        @close="showCreateModal = false"
        @created="onKeyCreated"
      />

      <!-- Delete Key Confirmation -->
      <DeleteConfirm
        v-if="showDeleteConfirm"
        title="Delete API Key"
        :message="`Are you sure you want to delete the key &quot;${deletingKey?.name}&quot;? This action cannot be undone.`"
        @close="showDeleteConfirm = false; deletingKey = null"
        @confirm="handleDeleteKey"
      />
    </template>
  </div>
</template>
