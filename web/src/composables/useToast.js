export function useToast() {
  const toast = {
    success: (title, description) => {
      if (typeof window !== "undefined" && window.toast) {
        window.toast.success(title, description)
      } else {
        console.log("Toast:", title, description)
      }
    },
    error: (title, description) => {
      if (typeof window !== "undefined" && window.toast) {
        window.toast.error(title, description)
      } else {
        console.error("Error:", title, description)
      }
    },
    warning: (title, description) => {
      if (typeof window !== "undefined" && window.toast) {
        window.toast.warning(title, description)
      } else {
        console.warn("Warning:", title, description)
      }
    },
    info: (title, description) => {
      if (typeof window !== "undefined" && window.toast) {
        window.toast.info(title, description)
      } else {
        console.info("Info:", title, description)
      }
    },
    show: (options) => {
      if (typeof window !== "undefined" && window.toast) {
        window.toast.show(options)
      }
    },
  }
  
  return { toast }
}

