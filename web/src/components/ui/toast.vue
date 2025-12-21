<template>
  <Teleport to="body">
    <TransitionGroup
      name="toast"
      tag="div"
      class="fixed top-0 z-[100] flex max-h-screen w-full flex-col-reverse p-4 sm:bottom-0 sm:right-0 sm:top-auto sm:flex-col md:max-w-[420px]"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        :class="cn(
          'group pointer-events-auto relative flex w-full items-center justify-between space-x-4 overflow-hidden rounded-md border p-6 pr-8 shadow-lg transition-all',
          toastVariantClasses(toast.variant)
        )"
      >
        <div class="grid gap-1">
          <div v-if="toast.title" class="text-sm font-semibold">{{ toast.title }}</div>
          <div v-if="toast.description" class="text-sm opacity-90">{{ toast.description }}</div>
        </div>
        <button
          class="absolute right-2 top-2 rounded-md p-1 text-foreground/50 opacity-0 transition-opacity hover:text-foreground focus:opacity-100 focus:outline-none focus:ring-2 group-hover:opacity-100"
          @click="removeToast(toast.id)"
        >
          <X class="h-4 w-4" />
        </button>
      </div>
    </TransitionGroup>
  </Teleport>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue"
import { Teleport, TransitionGroup } from "vue"
import { X } from "lucide-vue-next"
import { cn } from "@/lib/utils"

const toasts = ref([])
let toastIdCounter = 0

const toastVariantClasses = (variant = "default") => {
  const variants = {
    default: "border bg-background text-foreground",
    destructive: "destructive group border-destructive bg-destructive text-destructive-foreground",
    success: "border-green-500 bg-green-50 text-green-900 dark:bg-green-900 dark:text-green-50",
    warning: "border-yellow-500 bg-yellow-50 text-yellow-900 dark:bg-yellow-900 dark:text-yellow-50",
  }
  return variants[variant] || variants.default
}

const showToast = (toast) => {
  const id = ++toastIdCounter
  const newToast = {
    id,
    title: toast.title || "",
    description: toast.description || "",
    variant: toast.variant || "default",
    duration: toast.duration || 5000,
  }
  
  toasts.value.push(newToast)
  
  if (newToast.duration > 0) {
    setTimeout(() => {
      removeToast(id)
    }, newToast.duration)
  }
  
  return id
}

const removeToast = (id) => {
  const index = toasts.value.findIndex((t) => t.id === id)
  if (index > -1) {
    toasts.value.splice(index, 1)
  }
}

// 提供全局方法
const toast = {
  success: (title, description) => showToast({ title, description, variant: "success" }),
  error: (title, description) => showToast({ title, description, variant: "destructive" }),
  warning: (title, description) => showToast({ title, description, variant: "warning" }),
  info: (title, description) => showToast({ title, description, variant: "default" }),
  show: (options) => showToast(options),
}

// 暴露给全局
if (typeof window !== "undefined") {
  window.toast = toast
}

defineExpose({
  toast,
  showToast,
  removeToast,
})
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>

