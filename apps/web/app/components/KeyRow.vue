<script setup lang="ts">
import { useClipboard } from '@vueuse/core'

const props = defineProps<{
  apiKey: ApiKeyMasked
}>()

const emit = defineEmits<{
  'update:enabled': [value: boolean]
  'delete': []
}>()

const toast = useToastContext()
const { revealKey } = useApiKeys()
const { copy } = useClipboard()

const isRevealed = ref(false)
const revealedValue = ref<string | null>(null)
const revealing = ref(false)
const localEnabled = ref(props.apiKey.is_enabled)

watch(() => props.apiKey.is_enabled, (v) => {
  localEnabled.value = v
})

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

async function fetchReveal(): Promise<string | null> {
  if (revealedValue.value)
    return revealedValue.value

  revealing.value = true
  try {
    const result = await revealKey(props.apiKey.id)
    revealedValue.value = result.key_value
    return result.key_value
  }
  catch (err: any) {
    toast?.error(err?.data?.error || 'Failed to reveal key')
    return null
  }
  finally {
    revealing.value = false
  }
}

async function handleRevealToggle() {
  if (revealing.value)
    return

  if (isRevealed.value) {
    isRevealed.value = false
    return
  }

  const value = await fetchReveal()
  if (value)
    isRevealed.value = true
}

async function handleCopy() {
  const value = revealedValue.value ?? await fetchReveal()
  if (!value)
    return

  try {
    await copy(value)
    toast?.success('Copied to clipboard')
  }
  catch {
    toast?.error('Failed to copy')
  }
}

function handleEnabledUpdate(value: boolean) {
  localEnabled.value = value
  emit('update:enabled', value)
}
</script>

<template>
  <tr class="hover:bg-base-200/50">
    <td class="font-medium">
      {{ apiKey.name }}
    </td>
    <td>
      <code class="break-all font-mono text-sm">
        <template v-if="revealing">
          <span class="loading loading-dots loading-xs" />
        </template>
        <template v-else-if="isRevealed && revealedValue">
          {{ revealedValue }}
        </template>
        <template v-else>
          {{ apiKey.key_value_masked }}
        </template>
      </code>
    </td>
    <td>
      <ToggleSwitch
        :model-value="localEnabled"
        :key-id="apiKey.id"
        @update:model-value="handleEnabledUpdate"
      />
    </td>
    <td class="text-sm text-base-content/60">
      {{ formatDate(apiKey.created_at) }}
    </td>
    <td class="text-sm text-base-content/60">
      {{ apiKey.last_used_at ? formatDate(apiKey.last_used_at) : 'Never' }}
    </td>
    <td>
      <div class="flex gap-1">
        <button
          class="btn btn-ghost btn-xs"
          title="Reveal"
          :disabled="revealing"
          @click="handleRevealToggle"
        >
          <svg v-if="isRevealed" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21" />
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
        </button>
        <button
          class="btn btn-ghost btn-xs"
          title="Copy"
          @click="handleCopy"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
        </button>
        <button
          class="btn btn-ghost btn-xs text-error"
          title="Delete"
          @click="emit('delete')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </td>
  </tr>
</template>
