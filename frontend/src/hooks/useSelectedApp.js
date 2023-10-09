import { ref } from 'vue'

export const selectedAppID = ref(-1)

export const selectedAppInTable = ref({})

export function selectApp(id) {
  if (id == selectedAppID.value) {
    selectedAppID.value = -1
    selectedAppInTable.value = {}
  } else {
    selectedAppID.value = id
  }
}
