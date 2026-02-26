<script setup lang="ts">
defineProps<{
  title: string
  message: string
}>()

const emit = defineEmits<{
  close: []
  confirm: []
}>()

const loading = ref(false)

async function handleConfirm() {
  loading.value = true
  emit('confirm')
}
</script>

<template>
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="text-lg font-bold text-error">
        {{ title }}
      </h3>
      <p class="mt-4 text-base-content/70">
        {{ message }}
      </p>
      <div class="modal-action">
        <button class="btn" @click="emit('close')">
          Cancel
        </button>
        <button
          class="btn btn-error"
          :disabled="loading"
          @click="handleConfirm"
        >
          <span v-if="loading" class="loading loading-spinner loading-sm" />
          Delete
        </button>
      </div>
    </div>
    <div class="modal-backdrop" @click="emit('close')" />
  </div>
</template>
