import { ref } from "vue"

export function useConfirm() {
  const showDialog = ref(false)
  const dialogConfig = ref({
    title: "Are you sure?",
    description: "",
    confirmText: "Confirm",
    cancelText: "Cancel",
    confirmVariant: "default",
  })
  
  let resolvePromise = null
  
  const confirm = (config = {}) => {
    dialogConfig.value = {
      ...dialogConfig.value,
      ...config,
    }
    showDialog.value = true
    
    return new Promise((resolve) => {
      resolvePromise = resolve
    })
  }
  
  const handleConfirm = () => {
    if (resolvePromise) {
      resolvePromise(true)
      resolvePromise = null
    }
    showDialog.value = false
  }
  
  const handleCancel = () => {
    if (resolvePromise) {
      resolvePromise(false)
      resolvePromise = null
    }
    showDialog.value = false
  }
  
  return {
    showDialog,
    dialogConfig,
    confirm,
    handleConfirm,
    handleCancel,
  }
}

