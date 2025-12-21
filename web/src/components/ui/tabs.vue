<template>
  <div :class="cn('w-full', $attrs.class)" v-bind="$attrs">
    <slot :modelValue="modelValue" />
  </div>
</template>

<script setup>
import { provide } from "vue"
import { cn } from "@/lib/utils"

const props = defineProps({
  modelValue: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(["update:modelValue"])

provide("tabsValue", {
  modelValue: () => props.modelValue,
  updateValue: (value) => emit("update:modelValue", value),
})
</script>

