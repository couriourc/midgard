<template>
  <div
    role="alert"
    :class="cn(
      'relative w-full rounded-lg border p-4 [&>svg~*]:pl-7 [&>svg+div]:translate-y-[-3px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg]:text-foreground',
      variantClasses,
      $attrs.class
    )"
    v-bind="$attrs"
  >
    <component :is="icon" v-if="icon" class="h-4 w-4" />
    <div class="[&_p]:leading-relaxed">
      <slot />
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue"
import { AlertCircle, CheckCircle2, Info, AlertTriangle } from "lucide-vue-next"
import { cn } from "@/lib/utils"

const props = defineProps({
  variant: {
    type: String,
    default: "default",
    validator: (value) => ["default", "destructive", "success", "warning"].includes(value),
  },
})

const variantClasses = computed(() => {
  const variants = {
    default: "bg-background text-foreground",
    destructive: "border-destructive/50 text-destructive dark:border-destructive [&>svg]:text-destructive",
    success: "border-green-500/50 text-green-600 dark:border-green-500 [&>svg]:text-green-600",
    warning: "border-yellow-500/50 text-yellow-600 dark:border-yellow-500 [&>svg]:text-yellow-600",
  }
  return variants[props.variant] || variants.default
})

const icon = computed(() => {
  const icons = {
    default: Info,
    destructive: AlertCircle,
    success: CheckCircle2,
    warning: AlertTriangle,
  }
  return icons[props.variant] || icons.default
})
</script>

