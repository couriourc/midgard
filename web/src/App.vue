<template>
  <div class="min-h-screen bg-background">
    <Sidebar :collapsed="sidebarCollapsed" @toggle-collapse="sidebarCollapsed = !sidebarCollapsed" />
    <div :class="['transition-all duration-300', sidebarCollapsed ? 'md:pl-16' : 'md:pl-64']">
      <!-- 移动端菜单按钮 -->
      <div class="md:hidden sticky top-0 z-40 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div class="flex h-16 items-center px-4">
          <button
            class="inline-flex items-center justify-center rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground"
            @click="sidebarOpen = !sidebarOpen"
          >
            <Menu class="h-6 w-6" />
          </button>
          <div class="flex flex-1 items-center justify-center space-x-2">
            <img src="/midgard.png" alt="Midgard Gateway" class="h-6 w-auto object-contain" />
            <span class="text-sm font-medium">Midgard Gateway</span>
          </div>
        </div>
      </div>
      <main class="p-4 md:p-6 lg:p-8">
      <router-view />
    </main>
    </div>
    
    <!-- Mobile Sidebar -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 z-50 md:hidden"
      @click="sidebarOpen = false"
    >
      <div class="fixed inset-y-0 left-0 z-50 w-64 bg-background border-r border-border">
        <div class="flex flex-col h-full">
          <div class="flex items-center justify-between px-4 h-16 border-b border-border">
            <div class="flex items-center space-x-2">
              <img src="/midgard.png" alt="Midgard Gateway" class="h-8 w-auto object-contain" />
              <h1 class="text-xl font-bold">Midgard Gateway</h1>
            </div>
            <button
              class="p-2 rounded-md hover:bg-accent"
              @click="sidebarOpen = false"
            >
              <X class="h-5 w-5" />
            </button>
          </div>
          <nav class="flex-1 px-2 py-4 space-y-1 overflow-y-auto">
            <router-link
              v-for="item in navigation"
              :key="item.name"
              :to="item.href"
              :class="cn(
                'flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors',
                $route.path === item.href
                  ? 'bg-accent text-accent-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              )"
              @click="sidebarOpen = false"
            >
              <component :is="item.icon" class="mr-3 h-5 w-5" />
              {{ item.name }}
            </router-link>
          </nav>
        </div>
      </div>
    </div>
    
    <!-- Toast Container -->
    <Toast />
  </div>
</template>

<script setup>
import { ref } from "vue"
import { LayoutDashboard, FolderOpen, X, Menu } from "lucide-vue-next"
import { cn } from "@/lib/utils"
import Sidebar from "@/components/layout/sidebar.vue"
import Toast from "@/components/ui/toast.vue"

const sidebarOpen = ref(false)
const sidebarCollapsed = ref(false)

const navigation = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard },
  { name: "Collections", href: "/collections", icon: FolderOpen },
]
</script>
