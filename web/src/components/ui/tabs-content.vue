<template>
  <div
    v-show="isActive"
    :class="cn(
      'mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
      $attrs.class
    )"
    v-bind="$attrs"
  >
    <slot />
  </div>
</template>

<script setup>
import { computed, inject } from "vue"
import { cn } from "@/lib/utils"

const props = defineProps({
  value: {
    type: String,
    required: true,
  },
})

const tabsContext = inject("tabsValue", null)

const isActive = computed(() => {
  if (!tabsContext) return false
  return tabsContext.modelValue() === props.value
})
</script>

