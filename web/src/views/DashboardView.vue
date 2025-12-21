<template>
  <div :class="['space-y-6', isFullscreen ? 'fixed inset-0 z-[100] bg-background p-4 md:p-6 lg:p-8 overflow-y-auto' : '']">
    <div v-if="!isFullscreen">
      <h2 class="text-3xl font-bold tracking-tight">Dashboard</h2>
      <p class="text-muted-foreground">Overview of your gateway collections and requests</p>
    </div>

    <div v-if="!isFullscreen" class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Collections</CardTitle>
          <FolderOpen class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ collections.length }}</div>
          <p class="text-xs text-muted-foreground">Active gateway collections</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Active Collections</CardTitle>
          <CheckCircle class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ activeCollections }}</div>
          <p class="text-xs text-muted-foreground">{{ collections.length - activeCollections }} inactive</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Requests</CardTitle>
          <FileText class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ requestCount }}</div>
          <p class="text-xs text-muted-foreground">All time requests</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Endpoints</CardTitle>
          <Network class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalEndpoints }}</div>
          <p class="text-xs text-muted-foreground">Across all collections</p>
        </CardContent>
      </Card>
    </div>

    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <CardTitle>Recent Requests</CardTitle>
          <div class="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              @click="toggleFullscreen"
            >
              <Maximize v-if="!isFullscreen" class="h-4 w-4 mr-2" />
              <Minimize v-else class="h-4 w-4 mr-2" />
              {{ isFullscreen ? '退出全屏' : '全屏' }}
            </Button>
            <Button
              variant="outline"
              size="sm"
              @click="toggleAutoRefresh"
              :class="autoRefreshEnabled ? 'bg-primary text-primary-foreground' : ''"
            >
              <RefreshCw :class="['h-4 w-4 mr-2', autoRefreshEnabled ? 'animate-spin' : '']" />
              {{ autoRefreshEnabled ? '停止自动刷新' : '开启自动刷新' }}
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <!-- 筛选区域 -->
        <div class="mb-4 space-y-3">
          <div class="flex flex-col sm:flex-row gap-3">
            <div class="flex-1">
              <label class="text-sm font-medium mb-1 block">接口名 (Path)</label>
              <Input
                v-model="filters.path"
                placeholder="搜索接口路径..."
                @input="handleFilterChange"
              />
            </div>
            <div class="flex-1">
              <label class="text-sm font-medium mb-1 block">Collection</label>
              <select
                v-model="filters.collectionId"
                @change="handleFilterChange"
                class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              >
                <option value="">全部 Collection</option>
                <option
                  v-for="collection in collections"
                  :key="collection.id"
                  :value="collection.id"
                >
                  {{ collection.name }}
                </option>
              </select>
            </div>
          </div>
        </div>

        <div class="space-y-2">
          <div
            v-for="log in recentLogs"
            :key="log.id"
            class="border rounded-lg transition-colors hover:bg-muted/50"
          >
            <div
              class="flex items-center justify-between p-4 cursor-pointer"
              @click="toggleLogExpand(log.id)"
            >
              <div class="flex items-center space-x-4 flex-1">
                <ChevronDown
                  :class="['h-4 w-4 transition-transform text-muted-foreground', expandedLogs[log.id] ? 'rotate-180' : '']"
                />
                <div class="flex flex-col flex-1">
                  <div class="flex items-center space-x-2 flex-wrap">
                    <Badge :variant="getMethodVariant(log.method)" variant="outline" class="text-xs">
                      {{ log.method }}
                    </Badge>
                    <p class="text-sm font-medium">{{ log.path }}</p>
                    <Badge 
                      variant="secondary" 
                      class="text-xs cursor-pointer hover:bg-secondary/80" 
                      @click.stop="goToCollection(log.collection_id)"
                    >
                      {{ getCollectionName(log.collection_id) }}
                    </Badge>
                    <Badge v-if="log.from_cache" variant="secondary" class="text-xs">
                      缓存
                    </Badge>
                  </div>
                  <div class="flex items-center space-x-2 mt-1">
                    <p class="text-xs text-muted-foreground">
                      {{ formatDate(log.timestamp) }}
                    </p>
                    <span class="text-xs text-muted-foreground">•</span>
                    <p class="text-xs text-muted-foreground">
                      {{ log.duration !== null && log.duration !== undefined ? `${log.duration}ms` : '未返回' }}
                    </p>
                  </div>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <Badge :variant="getStatusVariant(log.status)">
                  {{ log.status }}
                </Badge>
              </div>
            </div>
            <div v-if="expandedLogs[log.id]" class="border-t bg-muted/30 p-4 space-y-3">
              <div class="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <span class="text-muted-foreground">目标 URL:</span>
                  <code class="ml-2 text-xs">{{ log.target_url }}</code>
                </div>
                <div>
                  <span class="text-muted-foreground">客户端 IP:</span>
                  <span class="ml-2">{{ log.client_ip }}</span>
                </div>
                <div>
                  <span class="text-muted-foreground">请求大小:</span>
                  <span class="ml-2">{{ formatBytes(log.request_size) }}</span>
                </div>
                <div>
                  <span class="text-muted-foreground">响应大小:</span>
                  <span class="ml-2">{{ formatBytes(log.response_size) }}</span>
                </div>
              </div>
              <div v-if="log.request_params || log.request_body" class="space-y-2">
                <div v-if="log.request_params">
                  <div class="text-sm font-medium mb-2">请求参数 (Params):</div>
                  <pre class="bg-background border rounded p-3 text-xs overflow-auto max-h-40">{{ formatJSON(log.request_params) }}</pre>
                </div>
                <div v-if="log.request_body">
                  <div class="flex items-center justify-between mb-2">
                    <div class="text-sm font-medium">请求体 (Body):</div>
                    <button
                      @click.stop="copyRequestBody(log.request_body)"
                      class="p-1 hover:bg-accent rounded transition-colors"
                      title="复制请求体"
                    >
                      <Copy class="h-4 w-4 text-muted-foreground" />
                    </button>
                  </div>
                  <pre class="bg-background border rounded p-3 text-xs overflow-auto max-h-40">{{ formatJSON(log.request_body) }}</pre>
                </div>
              </div>
              <div v-else class="text-sm text-muted-foreground">
                无请求参数
              </div>
            </div>
          </div>
          <div v-if="recentLogs.length === 0" class="text-sm text-muted-foreground text-center py-4">
            No requests yet
          </div>
        </div>

        <!-- 分页组件 -->
        <div v-if="totalLogs > 0" class="mt-4 flex items-center justify-between">
          <div class="text-sm text-muted-foreground">
            显示 {{ (currentPage - 1) * pageSize + 1 }} - {{ Math.min(currentPage * pageSize, totalLogs) }} 条，共 {{ totalLogs }} 条
          </div>
          <div class="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage === 1"
              @click="goToPage(currentPage - 1)"
            >
              上一页
            </Button>
            <div class="flex items-center space-x-1">
              <Button
                v-for="page in visiblePages"
                :key="page"
                variant="outline"
                size="sm"
                :class="page === currentPage ? 'bg-primary text-primary-foreground' : ''"
                @click="goToPage(page)"
              >
                {{ page }}
              </Button>
            </div>
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage >= totalPages"
              @click="goToPage(currentPage + 1)"
            >
              下一页
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue"
import { useRouter } from "vue-router"
import axios from "axios"
import { FolderOpen, CheckCircle, FileText, Network, RefreshCw, ChevronDown, Maximize, Minimize, Copy } from "lucide-vue-next"
import Card from "@/components/ui/card.vue"
import CardHeader from "@/components/ui/card-header.vue"
import CardTitle from "@/components/ui/card-title.vue"
import CardContent from "@/components/ui/card-content.vue"
import Badge from "@/components/ui/badge.vue"
import Button from "@/components/ui/button.vue"
import Input from "@/components/ui/input.vue"
import { useToast } from "@/composables/useToast"

const router = useRouter()
const { toast } = useToast()
const collections = ref([])
const requestLogs = ref([])
const expandedLogs = ref({})
const autoRefreshEnabled = ref(false)
const isFullscreen = ref(false)
let refreshInterval = null

// 分页和筛选状态
const currentPage = ref(1)
const pageSize = ref(5)
const totalLogs = ref(0)
const filters = ref({
  path: "",
  collectionId: "",
})

const fetchCollections = async () => {
  try {
    const response = await axios.get("/api/collections")
    collections.value = response.data
  } catch (error) {
    console.error("Failed to fetch collections:", error)
  }
}

const fetchLogs = async () => {
  try {
    const params = new URLSearchParams({
      page: currentPage.value.toString(),
      pageSize: pageSize.value.toString(),
    })
    
    if (filters.value.path) {
      params.append("path", filters.value.path)
    }
    if (filters.value.collectionId) {
      params.append("collection_id", filters.value.collectionId)
    }

    const response = await axios.get(`/api/logs?${params.toString()}`)
    
    // 处理新的 API 响应格式
    if (response.data.data) {
      requestLogs.value = response.data.data
      totalLogs.value = response.data.total || 0
      currentPage.value = response.data.page || 1
      pageSize.value = response.data.pageSize || 5
    } else {
      // 兼容旧格式（如果没有 data 字段，说明是旧 API）
      requestLogs.value = response.data
      totalLogs.value = response.data.length || 0
    }
  } catch (error) {
    console.error("Failed to fetch logs:", error)
  }
}

const activeCollections = computed(() => {
  return collections.value.filter((c) => c.active).length
})

const requestCount = computed(() => {
  return totalLogs.value || requestLogs.value.length
})

const totalEndpoints = computed(() => {
  return collections.value.reduce((sum, c) => sum + (c.endpoints?.length || 0), 0)
})

const recentLogs = computed(() => {
  return requestLogs.value
})

const totalPages = computed(() => {
  return Math.ceil(totalLogs.value / pageSize.value)
})

const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = Math.min(totalPages.value, start + maxVisible - 1)
  
  if (end - start < maxVisible - 1) {
    start = Math.max(1, end - maxVisible + 1)
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

const getCollectionName = (collectionId) => {
  const collection = collections.value.find(c => c.id === collectionId)
  return collection ? collection.name : 'Unknown'
}

const toggleLogExpand = (logId) => {
  expandedLogs.value[logId] = !expandedLogs.value[logId]
}

const goToCollection = (collectionId) => {
  router.push(`/collections/${collectionId}`)
}

const toggleAutoRefresh = () => {
  autoRefreshEnabled.value = !autoRefreshEnabled.value
  if (autoRefreshEnabled.value) {
    refreshInterval = setInterval(() => {
      fetchLogs()
    }, 5000)
  } else {
    if (refreshInterval) {
      clearInterval(refreshInterval)
      refreshInterval = null
    }
  }
}

const handleFilterChange = () => {
  currentPage.value = 1 // 重置到第一页
  fetchLogs()
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchLogs()
  }
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
  // 当进入全屏时，可以隐藏 body 滚动条
  if (isFullscreen.value) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
}

const copyRequestBody = (requestBody) => {
  try {
    // 如果 requestBody 是 JSON 字符串，先格式化再复制
    const formattedBody = formatJSON(requestBody)
    navigator.clipboard.writeText(formattedBody)
    toast.success("已复制", "请求体已复制到剪贴板")
  } catch (error) {
    console.error("Failed to copy request body:", error)
    toast.error("复制失败", error.message)
  }
}

const formatJSON = (str) => {
  if (!str) return ""
  try {
    const parsed = JSON.parse(str)
    return JSON.stringify(parsed, null, 2)
  } catch {
    return str
  }
}

const formatBytes = (bytes) => {
  if (!bytes || bytes === 0) return "0 B"
  const k = 1024
  const sizes = ["B", "KB", "MB", "GB"]
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
}

const formatDate = (dateString) => {
  if (!dateString) return "-"
  const date = new Date(dateString)
  return date.toLocaleString()
}

const getMethodVariant = (method) => {
  const variants = {
    GET: "default",
    POST: "default",
    PUT: "default",
    DELETE: "destructive",
  }
  return variants[method] || "secondary"
}

const getStatusVariant = (status) => {
  if (status >= 200 && status < 300) return "default"
  if (status >= 400 && status < 500) return "secondary"
  return "destructive"
}

onMounted(() => {
  fetchCollections()
  fetchLogs()
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
  // 清理全屏状态
  document.body.style.overflow = ''
})
</script>
