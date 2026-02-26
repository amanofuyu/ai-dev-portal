<script setup lang="ts">
import { useClipboard } from '@vueuse/core'

const props = defineProps<{
  projectId: number
}>()

const emit = defineEmits<{
  close: []
  created: []
}>()

const toast = useToastContext()
const { createKey } = useApiKeys()
const { copy } = useClipboard()

const name = ref('')
const loading = ref(false)
const createdKey = ref<ApiKeyFull | null>(null)

async function handleSubmit() {
  if (!name.value.trim()) {
    toast?.error('Key name is required')
    return
  }

  loading.value = true
  try {
    createdKey.value = await createKey(props.projectId, name.value.trim())
    toast?.success('API key created')
  }
  catch (err: any) {
    toast?.error(err?.data?.error || 'Failed to create key')
  }
  finally {
    loading.value = false
  }
}

async function handleCopy() {
  if (!createdKey.value)
    return
  try {
    await copy(createdKey.value.key_value)
    toast?.success('Copied to clipboard')
  }
  catch {
    toast?.error('Failed to copy')
  }
}

function handleClose() {
  if (createdKey.value)
    emit('created')
  emit('close')
}
</script>

<template>
  <div class="modal modal-open">
    <div class="modal-box">
      <!-- Input state -->
      <template v-if="!createdKey">
        <h3 class="text-lg font-bold">
          Generate New API Key
        </h3>
        <form class="mt-4 flex flex-col gap-4" @submit.prevent="handleSubmit">
          <label class="form-control w-full">
            <div class="label">
              <span class="label-text">Key Name *</span>
            </div>
            <input
              v-model="name"
              type="text"
              placeholder="e.g. Production Key"
              class="input input-bordered w-full"
              required
            >
          </label>
          <div class="modal-action">
            <button type="button" class="btn" @click="handleClose">
              Cancel
            </button>
            <button type="submit" class="btn btn-primary" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner loading-sm" />
              Generate
            </button>
          </div>
        </form>
      </template>

      <!-- Result state -->
      <template v-else>
        <h3 class="text-lg font-bold text-success">
          API Key Created
        </h3>

        <div class="mt-4 flex flex-col gap-4">
          <div class="alert alert-warning">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <span>Please copy this key now. You won't be able to see the full key again after closing this dialog.</span>
          </div>

          <div class="rounded-lg bg-base-300 p-4">
            <p class="mb-1 text-xs text-base-content/50">
              {{ createdKey.name }}
            </p>
            <code class="break-all font-mono text-sm text-primary">
              {{ createdKey.key_value }}
            </code>
          </div>

          <button class="btn btn-outline btn-sm gap-2" @click="handleCopy">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            Copy Key
          </button>

          <div class="modal-action">
            <button class="btn" @click="handleClose">
              Done
            </button>
          </div>
        </div>
      </template>
    </div>
    <div class="modal-backdrop" @click="handleClose" />
  </div>
</template>
