<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold tracking-tight">Collections</h2>
        <p class="text-muted-foreground">Manage your gateway collections</p>
      </div>
      <Button @click="showCreateModal = true">
        <Plus class="mr-2 h-4 w-4" />
        Add Collection
      </Button>
    </div>

    <Card>
      <CardContent class="p-0">
        <Table>
          <thead>
            <tr class="border-b">
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground w-12"></th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Name</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Prefix</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Proxy URL</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Base URL</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Status</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Endpoints</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="collection in collections" :key="collection.id">
              <tr class="border-b transition-colors hover:bg-muted/50">
                <td class="p-4 align-middle">
                  <Button
                    variant="ghost"
                    size="sm"
                    @click="toggleExpand(collection.id)"
                    class="h-8 w-8 p-0"
                  >
                    <ChevronDown
                      :class="['h-4 w-4 transition-transform', expandedRows[collection.id] ? 'rotate-180' : '']"
                    />
                  </Button>
                </td>
                <td class="p-4 align-middle font-medium">{{ collection.name }}</td>
                <td class="p-4 align-middle">
                  <code class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm">{{ collection.prefix }}</code>
                </td>
                <td class="p-4 align-middle">
                  <div class="flex items-center space-x-2">
                    <code class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-xs">
                      /proxy/{{ collection.prefix }}/*
                    </code>
                    <button
                      @click="copyProxyUrl(collection.prefix)"
                      class="p-1 hover:bg-accent rounded"
                      title="Copy proxy URL"
                    >
                      <Copy class="h-3 w-3 text-muted-foreground" />
                    </button>
                  </div>
                </td>
                <td class="p-4 align-middle text-sm text-muted-foreground">{{ collection.base_url }}</td>
                <td class="p-4 align-middle">
                  <Badge :variant="collection.active ? 'default' : 'secondary'">
                  {{ collection.active ? 'Active' : 'Inactive' }}
                  </Badge>
                </td>
                <td class="p-4 align-middle">{{ collection.endpoints ? collection.endpoints.length : 0 }}</td>
                <td class="p-4 align-middle">
                  <div class="flex items-center space-x-2">
                    <Button variant="ghost" size="sm" @click="$router.push(`/collections/${collection.id}`)">
                      <Eye class="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      @click="toggleCollection(collection)"
                    >
                      <component :is="collection.active ? Pause : Play" class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm" @click="deleteCollection(collection.id)">
                      <Trash2 class="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </td>
              </tr>
              <tr v-if="expandedRows[collection.id]" class="border-b bg-muted/30">
                <td colspan="8" class="p-4">
                  <div v-if="latestLogs[collection.id] === null" class="text-sm text-muted-foreground">
                    暂无日志记录
                  </div>
                  <div v-else-if="latestLogs[collection.id]" class="space-y-4">
                    <div class="flex items-center justify-between">
                      <h4 class="font-semibold">最新日志记录</h4>
                      <Button
                        variant="ghost"
                        size="sm"
                        @click="toggleLogDetail(collection.id)"
                      >
                        <ChevronDown
                          :class="['h-4 w-4 transition-transform mr-2', expandedLogDetails[collection.id] ? 'rotate-180' : '']"
                        />
                        {{ expandedLogDetails[collection.id] ? '收起' : '展开详情' }}
                      </Button>
                    </div>
                    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                      <div>
                        <span class="text-muted-foreground">方法:</span>
                        <Badge :variant="getMethodVariant(latestLogs[collection.id].method)" variant="outline" class="ml-2">
                          {{ latestLogs[collection.id].method }}
                        </Badge>
                      </div>
                      <div>
                        <span class="text-muted-foreground">路径:</span>
                        <code class="ml-2 text-xs">{{ latestLogs[collection.id].path }}</code>
                      </div>
                      <div>
                        <span class="text-muted-foreground">状态:</span>
                        <Badge :variant="getStatusVariant(latestLogs[collection.id].status)" class="ml-2">
                          {{ latestLogs[collection.id].status }}
                        </Badge>
                      </div>
                      <div>
                        <span class="text-muted-foreground">耗时:</span>
                        <span class="ml-2">
                          {{ latestLogs[collection.id].duration !== null && latestLogs[collection.id].duration !== undefined ? `${latestLogs[collection.id].duration}ms` : '未返回' }}
                </span>
                        <Badge v-if="latestLogs[collection.id].from_cache" variant="secondary" class="ml-2">
                          缓存
                        </Badge>
                      </div>
                    </div>
                    <div v-if="expandedLogDetails[collection.id]" class="space-y-2 border-t pt-4">
                      <div v-if="latestLogs[collection.id].request_params">
                        <div class="text-sm font-medium mb-2">请求参数 (Params):</div>
                        <pre class="bg-muted p-3 rounded text-xs overflow-auto max-h-40">{{ formatJSON(latestLogs[collection.id].request_params) }}</pre>
                      </div>
                      <div v-if="latestLogs[collection.id].request_body">
                        <div class="text-sm font-medium mb-2">请求体 (Body):</div>
                        <pre class="bg-muted p-3 rounded text-xs overflow-auto max-h-40">{{ formatJSON(latestLogs[collection.id].request_body) }}</pre>
                      </div>
                      <div v-if="!latestLogs[collection.id].request_params && !latestLogs[collection.id].request_body" class="text-sm text-muted-foreground">
                        无请求参数
                      </div>
                    </div>
                  </div>
              </td>
              </tr>
            </template>
            <tr v-if="collections.length === 0">
              <td colspan="8" class="p-8 text-center text-muted-foreground">
                No collections found. Create one to get started.
              </td>
            </tr>
          </tbody>
        </Table>
      </CardContent>
    </Card>

    <!-- Create Collection Dialog -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center"
      @click.self="showCreateModal = false"
    >
      <div class="fixed inset-0 bg-background/80 backdrop-blur-sm" />
      <Card class="relative z-50 w-full max-w-3xl max-h-[90vh] overflow-hidden flex flex-col">
        <CardHeader>
          <CardTitle>Create New Collection</CardTitle>
        </CardHeader>
        <CardContent class="flex-1 overflow-y-auto">
          <form @submit.prevent="createCollection">
            <Tabs v-model="createTab" class="w-full">
              <TabsList class="grid w-full grid-cols-4">
                <TabsTrigger value="basic">Basic</TabsTrigger>
                <TabsTrigger value="health">Health</TabsTrigger>
                <TabsTrigger value="logging">Logging</TabsTrigger>
                <TabsTrigger value="caching">Caching</TabsTrigger>
              </TabsList>
              
              <TabsContent value="basic" class="space-y-4 mt-4">
                <div class="space-y-2">
                  <label class="text-sm font-medium">Name *</label>
                  <Input 
                    v-model="newCollection.name" 
                    required 
                    placeholder="My API Collection"
                    @input="handleNameChange"
                  />
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Description</label>
                  <textarea
                    v-model="newCollection.description"
                    class="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    placeholder="Collection description"
                  />
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Prefix *</label>
                  <div class="flex items-center space-x-2">
                    <Input 
                      v-model="newCollection.prefix" 
                      required 
                      placeholder="api-v1"
                      :class="prefixError ? 'border-red-500' : ''"
                      @input="() => { prefixManuallyEdited = true; checkPrefix(); }"
                    />
                    <Button 
                      type="button" 
                      variant="outline" 
                      size="sm"
                      @click="generatePrefixFromName"
                      title="根据名称自动生成"
                    >
                      自动生成
                    </Button>
                  </div>
                  <p v-if="prefixError" class="text-xs text-red-500">{{ prefixError }}</p>
                  <p v-else class="text-xs text-muted-foreground">External gateway prefix (e.g., api-v1). 不能包含 '/' 字符，将自动替换为 '-'</p>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Base URL *</label>
                  <Input v-model="newCollection.base_url" required placeholder="http://localhost:3000" />
                  <p class="text-xs text-muted-foreground">Target service base URL</p>
                </div>
              </TabsContent>
              
              <TabsContent value="health" class="space-y-4 mt-4">
                <div class="space-y-2">
                  <label class="text-sm font-medium">Health Check Path</label>
                  <Input v-model="newCollection.health_path" placeholder="/health" />
                  <p class="text-xs text-muted-foreground">Path to check service health status</p>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Health Check Interval (seconds)</label>
                  <Input v-model.number="newCollection.health_interval" type="number" min="10" />
                  <p class="text-xs text-muted-foreground">How often to check health status</p>
                </div>
              </TabsContent>
              
              <TabsContent value="logging" class="space-y-4 mt-4">
                <div class="flex items-center space-x-2">
                  <input
                    v-model="newCollection.log_enabled"
                    type="checkbox"
                    id="log_enabled"
                    class="h-4 w-4 rounded border-gray-300"
                  />
                  <label for="log_enabled" class="text-sm font-medium">Enable Logging</label>
                </div>
                <div class="flex items-center space-x-2">
                  <input
                    v-model="newCollection.log_rolling"
                    type="checkbox"
                    id="log_rolling"
                    class="h-4 w-4 rounded border-gray-300"
                  />
                  <label for="log_rolling" class="text-sm font-medium">Enable Rolling Logs</label>
                  <p class="text-xs text-muted-foreground">Automatically remove old logs when limit is reached</p>
                </div>
                <div class="space-y-2">
                  <label class="text-sm font-medium">Max Log Entries</label>
                  <Input v-model.number="newCollection.log_max_entries" type="number" min="100" />
                  <p class="text-xs text-muted-foreground">Maximum number of log entries to keep</p>
            </div>
              </TabsContent>
              
              <TabsContent value="caching" class="space-y-4 mt-4">
                <div class="flex items-center space-x-2">
                  <input
                    v-model="newCollection.cache_enabled"
                    type="checkbox"
                    id="cache_enabled"
                    class="h-4 w-4 rounded border-gray-300"
                  />
                  <label for="cache_enabled" class="text-sm font-medium">Enable Caching</label>
            </div>
                <div v-if="newCollection.cache_enabled" class="space-y-2">
                  <label class="text-sm font-medium">Cache TTL (seconds)</label>
                  <Input v-model.number="newCollection.cache_ttl" type="number" min="60" />
                  <p class="text-xs text-muted-foreground">Time to live for cached responses</p>
            </div>
                <div v-if="newCollection.cache_enabled" class="space-y-2">
                  <label class="text-sm font-medium">Cache Key Strategy</label>
                  <select v-model="newCollection.cache_key_strategy" class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
                    <option value="params">Query Parameters Only</option>
                    <option value="body">Request Body Only</option>
                    <option value="all">All (Params + Body)</option>
                  </select>
                  <p class="text-xs text-muted-foreground">How to generate cache keys</p>
            </div>
              </TabsContent>
            </Tabs>
            
            <div class="flex justify-end space-x-2 mt-6 pt-4 border-t">
              <Button type="button" variant="outline" @click="showCreateModal = false">Cancel</Button>
              <Button type="submit">Create</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
    
    <!-- Confirm Dialog -->
    <AlertDialog
      v-model="showDialog"
      :title="dialogConfig.title"
      :description="dialogConfig.description"
      :confirm-text="dialogConfig.confirmText"
      :cancel-text="dialogConfig.cancelText"
      :confirm-variant="dialogConfig.confirmVariant"
      @confirm="handleConfirm"
      @cancel="handleCancel"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue"
import axios from "axios"
import { pinyin } from "pinyin-pro"
import { Plus, Eye, Pause, Play, Trash2, Copy, ChevronDown } from "lucide-vue-next"
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
import AlertDialog from "@/components/ui/alert-dialog.vue"
import { useToast } from "@/composables/useToast"
import { useConfirm } from "@/composables/useConfirm"

const collections = ref([])
const showCreateModal = ref(false)
const createTab = ref("basic")
const expandedRows = ref({})
const expandedLogDetails = ref({})
const latestLogs = ref({})
const { toast } = useToast()
const { showDialog, dialogConfig, confirm, handleConfirm, handleCancel } = useConfirm()

const newCollection = ref({
  name: "",
  description: "",
  prefix: "",
  base_url: "",
  health_path: "",
  health_interval: 30,
  log_enabled: true,
  log_rolling: true,
  log_max_entries: 1000,
  cache_enabled: false,
  cache_ttl: 300,
  cache_key_strategy: "all",
})

const prefixError = ref("")
const prefixChecking = ref(false)
const prefixManuallyEdited = ref(false)

// 将字符串转换为拼音并生成 Prefix
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

// 生成 Prefix
const generatePrefixFromName = () => {
  if (!newCollection.value.name) {
    toast.error("请先输入名称", "名称不能为空")
    return
  }
  
  let prefix = toPinyin(newCollection.value.name)
  
  // 确保前缀不为空
  if (!prefix) {
    prefix = "collection"
  }
  
  newCollection.value.prefix = prefix
  prefixManuallyEdited.value = false
  checkPrefix()
}

// 处理名称变化
const handleNameChange = () => {
  // 如果 prefix 为空或未被手动编辑，自动生成
  if (!newCollection.value.prefix || !prefixManuallyEdited.value) {
    generatePrefixFromName()
  }
}

// 检查 Prefix 唯一性
const checkPrefix = async () => {
  let prefix = newCollection.value.prefix
  
  // 清理前缀：替换 / 为 -
  if (prefix.includes('/')) {
    prefix = prefix.replace(/\//g, '-')
    newCollection.value.prefix = prefix
  }
  
  if (!prefix || prefix.trim() === '') {
    prefixError.value = ""
    return
  }
  
  // 检查是否包含斜杠
  if (prefix.includes('/')) {
    prefixError.value = "Prefix 不能包含 '/' 字符，已自动替换为 '-'"
    return
  }
  
  prefixChecking.value = true
  try {
    const response = await axios.get(`/api/collections/check-prefix/${encodeURIComponent(prefix)}`)
    if (response.data.exists) {
      prefixError.value = "该 Prefix 已存在，请使用其他值"
    } else {
      prefixError.value = ""
    }
  } catch (error) {
    console.error("Failed to check prefix:", error)
    prefixError.value = ""
  } finally {
    prefixChecking.value = false
  }
}

const fetchCollections = async () => {
  try {
    const response = await axios.get("/api/collections")
    collections.value = response.data
    // Fetch latest logs for all collections
    for (const collection of collections.value) {
      await fetchLatestLog(collection.id)
    }
  } catch (error) {
    console.error("Failed to fetch collections:", error)
    toast.error("Failed to fetch collections", error.message)
  }
}

const fetchLatestLog = async (collectionId) => {
  try {
    const response = await axios.get(`/api/logs/${collectionId}/latest`)
    latestLogs.value[collectionId] = response.data
  } catch (error) {
    latestLogs.value[collectionId] = null
  }
}

const toggleExpand = async (collectionId) => {
  if (!expandedRows.value[collectionId]) {
    expandedRows.value[collectionId] = true
    await fetchLatestLog(collectionId)
  } else {
    expandedRows.value[collectionId] = false
    expandedLogDetails.value[collectionId] = false
  }
}

const toggleLogDetail = (collectionId) => {
  expandedLogDetails.value[collectionId] = !expandedLogDetails.value[collectionId]
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

const createCollection = async (e) => {
  if (e) {
    e.preventDefault()
  }
  
  // 验证 prefix
  if (prefixError.value) {
    toast.error("Prefix 验证失败", prefixError.value)
    return
  }
  
  // 确保 prefix 不包含斜杠
  if (newCollection.value.prefix.includes('/')) {
    newCollection.value.prefix = newCollection.value.prefix.replace(/\//g, '-')
  }
  
  try {
    console.log(newCollection.value)
    // Ensure all required fields are present
    const payload = {
      name: newCollection.value.name,
      description: newCollection.value.description || "",
      prefix: newCollection.value.prefix,
      base_url: newCollection.value.base_url,
      health_path: newCollection.value.health_path || "",
      health_interval: newCollection.value.health_interval || 30,
      log_enabled: newCollection.value.log_enabled !== undefined ? newCollection.value.log_enabled : true,
      log_rolling: newCollection.value.log_rolling !== undefined ? newCollection.value.log_rolling : true,
      log_max_entries: newCollection.value.log_max_entries || 1000,
      cache_enabled: newCollection.value.cache_enabled !== undefined ? newCollection.value.cache_enabled : false,
      cache_ttl: newCollection.value.cache_ttl || 300,
      cache_key_strategy: newCollection.value.cache_key_strategy || "all",
    }
    const response = await axios.post("/api/collections", payload)
    collections.value.push(response.data)
    showCreateModal.value = false
    toast.success("Collection created", "Collection has been created successfully")
    prefixError.value = ""
    prefixManuallyEdited.value = false
    newCollection.value = {
      name: "",
      description: "",
      prefix: "",
      base_url: "",
      health_path: "",
      health_interval: 30,
      log_enabled: true,
      log_rolling: true,
      log_max_entries: 1000,
      cache_enabled: false,
      cache_ttl: 300,
      cache_key_strategy: "all",
    }
  } catch (error) {
    console.error("Failed to create collection:", error)
    if (error.response && error.response.status === 409) {
      prefixError.value = error.response.data.error || "Prefix 已存在"
      toast.error("创建失败", error.response.data.error || "Prefix 已存在")
    } else {
      toast.error("Failed to create collection", error.message)
    }
  }
}

const toggleCollection = async (collection) => {
  try {
    await axios.post(`/api/collections/${collection.id}/toggle`)
    await fetchCollections()
    toast.success(
      "Collection updated",
      `Collection has been ${collection.active ? "deactivated" : "activated"}`
    )
  } catch (error) {
    console.error("Failed to toggle collection:", error)
    toast.error("Failed to toggle collection", error.message)
  }
}

const deleteCollection = async (id) => {
  const confirmed = await confirm({
    title: "Delete Collection",
    description: "Are you sure you want to delete this collection? This action cannot be undone.",
    confirmText: "Delete",
    cancelText: "Cancel",
    confirmVariant: "destructive",
  })
  
  if (confirmed) {
    try {
      await axios.delete(`/api/collections/${id}`)
      collections.value = collections.value.filter((c) => c.id !== id)
      toast.success("Collection deleted", "Collection has been deleted successfully")
    } catch (error) {
      console.error("Failed to delete collection:", error)
      toast.error("Failed to delete collection", error.message)
    }
  }
}

const copyProxyUrl = (prefix) => {
  const url = `${window.location.origin}/proxy/${prefix}/*`
  navigator.clipboard.writeText(url)
  toast.success("Copied", "Proxy URL copied to clipboard")
}

onMounted(() => {
  fetchCollections()
})
</script>

