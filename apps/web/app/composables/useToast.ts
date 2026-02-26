import type { InjectionKey } from 'vue'

export interface ToastItem {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
}

export interface ToastContext {
  toasts: Ref<ToastItem[]>
  success: (message: string) => void
  error: (message: string) => void
  info: (message: string) => void
}

const TOAST_KEY = Symbol('toast') as InjectionKey<ToastContext>

let nextId = 0

export function provideToastContext(): ToastContext {
  const toasts = ref<ToastItem[]>([])

  function add(message: string, type: ToastItem['type']) {
    const id = nextId++
    toasts.value.push({ id, message, type })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, 3000)
  }

  const context: ToastContext = {
    toasts,
    success: (msg: string) => add(msg, 'success'),
    error: (msg: string) => add(msg, 'error'),
    info: (msg: string) => add(msg, 'info'),
  }

  provide(TOAST_KEY, context)
  return context
}

export function useToastContext(): ToastContext | null {
  return inject(TOAST_KEY, null)
}
