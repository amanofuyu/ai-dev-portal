<script setup lang="ts">
defineProps<{
  keys: ApiKeyMasked[]
}>()

const emit = defineEmits<{
  create: []
  delete: [key: ApiKeyMasked]
}>()
</script>

<template>
  <div v-if="keys.length === 0" class="flex flex-col items-center gap-4 rounded-xl border border-base-300 py-16 text-base-content/50">
    <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1">
      <path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
    </svg>
    <p class="text-lg">
      No API keys yet
    </p>
    <button class="btn btn-primary btn-sm" @click="emit('create')">
      Generate Your First Key
    </button>
  </div>

  <div v-else class="overflow-x-auto rounded-xl border border-base-300">
    <table class="table">
      <thead>
        <tr class="bg-base-200/50">
          <th>Name</th>
          <th class="min-w-56">
            Key
          </th>
          <th>Status</th>
          <th>Created</th>
          <th>Last Used</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <KeyRow
          v-for="key in keys"
          :key="key.id"
          :api-key="key"
          @delete="emit('delete', key)"
        />
      </tbody>
    </table>
  </div>
</template>
