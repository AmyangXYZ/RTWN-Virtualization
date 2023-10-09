import { shallowRef } from 'vue'

export const network = shallowRef(null)

export async function fetchNetwork() {
  const resp = await (await fetch('http://localhost:8000/api/network/0')).json()
  network.value = resp.data
}
