import { ref } from 'vue'

export function useDrawStats() {
  const stats = ref(null)
  return { stats }
}
