<script setup lang="ts">
provideToastContext()
</script>

<template>
  <div class="min-h-screen bg-base-100">
    <nav class="navbar bg-base-200 border-b border-base-300 px-4">
      <div class="mx-auto flex w-full max-w-5xl items-center">
        <NuxtLink to="/" class="flex items-center gap-2 text-lg font-bold text-primary">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
          </svg>
          Dev Portal
        </NuxtLink>
      </div>
    </nav>

    <main class="mx-auto w-full max-w-5xl px-4 py-8">
      <slot />
    </main>

    <AppToast />
  </div>
</template>
