import { ref } from 'vue'
import { network, fetchNetwork } from '../hooks/useNetwork.js'

export function useSeed(flag) {
  const seed = ref(0)
  if (flag == 'topo') seed.value = network.value.settings.seed_topo
  else seed.value = network.value.settings.seed_app

  function updateSeed() {
    if (
      (flag == 'topo' && seed.value != network.value.settings.seed_topo) ||
      (flag == 'app' && seed.value != network.value.settings.seed_app)
    ) {
      fetch(`http://localhost:8000/api/new/${flag}/${seed.value}`, { method: 'put' }).then(() => {
        fetchNetwork()
      })
    }
  }
  return { seed, updateSeed }
}
