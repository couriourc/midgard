<template>
  <button
    type="button"
    :class="cn(
      'inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
      isActive
        ? 'bg-background text-foreground shadow-sm'
        : 'text-muted-foreground hover:bg-background/50',
      $attrs.class
    )"
    @click.prevent="updateValue"
    v-bind="$attrs"
  >
    <slot />
  </button>
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

const updateValue = () => {
  if (tabsContext) {
    tabsContext.updateValue(props.value)
  }
}
</script>

