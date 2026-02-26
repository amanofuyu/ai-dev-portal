<script setup lang="ts">
const { fetchProjects, deleteProject } = useProjects()
const { data: projects, status, error, refresh } = fetchProjects()
const toast = useToastContext()

const showForm = ref(false)
const editingProject = ref<Project | null>(null)
const showDeleteConfirm = ref(false)
const deletingProject = ref<Project | null>(null)

function openCreateForm() {
  editingProject.value = null
  showForm.value = true
}

function openEditForm(project: Project) {
  editingProject.value = project
  showForm.value = true
}

function openDeleteConfirm(project: Project) {
  deletingProject.value = project
  showDeleteConfirm.value = true
}

async function handleDelete() {
  if (!deletingProject.value)
    return
  try {
    await deleteProject(deletingProject.value.id)
    toast?.success('Project deleted')
    showDeleteConfirm.value = false
    deletingProject.value = null
    await refresh()
  }
  catch (err: any) {
    toast?.error(err?.data?.error || 'Failed to delete project')
  }
}

function onFormSaved() {
  showForm.value = false
  editingProject.value = null
  refresh()
}
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold">
        Projects
      </h1>
      <button class="btn btn-primary btn-sm" @click="openCreateForm">
        + New Project
      </button>
    </div>

    <!-- Loading -->
    <div v-if="status === 'pending'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="i in 3" :key="i" class="skeleton h-40 rounded-xl" />
    </div>

    <!-- Error -->
    <div v-else-if="error" class="alert alert-error">
      <span>Failed to load projects: {{ error.message }}</span>
      <button class="btn btn-ghost btn-sm" @click="refresh()">
        Retry
      </button>
    </div>

    <!-- Empty -->
    <div v-else-if="projects && projects.length === 0" class="flex flex-col items-center gap-4 py-20 text-base-content/50">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
      </svg>
      <p class="text-lg">
        No projects yet
      </p>
      <button class="btn btn-primary" @click="openCreateForm">
        Create Your First Project
      </button>
    </div>

    <!-- Project Grid -->
    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <ProjectCard
        v-for="project in projects"
        :key="project.id"
        :project="project"
        @edit="openEditForm"
        @delete="openDeleteConfirm"
      />
    </div>

    <!-- Project Form Modal -->
    <ProjectForm
      v-if="showForm"
      :project="editingProject"
      @close="showForm = false"
      @saved="onFormSaved"
    />

    <!-- Delete Confirmation -->
    <DeleteConfirm
      v-if="showDeleteConfirm"
      title="Delete Project"
      message="This will permanently delete the project and all its API keys. This action cannot be undone."
      @close="showDeleteConfirm = false; deletingProject = null"
      @confirm="handleDelete"
    />
  </div>
</template>
