<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold tracking-tight">Request Logs</h2>
        <p class="text-muted-foreground">View and manage request logs</p>
      </div>
      <Button variant="destructive" @click="clearLogs">
        <Trash2 class="mr-2 h-4 w-4" />
        Clear Logs
      </Button>
    </div>

    <Card>
      <CardContent class="p-0">
        <Table>
          <thead>
            <tr class="border-b">
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Timestamp</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Method</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Path</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Status</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Duration</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Size</th>
              <th class="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Client IP</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="log in requestLogs"
              :key="log.id"
              class="border-b transition-colors hover:bg-muted/50"
            >
              <td class="p-4 align-middle text-sm">{{ formatTimestamp(log.timestamp) }}</td>
              <td class="p-4 align-middle">
                <Badge :variant="getMethodVariant(log.method)" variant="outline">
                  {{ log.method }}
                </Badge>
              </td>
              <td class="p-4 align-middle font-mono text-sm">{{ log.path }}</td>
              <td class="p-4 align-middle">
                <Badge :variant="getStatusVariant(log.status)">
                  {{ log.status }}
                </Badge>
              </td>
              <td class="p-4 align-middle text-sm">{{ log.duration }}ms</td>
              <td class="p-4 align-middle text-sm">{{ formatBytes(log.response_size) }}</td>
              <td class="p-4 align-middle text-sm text-muted-foreground">{{ log.client_ip }}</td>
            </tr>
            <tr v-if="requestLogs.length === 0">
              <td colspan="7" class="p-8 text-center text-muted-foreground">
                No logs available
              </td>
            </tr>
          </tbody>
        </Table>
      </CardContent>
    </Card>
    
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
import { ref, onMounted, onUnmounted } from "vue"
import axios from "axios"
import { Trash2 } from "lucide-vue-next"
import Card from "@/components/ui/card.vue"
import CardContent from "@/components/ui/card-content.vue"
import Button from "@/components/ui/button.vue"
import Table from "@/components/ui/table.vue"
import Badge from "@/components/ui/badge.vue"
import AlertDialog from "@/components/ui/alert-dialog.vue"
import { useToast } from "@/composables/useToast"
import { useConfirm } from "@/composables/useConfirm"

const requestLogs = ref([])
const loading = ref(false)
let refreshInterval = null
const { toast } = useToast()
const { showDialog, dialogConfig, confirm, handleConfirm, handleCancel } = useConfirm()

const formatTimestamp = (timestamp) => {
  if (!timestamp) return "-"
  const date = new Date(timestamp)
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

const getStatusVariant = (status) => {
  if (status >= 200 && status < 300) return "default"
  if (status >= 400 && status < 500) return "secondary"
  return "destructive"
  }

const fetchLogs = async () => {
  try {
    loading.value = true
    const response = await axios.get("/api/logs?limit=100")
    requestLogs.value = response.data
  } catch (error) {
    console.error("Failed to fetch logs:", error)
  } finally {
    loading.value = false
  }
}

const clearLogs = async () => {
  const confirmed = await confirm({
    title: "Clear All Logs",
    description: "Are you sure you want to clear all logs? This action cannot be undone.",
    confirmText: "Clear",
    cancelText: "Cancel",
    confirmVariant: "destructive",
  })
  
  if (confirmed) {
    try {
      await axios.delete("/api/logs")
      requestLogs.value = []
      toast.success("Logs cleared", "All logs have been cleared successfully")
    } catch (error) {
      console.error("Failed to clear logs:", error)
      toast.error("Failed to clear logs", error.message)
    }
  }
}

onMounted(() => {
  fetchLogs()
  refreshInterval = setInterval(fetchLogs, 5000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
