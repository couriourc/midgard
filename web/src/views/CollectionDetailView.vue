<template>
  <div class="space-y-6">
    <div v-if="loading" class="flex items-center justify-center py-12">
      <p class="text-muted-foreground">Loading...</p>
    </div>
    <div v-else-if="collection">
      <div class="flex items-center justify-between">
          <div>
          <h2 class="text-3xl font-bold tracking-tight">{{ collection.name }}</h2>
          <p class="text-muted-foreground mt-1">{{ collection.description || "No description" }}</p>
        </div>
        <div class="flex space-x-2">
          <Button variant="outline" @click="showEditModal = true">
            <Edit class="mr-2 h-4 w-4" />
            Edit
          </Button>
          <Button @click="showImportModal = true">
            <Upload class="mr-2 h-4 w-4" />
            Import OpenAPI
          </Button>
        </div>
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Basic Information</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div>
              <label class="text-sm font-medium text-muted-foreground">Prefix</label>
              <p class="mt-1">
                <code class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm">{{ collection.prefix }}</code>
              </p>
            </div>
              <div>
              <label class="text-sm font-medium text-muted-foreground">Base URL</label>
              <p class="mt-1 text-sm">{{ collection.base_url }}</p>
              </div>
              <div>
              <label class="text-sm font-medium text-muted-foreground">Proxy URL</label>
              <div class="mt-1 flex items-center space-x-2">
                <code class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-xs">
                  {{ proxyUrl }}
                </code>
                <button
                  @click="copyProxyUrl(collection.prefix)"
                  class="p-1 hover:bg-accent rounded"
                  title="Copy proxy URL"
                >
                  <Copy class="h-3 w-3 text-muted-foreground" />
                </button>
              </div>
            </div>
            <div>
              <label class="text-sm font-medium text-muted-foreground">Status</label>
              <p class="mt-1">
                <Badge :variant="collection.active ? 'default' : 'secondary'">
                  {{ collection.active ? "Active" : "Inactive" }}
                </Badge>
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Configuration</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
              <div>
              <label class="text-sm font-medium text-muted-foreground">Health Check</label>
              <p class="mt-1 text-sm">{{ collection.health_path || "Not configured" }}</p>
            </div>
        <div>
              <label class="text-sm font-medium text-muted-foreground">Logging</label>
              <p class="mt-1">
                <Badge :variant="collection.log_enabled ? 'default' : 'secondary'">
                  {{ collection.log_enabled ? "Enabled" : "Disabled" }}
                </Badge>
              </p>
            </div>
              <div>
              <label class="text-sm font-medium text-muted-foreground">Caching</label>
              <p class="mt-1">
                <Badge :variant="collection.cache_enabled ? 'default' : 'secondary'">
                  {{ collection.cache_enabled ? "Enabled" : "Disabled" }}
                </Badge>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Endpoints ({{ collection.endpoints ? collection.endpoints.length : 0 }})</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <thead>
              <tr class="border-b">
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Method</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Path</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Summary</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">请求次数</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">平均耗时</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="endpoint in collection.endpoints"
                :key="endpoint.id"
                class="border-b transition-colors hover:bg-muted/50"
              >
                <td class="p-4 align-middle">
                  <Badge :variant="getMethodColorVariant(endpoint.method)" :class="getMethodColorClass(endpoint.method)">
                    {{ endpoint.method }}
                  </Badge>
                </td>
                <td class="p-4 align-middle font-mono text-sm">{{ endpoint.path }}</td>
                <td class="p-4 align-middle text-sm text-muted-foreground">{{ endpoint.summary || "-" }}</td>
                <td class="p-4 align-middle text-sm">
                  {{ getEndpointStat(endpoint.path, endpoint.method)?.request_count || 0 }}
                </td>
                <td class="p-4 align-middle text-sm">
                  {{ getEndpointStat(endpoint.path, endpoint.method)?.avg_duration 
                    ? `${Math.round(getEndpointStat(endpoint.path, endpoint.method).avg_duration)}ms`
                    : '-' }}
                </td>
              </tr>
              <tr v-if="!collection.endpoints || collection.endpoints.length === 0">
                <td colspan="5" class="p-8 text-center text-muted-foreground">
                  No endpoints. Import OpenAPI to add endpoints.
                </td>
              </tr>
            </tbody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div class="flex items-center justify-between">
            <CardTitle>Request Logs</CardTitle>
            <Button variant="ghost" size="sm" @click="fetchLogs">
              <RefreshCw class="h-4 w-4" />
            </Button>
      </div>
        </CardHeader>
        <CardContent>
          <Table>
          <thead>
              <tr class="border-b">
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Method</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Path</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Status</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Duration</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Size</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Client IP</th>
                <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Timestamp</th>
            </tr>
          </thead>
          <tbody>
              <tr
                v-for="log in logs"
                :key="log.id"
                class="border-b transition-colors hover:bg-muted/50"
              >
                <td class="p-4 align-middle">
                  <Badge :variant="getMethodVariant(log.method)" variant="outline">
                    {{ log.method }}
                  </Badge>
                </td>
                <td class="p-4 align-middle font-mono text-sm">{{ log.path }}</td>
                <td class="p-4 align-middle">
                  <Badge :variant="getStatusVariant(log.status)">{{ log.status }}</Badge>
                </td>
                <td class="p-4 align-middle text-sm">{{ log.duration }}ms</td>
                <td class="p-4 align-middle text-sm">{{ formatBytes(log.response_size) }}</td>
                <td class="p-4 align-middle text-sm text-muted-foreground">{{ log.client_ip }}</td>
                <td class="p-4 align-middle text-sm text-muted-foreground">{{ formatDate(log.timestamp) }}</td>
              </tr>
              <tr v-if="logs.length === 0">
                <td colspan="7" class="p-8 text-center text-muted-foreground">No logs yet</td>
            </tr>
          </tbody>
          </Table>
        </CardContent>
      </Card>
    </div>

    <!-- Edit Modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 z-50 flex items-center justify-center"
      @click.self="showEditModal = false"
    >
      <div class="fixed inset-0 bg-background/80 backdrop-blur-sm" />
      <Card class="relative z-50 w-full max-w-3xl max-h-[90vh] overflow-hidden flex flex-col">
        <CardHeader>
          <CardTitle>Edit Collection</CardTitle>
        </CardHeader>
        <CardContent class="flex-1 overflow-y-auto">
          <form @submit.prevent="updateCollection">
            <Tabs v-model="editTab" class="w-full">
              <TabsList class="grid w-full grid-cols-4">
                <TabsTrigger value="basic">Basic</TabsTrigger>
                <TabsTrigger value="health">Health</TabsTrigger>
                <TabsTrigger value="logging">Logging</TabsTrigger>
                <TabsTrigger value="caching">Caching</TabsTrigger>
              </TabsList>
              
              <TabsContent value="basic" class="space-y-4 mt-4">
                <div class="space-y-2">
                  <label class="text-sm font-medium">Name *</label>
                  <Input v-model="editForm.name" required />
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Description</label>
                  <textarea
                    v-model="editForm.description"
                    class="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  />
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Prefix *</label>
                  <div class="flex items-center space-x-2">
                    <Input 
                      v-model="editForm.prefix" 
                      required 
                      placeholder="api-v1"
                      :class="editPrefixError ? 'border-red-500' : ''"
                      @input="() => { checkEditPrefix(); }"
                    />
                    <Button 
                      type="button" 
                      variant="outline" 
                      size="sm"
                      @click="generateEditPrefixFromName"
                      title="根据名称自动生成"
                    >
                      自动生成
                    </Button>
                  </div>
                  <p v-if="editPrefixError" class="text-xs text-red-500">{{ editPrefixError }}</p>
                  <p v-else class="text-xs text-muted-foreground">External gateway prefix (e.g., api-v1). 不能包含 '/' 字符，将自动替换为 '-'</p>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Base URL *</label>
                  <Input v-model="editForm.base_url" required />
                  <p class="text-xs text-muted-foreground">Target service base URL</p>
                </div>
                <div class="flex items-center space-x-2">
                  <input v-model="editForm.active" type="checkbox" id="active" class="h-4 w-4 rounded border-gray-300" />
                  <label for="active" class="text-sm font-medium">Active</label>
                </div>
              </TabsContent>
              
              <TabsContent value="health" class="space-y-4 mt-4">
                <div class="space-y-2">
                  <label class="text-sm font-medium">Health Check Path</label>
                  <Input v-model="editForm.health_path" placeholder="/health" />
                  <p class="text-xs text-muted-foreground">Path to check service health status</p>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Health Check Interval (seconds)</label>
                  <Input v-model.number="editForm.health_interval" type="number" min="10" />
                  <p class="text-xs text-muted-foreground">How often to check health status</p>
                </div>
              </TabsContent>
              
              <TabsContent value="logging" class="space-y-4 mt-4">
                <div class="flex items-center space-x-2">
                  <input v-model="editForm.log_enabled" type="checkbox" id="edit_log_enabled" class="h-4 w-4 rounded border-gray-300" />
                  <label for="edit_log_enabled" class="text-sm font-medium">Enable Logging</label>
                </div>
                <div class="flex items-center space-x-2">
                  <input v-model="editForm.log_rolling" type="checkbox" id="edit_log_rolling" class="h-4 w-4 rounded border-gray-300" />
                  <label for="edit_log_rolling" class="text-sm font-medium">Enable Rolling Logs</label>
                  <p class="text-xs text-muted-foreground">Automatically remove old logs when limit is reached</p>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Max Log Entries</label>
                  <Input v-model.number="editForm.log_max_entries" type="number" min="100" />
                  <p class="text-xs text-muted-foreground">Maximum number of log entries to keep</p>
                </div>
              </TabsContent>
              
              <TabsContent value="caching" class="space-y-4 mt-4">
                <div class="flex items-center space-x-2">
                  <input v-model="editForm.cache_enabled" type="checkbox" id="edit_cache_enabled" class="h-4 w-4 rounded border-gray-300" />
                  <label for="edit_cache_enabled" class="text-sm font-medium">Enable Caching</label>
                </div>
                <div v-if="editForm.cache_enabled" class="space-y-2">
                  <label class="text-sm font-medium">Cache TTL (seconds)</label>
                  <Input v-model.number="editForm.cache_ttl" type="number" min="60" />
                  <p class="text-xs text-muted-foreground">Time to live for cached responses</p>
                </div>
                <div v-if="editForm.cache_enabled" class="space-y-2">
                  <label class="text-sm font-medium">Cache Key Strategy</label>
                  <select v-model="editForm.cache_key_strategy" class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
                    <option value="params">Query Parameters Only</option>
                    <option value="body">Request Body Only</option>
                    <option value="all">All (Params + Body)</option>
                  </select>
                  <p class="text-xs text-muted-foreground">How to generate cache keys</p>
                </div>
              </TabsContent>
            </Tabs>
            
            <div class="flex justify-end space-x-2 mt-6 pt-4 border-t">
              <Button type="button" variant="outline" @click="showEditModal = false">Cancel</Button>
              <Button type="submit">Update</Button>
      </div>
          </form>
        </CardContent>
      </Card>
    </div>

    <!-- Import OpenAPI Modal -->
    <div
      v-if="showImportModal"
      class="fixed inset-0 z-50 flex items-center justify-center"
      @click.self="showImportModal = false"
    >
      <div class="fixed inset-0 bg-background/80 backdrop-blur-sm" />
      <Card class="relative z-50 w-full max-w-2xl">
        <CardHeader>
          <CardTitle>Import OpenAPI</CardTitle>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="importOpenAPI" class="space-y-4">
            <div class="space-y-2">
              <label class="text-sm font-medium">OpenAPI URL</label>
              <Input v-model="importForm.openapi_url" type="url" placeholder="https://api.example.com/openapi.json" />
        </div>
            <div class="space-y-2">
              <label class="text-sm font-medium">Or Paste OpenAPI JSON</label>
              <textarea
                v-model="importForm.openapi_json"
                class="flex min-h-[200px] w-full rounded-md border border-input bg-background px-3 py-2 font-mono text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                placeholder='{"openapi": "3.0.0", ...}'
              />
          </div>
          <div class="flex justify-end space-x-2">
              <Button type="button" variant="outline" @click="showImportModal = false">Cancel</Button>
              <Button type="submit">Import</Button>
          </div>
          </form>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import axios from "axios"
import { pinyin } from "pinyin-pro"
import { Edit, Upload, RefreshCw, Copy } from "lucide-vue-next"
import Card from "@/components/ui/card.vue"
import CardHeader from "@/components/ui/card-header.vue"
import CardTitle from "@/components/ui/card-title.vue"
import CardContent from "@/components/ui/card-content.vue"
import Button from "@/components/ui/button.vue"
import Input from "@/components/ui/input.vue"
import Table from "@/components/ui/table.vue"
import Badge from "@/components/ui/badge.vue"
import Tabs from "@/components/ui/tabs.vue"
import TabsList from "@/components/ui/tabs-list.vue"
import TabsTrigger from "@/components/ui/tabs-trigger.vue"
import TabsContent from "@/components/ui/tabs-content.vue"
import { useToast } from "@/composables/useToast"

const route = useRoute()
const router = useRouter()

const collection = ref(null)
const logs = ref([])
const endpointStats = ref([])
const loading = ref(true)
const showEditModal = ref(false)
const showImportModal = ref(false)
const editTab = ref("basic")
const { toast } = useToast()

const proxyUrl = computed(() => {
  if (!collection.value) return ""
  if (typeof window !== "undefined") {
    return `${window.location.origin}/proxy/${collection.value.prefix}/*`
  }
  return `/proxy/${collection.value.prefix}/*`
})

const editForm = ref({})
const editPrefixError = ref("")
const editPrefixChecking = ref(false)
const importForm = ref({
  openapi_url: "",
  openapi_json: "",
})

// 使用 pinyin-pro 将字符串转换为拼音并生成 Prefix
const toPinyin = (str) => {
  if (!str) return ""
  
  let result = ""
  for (let i = 0; i < str.length; i++) {
    const char = str[i]
    const code = char.charCodeAt(0)
    
    // ASCII 字母数字和连字符直接保留
    if ((code >= 48 && code <= 57) || // 0-9
        (code >= 65 && code <= 90) || // A-Z
        (code >= 97 && code <= 122) || // a-z
        char === '-' || char === '_') {
      result += char.toLowerCase()
    } else if (code >= 0x4e00 && code <= 0x9fff) {
      // 中文字符，使用 pinyin-pro 获取拼音首字母
      try {
        const pinyinStr = pinyin(char, { toneType: 'none', type: 'first' })
        if (pinyinStr) {
          result += pinyinStr.toLowerCase()
        }
      } catch (error) {
        console.warn("Failed to convert pinyin:", error)
      }
    } else if (char === ' ' || char === '/') {
      // 空格和斜杠转换为连字符
      result += '-'
    }
    // 其他字符忽略
  }
  
  // 清理：移除多余的连字符，转换为小写
  result = result.toLowerCase().replace(/-+/g, '-').replace(/^-|-$/g, '')
  
  return result
}

// 生成编辑表单的 Prefix
const generateEditPrefixFromName = () => {
  if (!editForm.value.name) {
    toast.error("请先输入名称", "名称不能为空")
    return
  }
  
  let prefix = toPinyin(editForm.value.name)
  if (!prefix) {
    prefix = "collection"
  }
  
  editForm.value.prefix = prefix
  checkEditPrefix()
}

// 检查编辑表单的 Prefix 唯一性
const checkEditPrefix = async () => {
  let prefix = editForm.value.prefix
  
  // 清理前缀：替换 / 为 -
  if (prefix.includes('/')) {
    prefix = prefix.replace(/\//g, '-')
    editForm.value.prefix = prefix
  }
  
  if (!prefix || prefix.trim() === '') {
    editPrefixError.value = ""
    return
  }
  
  if (prefix.includes('/')) {
    editPrefixError.value = "Prefix 不能包含 '/' 字符，已自动替换为 '-'"
    return
  }
  
  editPrefixChecking.value = true
  try {
    const excludeId = collection.value?.id || ""
    const response = await axios.get(`/api/collections/check-prefix/${encodeURIComponent(prefix)}?exclude_id=${excludeId}`)
    if (response.data.exists) {
      editPrefixError.value = "该 Prefix 已存在，请使用其他值"
    } else {
      editPrefixError.value = ""
    }
  } catch (error) {
    console.error("Failed to check prefix:", error)
    editPrefixError.value = ""
  } finally {
    editPrefixChecking.value = false
  }
}

const fetchCollection = async () => {
  try {
    loading.value = true
    const response = await axios.get(`/api/collections/${route.params.id}`)
    collection.value = response.data
    editForm.value = {
      name: collection.value.name || "",
      description: collection.value.description || "",
      prefix: collection.value.prefix || "",
      base_url: collection.value.base_url || "",
      health_path: collection.value.health_path || "",
      health_interval: collection.value.health_interval || 30,
      log_enabled: collection.value.log_enabled !== undefined ? collection.value.log_enabled : true,
      log_rolling: collection.value.log_rolling !== undefined ? collection.value.log_rolling : true,
      log_max_entries: collection.value.log_max_entries || 1000,
      cache_enabled: collection.value.cache_enabled !== undefined ? collection.value.cache_enabled : false,
      cache_ttl: collection.value.cache_ttl || 300,
      cache_key_strategy: collection.value.cache_key_strategy || "all",
      active: collection.value.active !== undefined ? collection.value.active : true,
    }
    // 重置错误状态
    editPrefixError.value = ""
    await Promise.all([fetchLogs(), fetchEndpointStats()])
  } catch (error) {
    console.error("Failed to fetch collection:", error)
    toast.error("Failed to fetch collection", error.message)
    router.push("/collections")
  } finally {
    loading.value = false
  }
}

const fetchEndpointStats = async () => {
  try {
    const response = await axios.get(`/api/collections/${route.params.id}/endpoint-stats`)
    endpointStats.value = response.data
  } catch (error) {
    console.error("Failed to fetch endpoint stats:", error)
  }
}

const getEndpointStat = (path, method) => {
  return endpointStats.value.find(stat => stat.path === path && stat.method === method)
}

const fetchLogs = async () => {
  try {
    const response = await axios.get(`/api/logs/${route.params.id}?limit=50`)
    logs.value = response.data
  } catch (error) {
    console.error("Failed to fetch logs:", error)
  }
}

const updateCollection = async (e) => {
  if (e) {
    e.preventDefault()
  }
  
  // 验证 prefix
  if (editPrefixError.value) {
    toast.error("Prefix 验证失败", editPrefixError.value)
    return
  }
  
  // 确保 prefix 不包含斜杠
  if (editForm.value.prefix.includes('/')) {
    editForm.value.prefix = editForm.value.prefix.replace(/\//g, '-')
  }
  
  try {
    // Ensure all fields are present
    const payload = {
      name: editForm.value.name,
      description: editForm.value.description || "",
      prefix: editForm.value.prefix,
      base_url: editForm.value.base_url,
      health_path: editForm.value.health_path || "",
      health_interval: editForm.value.health_interval || 30,
      log_enabled: editForm.value.log_enabled !== undefined ? editForm.value.log_enabled : true,
      log_rolling: editForm.value.log_rolling !== undefined ? editForm.value.log_rolling : true,
      log_max_entries: editForm.value.log_max_entries || 1000,
      cache_enabled: editForm.value.cache_enabled !== undefined ? editForm.value.cache_enabled : false,
      cache_ttl: editForm.value.cache_ttl || 300,
      cache_key_strategy: editForm.value.cache_key_strategy || "all",
      active: editForm.value.active !== undefined ? editForm.value.active : true,
    }
    await axios.put(`/api/collections/${route.params.id}`, payload)
    await fetchCollection()
    showEditModal.value = false
    editPrefixError.value = ""
    toast.success("Collection updated", "Collection has been updated successfully")
  } catch (error) {
    console.error("Failed to update collection:", error)
    if (error.response && error.response.status === 409) {
      editPrefixError.value = error.response.data.error || "Prefix 已存在"
      toast.error("更新失败", error.response.data.error || "Prefix 已存在")
    } else {
      toast.error("Failed to update collection", error.message)
    }
  }
}

const importOpenAPI = async () => {
  try {
    const payload = {
      openapi_url: importForm.value.openapi_url || undefined,
      openapi_json: importForm.value.openapi_json ? JSON.parse(importForm.value.openapi_json) : undefined,
    }
    await axios.post(`/api/collections/${route.params.id}/import-openapi`, payload)
    await fetchCollection()
    showImportModal.value = false
    importForm.value = { openapi_url: "", openapi_json: "" }
    toast.success("OpenAPI imported", "OpenAPI specification has been imported successfully")
  } catch (error) {
    console.error("Failed to import OpenAPI:", error)
    toast.error("Failed to import OpenAPI", error.message)
  }
}

const formatDate = (dateString) => {
  if (!dateString) return "-"
  const date = new Date(dateString)
  return date.toLocaleString()
}

const formatBytes = (bytes) => {
  if (!bytes || bytes === 0) return "0 B"
  const k = 1024
  const sizes = ["B", "KB", "MB", "GB"]
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
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

const getMethodColorVariant = (method) => {
  // Use outline variant for all methods to apply custom colors
  return "outline"
}

const getMethodColorClass = (method) => {
  const colorMap = {
    GET: "bg-blue-100 text-blue-800 border-blue-300 dark:bg-blue-900 dark:text-blue-200 dark:border-blue-700",
    POST: "bg-green-100 text-green-800 border-green-300 dark:bg-green-900 dark:text-green-200 dark:border-green-700",
    PUT: "bg-yellow-100 text-yellow-800 border-yellow-300 dark:bg-yellow-900 dark:text-yellow-200 dark:border-yellow-700",
    PATCH: "bg-orange-100 text-orange-800 border-orange-300 dark:bg-orange-900 dark:text-orange-200 dark:border-orange-700",
    DELETE: "bg-red-100 text-red-800 border-red-300 dark:bg-red-900 dark:text-red-200 dark:border-red-700",
  }
  return colorMap[method] || "bg-gray-100 text-gray-800 border-gray-300 dark:bg-gray-900 dark:text-gray-200 dark:border-gray-700"
}

const getStatusVariant = (status) => {
  if (status >= 200 && status < 300) return "default"
  if (status >= 400 && status < 500) return "secondary"
  return "destructive"
}

const copyProxyUrl = (prefix) => {
  const url = `${window.location.origin}/proxy/${prefix}/*`
  navigator.clipboard.writeText(url)
  toast.success("Copied", "Proxy URL copied to clipboard")
}

let logInterval = null

onMounted(() => {
  fetchCollection()
  // Auto refresh logs and stats every 5 seconds
  logInterval = setInterval(() => {
    if (collection.value) {
      fetchLogs()
      fetchEndpointStats()
    }
  }, 5000)
})

onUnmounted(() => {
  if (logInterval) {
    clearInterval(logInterval)
  }
})
</script>
