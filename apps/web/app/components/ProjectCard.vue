<script setup lang="ts">
const props = defineProps<{
  project: Project
}>()

const emit = defineEmits<{
  edit: [project: Project]
  delete: [project: Project]
}>()

const statusBadge = computed(() =>
  props.project.status === 'Active'
    ? 'badge-success'
    : 'badge-ghost',
)
</script>

<template>
  <div class="card bg-base-200 shadow-md transition-shadow hover:shadow-lg">
    <NuxtLink :to="`/projects/${project.id}`" class="card-body gap-3">
      <div class="flex items-start justify-between">
        <h2 class="card-title text-base">
          {{ project.name }}
        </h2>
        <span class="badge badge-sm" :class="statusBadge">
          {{ project.status }}
        </span>
      </div>

      <p class="line-clamp-2 text-sm text-base-content/60">
        {{ project.description || 'No description' }}
      </p>

      <div class="flex items-center justify-between">
        <span class="text-xs text-base-content/50">
          {{ project.key_count }} {{ project.key_count === 1 ? 'key' : 'keys' }}
        </span>

        <div class="flex gap-1" @click.prevent>
          <button
            class="btn btn-ghost btn-xs"
            title="Edit"
            @click="emit('edit', project)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </button>
          <button
            class="btn btn-ghost btn-xs text-error"
            title="Delete"
            @click="emit('delete', project)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>
      </div>
    </NuxtLink>
  </div>
</template>
