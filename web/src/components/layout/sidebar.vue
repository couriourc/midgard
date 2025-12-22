<template>
  <aside 
    :class="[
      'hidden md:flex md:flex-col md:fixed md:inset-y-0 transition-all duration-300',
      collapsed ? 'md:w-16' : 'md:w-64'
    ]"
  >
    <div class="flex flex-col flex-grow border-r border-border bg-background">
      <div class="flex items-center flex-shrink-0 h-16 border-b border-border" :class="collapsed ? 'px-2 justify-center' : 'px-4'">
        <div v-if="!collapsed" class="flex items-center space-x-3">
          <h1 class="text-xl font-bold whitespace-nowrap">Midgard Gateway</h1>
        </div>
        <div v-else class="flex items-center justify-center w-full">
          <img src="/midgard.png" alt="MG" class="h-8 w-auto object-contain" />
        </div>
      </div>
      <nav class="flex-1 py-4 space-y-1 overflow-y-auto" :class="collapsed ? 'px-2' : 'px-2'">
        <!-- Dashboard -->
        <router-link
          to="/"
          :class="cn(
            'flex items-center rounded-md transition-colors',
            collapsed ? 'px-2 py-2 justify-center' : 'px-3 py-2',
            $route.path === '/'
              ? 'bg-accent text-accent-foreground'
              : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
          )"
          :title="collapsed ? 'Dashboard' : ''"
        >
          <LayoutDashboard :class="['h-5 w-5 flex-shrink-0', collapsed ? '' : 'mr-3']" />
          <span v-if="!collapsed" class="text-sm font-medium whitespace-nowrap overflow-hidden">Dashboard</span>
        </router-link>

        <!-- Collections with submenu -->
        <div>
          <router-link
            to="/collections"
            :class="cn(
              'flex items-center rounded-md transition-colors',
              collapsed ? 'px-2 py-2 justify-center' : 'px-3 py-2',
              $route.path === '/collections'
                ? 'bg-accent text-accent-foreground'
                : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
            )"
            :title="collapsed ? 'Collections' : ''"
          >
            <FolderOpen :class="['h-5 w-5 flex-shrink-0', collapsed ? '' : 'mr-3']" />
            <span v-if="!collapsed" class="text-sm font-medium whitespace-nowrap overflow-hidden">Collections</span>
          </router-link>

          <!-- Collections submenu -->
          <div 
            v-if="!collapsed" 
            class="ml-6 mt-1 space-y-1 overflow-hidden transition-all"
          >
            <div
              v-for="collection in displayedCollections"
              :key="collection.id"
              :class="[
                'flex items-center justify-between px-3 py-1.5 rounded-md transition-colors cursor-pointer group',
                $route.params.id === collection.id 
                  ? 'bg-accent text-accent-foreground' 
                  : 'hover:bg-accent/50 text-muted-foreground'
              ]"
              @click="goToCollection(collection.id)"
            >
              <div class="flex items-center flex-1 min-w-0">
                <div 
                  :class="[
                    'h-2 w-2 rounded-full flex-shrink-0 mr-2',
                    collection.active ? 'bg-green-500' : 'bg-gray-400'
                  ]"
                  :title="collection.active ? '激活' : '未激活'"
                />
                <span 
                  :class="[
                    'text-xs truncate',
                    $route.params.id === collection.id ? 'font-medium' : ''
                  ]"
                  :title="collection.name"
                >
                  {{ collection.name }}
                </span>
              </div>
            </div>
            <div
              v-if="collections.length > defaultDisplayCount"
              class="flex items-center justify-center px-3 py-1.5"
            >
              <button
                class="text-xs text-muted-foreground hover:text-accent-foreground flex items-center transition-colors"
                @click.stop="showAllCollections = !showAllCollections"
              >
                <ChevronDown 
                  :class="['h-3 w-3 mr-1 transition-transform duration-200', showAllCollections ? 'rotate-180' : '']" 
                />
                <span class="whitespace-nowrap">
                  {{ showAllCollections ? '收起' : `查看更多 (${collections.length - defaultDisplayCount})` }}
                </span>
              </button>
            </div>
            <div v-if="collections.length === 0" class="px-3 py-1.5 text-xs text-muted-foreground">
              暂无 Collections
            </div>
          </div>
        </div>
      </nav>
      <!-- 帮助按钮 -->
      <div class="p-2 border-t border-border">
        <button
          class="w-full flex items-center rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
          :class="collapsed ? 'justify-center' : 'px-3'"
          @click="showHelpDialog = true"
          :title="collapsed ? '使用说明' : ''"
        >
          <HelpCircle :class="['h-5 w-5 flex-shrink-0', collapsed ? '' : 'mr-3']" />
          <span v-if="!collapsed" class="text-sm font-medium whitespace-nowrap">使用说明</span>
        </button>
      </div>
      <!-- 折叠按钮 -->
      <div class="p-2 border-t border-border">
        <button
          class="w-full flex items-center justify-center rounded-md p-2 text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
          @click="$emit('toggle-collapse')"
          :title="collapsed ? '展开侧边栏' : '折叠侧边栏'"
        >
          <PanelLeftClose v-if="!collapsed" class="h-5 w-5" />
          <PanelLeftOpen v-else class="h-5 w-5" />
        </button>
      </div>
    </div>
  </aside>

  <!-- Help Dialog -->
  <Teleport to="body">
    <div
      v-if="showHelpDialog"
      class="fixed inset-0 z-50 flex items-center justify-center"
      @click.self="showHelpDialog = false"
    >
      <div class="fixed inset-0 bg-background/80 backdrop-blur-sm" />
      <div class="relative z-50 w-full max-w-3xl max-h-[90vh] rounded-lg border bg-background shadow-lg overflow-hidden flex flex-col">
        <div class="flex items-center justify-between p-6 border-b">
          <h2 class="text-2xl font-semibold">Midgard Gateway 使用说明</h2>
          <button
            class="h-8 w-8 flex items-center justify-center rounded-md hover:bg-accent transition-colors"
            @click="showHelpDialog = false"
          >
            <X class="h-4 w-4" />
          </button>
        </div>
        <div class="flex-1 overflow-y-auto p-6 space-y-6">
          <div>
            <h3 class="text-lg font-semibold mb-2">项目简介</h3>
            <p class="text-sm text-muted-foreground">
              Midgard Gateway 是一个基于 Go、Traefik、Redis 和 Vue 3 的网关代理管理工具。
              支持通过导入 OpenAPI 规范自动生成代理端点，并提供完整的请求日志、健康检查和缓存功能。
            </p>
          </div>

          <div>
            <h3 class="text-lg font-semibold mb-2">快速开始</h3>
            <ol class="list-decimal list-inside space-y-2 text-sm text-muted-foreground">
              <li>创建 Collection：点击 "Collections" 页面，创建新的集合并配置基本信息</li>
              <li>设置代理前缀：为每个 Collection 设置唯一的 Prefix（如：api-v1），系统会根据名称自动生成</li>
              <li>配置目标服务：设置 Base URL 指向要代理的后端服务</li>
              <li>导入 OpenAPI（可选）：上传或输入 OpenAPI JSON，自动生成端点列表</li>
              <li>启用 Collection：确保 Collection 处于激活状态</li>
              <li>访问代理：通过 <code class="bg-muted px-1 py-0.5 rounded">/proxy/{prefix}/{path}</code> 访问代理服务</li>
            </ol>
          </div>

          <div>
            <h3 class="text-lg font-semibold mb-2">主要功能</h3>
            <div class="space-y-3 text-sm">
              <div>
                <h4 class="font-medium mb-1">1. Collection 管理</h4>
                <p class="text-muted-foreground">
                  创建、编辑、删除集合。每个 Collection 相当于一组代理端点，可以设置对外的网关前缀。
                  Prefix 不能包含 '/' 字符，系统会自动替换为 '-'。
                </p>
              </div>
              <div>
                <h4 class="font-medium mb-1">2. 健康检查</h4>
                <p class="text-muted-foreground">
                  配置健康检查路径和间隔，系统会自动监控后端服务状态。如果服务不健康，代理请求将被拒绝。
                </p>
              </div>
              <div>
                <h4 class="font-medium mb-1">3. 日志记录</h4>
                <p class="text-muted-foreground">
                  记录所有请求的详细信息，包括路径、方法、状态码、耗时、请求参数和响应等。
                  支持滚动日志和条目限制，避免日志过多占用存储空间。
                </p>
              </div>
              <div>
                <h4 class="font-medium mb-1">4. 缓存支持</h4>
                <p class="text-muted-foreground">
                  通过 Redis 缓存请求响应，可配置缓存策略（基于参数、请求体或全部）。
                  缓存命中时会显著提升响应速度。
                </p>
              </div>
            </div>
          </div>

          <div>
            <h3 class="text-lg font-semibold mb-2">代理请求格式</h3>
            <div class="bg-muted p-4 rounded-lg">
              <code class="text-sm">
                {base_url}/proxy/{prefix}/{path}
              </code>
            </div>
            <p class="text-sm text-muted-foreground mt-2">
              其中 <code class="bg-muted px-1 py-0.5 rounded">{prefix}</code> 是 Collection 配置的对外网关前缀，
              <code class="bg-muted px-1 py-0.5 rounded">{path}</code> 是具体的 API 路径。
            </p>
          </div>

          <div>
            <h3 class="text-lg font-semibold mb-2">Dashboard 功能</h3>
            <ul class="list-disc list-inside space-y-2 text-sm text-muted-foreground">
              <li><strong>统计信息</strong>：显示总集合数、活跃集合数、总请求数和总端点数</li>
              <li><strong>最近请求</strong>：查看最近的请求日志，支持展开查看详细信息</li>
              <li><strong>筛选和分页</strong>：可按接口名和 Collection 筛选，支持分页浏览</li>
              <li><strong>自动刷新</strong>：可开启自动刷新功能，实时查看最新请求</li>
              <li><strong>全屏模式</strong>：支持全屏查看请求日志，便于详细分析</li>
            </ul>
          </div>

          <div>
            <h3 class="text-lg font-semibold mb-2">注意事项</h3>
            <ul class="list-disc list-inside space-y-2 text-sm text-muted-foreground">
              <li>Prefix 必须唯一，系统会在创建时自动检查</li>
              <li>Prefix 不能包含 '/' 字符，系统会自动替换为 '-'</li>
              <li>确保目标服务的 Base URL 正确配置</li>
              <li>健康检查路径是可选的，但建议配置以确保服务可用性</li>
              <li>缓存功能需要 Redis 服务正常运行</li>
            </ul>
          </div>
        </div>
        <div class="flex justify-end p-6 border-t">
          <button
            class="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
            @click="showHelpDialog = false"
          >
            我知道了
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted } from "vue"
import { Teleport } from "vue"
import { useRouter } from "vue-router"
import { cn } from "@/lib/utils"
import { LayoutDashboard, FolderOpen, PanelLeftClose, PanelLeftOpen, ChevronDown, HelpCircle, X } from "lucide-vue-next"
import axios from "axios"

const router = useRouter()

defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})

defineEmits(["toggle-collapse"])

const collections = ref([])
const showAllCollections = ref(false)
const defaultDisplayCount = 5
const showHelpDialog = ref(false)

const displayedCollections = computed(() => {
  if (showAllCollections.value) {
    return collections.value
  }
  return collections.value.slice(0, defaultDisplayCount)
})

const fetchCollections = async () => {
  try {
    const response = await axios.get("/api/collections")
    collections.value = response.data
  } catch (error) {
    console.error("Failed to fetch collections:", error)
  }
}

const goToCollection = (id) => {
  router.push(`/collections/${id}`)
}

onMounted(() => {
  fetchCollections()
  // 定期刷新 collections 列表
  setInterval(() => {
    fetchCollections()
  }, 30000) // 每30秒刷新一次
})
</script>

