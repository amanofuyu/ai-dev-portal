<script setup lang="ts">
const props = defineProps<{
  project?: Project | null
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const toast = useToastContext()
const { createProject, updateProject } = useProjects()

const isEdit = computed(() => !!props.project)
const name = ref(props.project?.name ?? '')
const description = ref(props.project?.description ?? '')
const status = ref<'Active' | 'Archived'>(props.project?.status ?? 'Active')
const loading = ref(false)
const dropdownOpen = ref(false)

function selectStatus(value: 'Active' | 'Archived') {
  status.value = value
  dropdownOpen.value = false
}

async function handleSubmit() {
  if (!name.value.trim()) {
    toast?.error('Name is required')
    return
  }

  loading.value = true
  try {
    if (isEdit.value && props.project) {
      await updateProject(props.project.id, {
        name: name.value.trim(),
        description: description.value,
        status: status.value,
      })
      toast?.success('Project updated')
    }
    else {
      await createProject({
        name: name.value.trim(),
        description: description.value,
      })
      toast?.success('Project created')
    }
    emit('saved')
  }
  catch (err: any) {
    toast?.error(err?.data?.error || 'Operation failed')
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="text-lg font-bold">
        {{ isEdit ? 'Edit Project' : 'New Project' }}
      </h3>

      <form class="mt-4 flex flex-col gap-4" @submit.prevent="handleSubmit">
        <label class="form-control w-full">
          <div class="label">
            <span class="label-text">Name *</span>
          </div>
          <input
            v-model="name"
            type="text"
            placeholder="Project name"
            class="input input-bordered w-full"
            required
          >
        </label>

        <label class="form-control w-full">
          <div class="label">
            <span class="label-text">Description</span>
          </div>
          <textarea
            v-model="description"
            placeholder="Optional description"
            class="textarea textarea-bordered w-full"
            rows="3"
          />
        </label>

        <div v-if="isEdit" class="form-control w-full">
          <div class="label">
            <span class="label-text">Status</span>
          </div>
          <div class="relative w-full">
            <button
              type="button"
              class="btn btn-block justify-between border border-base-content/20 bg-base-100 font-normal"
              @click="dropdownOpen = !dropdownOpen"
            >
              {{ status }}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            <ul v-if="dropdownOpen" class="menu absolute z-10 mt-1 w-full rounded-lg border border-base-content/20 bg-base-200 p-1 shadow-lg">
              <li>
                <button
                  type="button"
                  :class="{ active: status === 'Active' }"
                  @click="selectStatus('Active')"
                >
                  Active
                </button>
              </li>
              <li>
                <button
                  type="button"
                  :class="{ active: status === 'Archived' }"
                  @click="selectStatus('Archived')"
                >
                  Archived
                </button>
              </li>
            </ul>
          </div>
        </div>

        <div class="modal-action">
          <button type="button" class="btn" @click="emit('close')">
            Cancel
          </button>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            <span v-if="loading" class="loading loading-spinner loading-sm" />
            {{ isEdit ? 'Save' : 'Create' }}
          </button>
        </div>
      </form>
    </div>
    <div class="modal-backdrop" @click="emit('close')" />
  </div>
</template>
