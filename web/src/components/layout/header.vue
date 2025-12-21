<template>
  <header class="sticky top-0 z-40 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
    <div class="flex h-16 items-center px-4 md:px-6">
      <button
        class="inline-flex items-center justify-center rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground md:hidden"
        @click="$emit('toggle-sidebar')"
      >
        <Menu class="h-6 w-6" />
      </button>
      <button
        class="hidden md:inline-flex items-center justify-center rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground"
        @click="$emit('toggle-sidebar-collapse')"
      >
        <PanelLeftClose v-if="!sidebarCollapsed" class="h-5 w-5" />
        <PanelLeftOpen v-else class="h-5 w-5" />
      </button>
      <div class="flex flex-1 items-center justify-end space-x-4">
        <div class="flex items-center space-x-2">
          <span class="text-sm text-muted-foreground">Midgard Gateway</span>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { Menu, PanelLeftClose, PanelLeftOpen } from "lucide-vue-next"

defineProps({
  sidebarCollapsed: {
    type: Boolean,
    default: false
  }
})

defineEmits(["toggle-sidebar", "toggle-sidebar-collapse"])
</script>

