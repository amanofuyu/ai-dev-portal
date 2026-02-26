<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean
  keyId: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const toast = useToastContext()
const { toggleKey } = useApiKeys()
const loading = ref(false)

async function handleToggle() {
  if (loading.value)
    return

  const oldValue = props.modelValue
  const newValue = !oldValue
  emit('update:modelValue', newValue)
  loading.value = true

  try {
    await toggleKey(props.keyId, newValue)
  }
  catch (err: any) {
    emit('update:modelValue', oldValue)
    toast?.error(err?.data?.error || 'Failed to update key status')
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <input
    type="checkbox"
    class="toggle toggle-sm toggle-primary"
    :checked="modelValue"
    :disabled="loading"
    @change="handleToggle"
  >
</template>
