<template>
  <input
    :class="cn(
      'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
      $attrs.class
    )"
    :value="value"
    @input="value = $event.target.value"
    v-bind="$attrs"
  />
</template>

<script setup>
import { cn } from "@/lib/utils"

// 简化 defineModel 声明（Vue 3.4+ 推荐写法）
const value = defineModel({
  type: [String, Number, Boolean, Object, Array],
  required: false,
})
</script>
