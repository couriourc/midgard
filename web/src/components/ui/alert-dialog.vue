<template>
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="fixed inset-0 z-50 flex items-center justify-center"
      @click.self="handleCancel"
    >
      <div class="fixed inset-0 bg-background/80 backdrop-blur-sm" />
      <div class="relative z-50 w-full max-w-lg rounded-lg border bg-background p-6 shadow-lg">
        <div class="flex flex-col space-y-2 text-center sm:text-left">
          <h2 class="text-lg font-semibold">{{ title }}</h2>
          <p class="text-sm text-muted-foreground">{{ description }}</p>
        </div>
        <div class="flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2 mt-4">
          <Button variant="outline" @click="handleCancel">{{ cancelText }}</Button>
          <Button :variant="confirmVariant" @click="handleConfirm">{{ confirmText }}</Button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { Teleport } from "vue"
import Button from "./button.vue"

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: "Are you sure?",
  },
  description: {
    type: String,
    default: "",
  },
  confirmText: {
    type: String,
    default: "Confirm",
  },
  cancelText: {
    type: String,
    default: "Cancel",
  },
  confirmVariant: {
    type: String,
    default: "default",
  },
})

const emit = defineEmits(["update:modelValue", "confirm", "cancel"])

const handleConfirm = () => {
  emit("confirm")
  emit("update:modelValue", false)
}

const handleCancel = () => {
  emit("cancel")
  emit("update:modelValue", false)
}
</script>

