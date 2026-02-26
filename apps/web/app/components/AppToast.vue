<script setup lang="ts">
const toast = useToastContext()
const toasts = computed(() => toast?.toasts.value ?? [])

const alertClass: Record<string, string> = {
  success: 'alert-success',
  error: 'alert-error',
  info: 'alert-info',
}
</script>

<template>
  <div class="fixed right-4 top-4 z-50 flex flex-col gap-2">
    <TransitionGroup
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="translate-x-full opacity-0"
      enter-to-class="translate-x-0 opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="translate-x-0 opacity-100"
      leave-to-class="translate-x-full opacity-0"
    >
      <div
        v-for="item in toasts"
        :key="item.id"
        class="alert w-80 shadow-lg"
        :class="alertClass[item.type]"
      >
        <span>{{ item.message }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>
