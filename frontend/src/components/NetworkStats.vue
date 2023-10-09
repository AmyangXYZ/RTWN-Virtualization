<script setup>
import { onMounted, nextTick, watch } from 'vue'
import { useDrawSupplyFunc } from '../hooks/useDrawSupplyFunc.js'
import { network } from '../hooks/useNetwork.js'

onMounted(() => {
  for (let app of network.value.apps) {
    useDrawSupplyFunc(app.id, -1)
    for (let node in app.critical_nodes) {
      useDrawSupplyFunc(app.id, node)
    }
  }
})

watch(network, async () => {
  await nextTick()
  for (let app of network.value.apps) {
    useDrawSupplyFunc(app.id, -1)
    for (let node in app.critical_nodes) {
      useDrawSupplyFunc(app.id, node)
    }
  }
})
</script>

<template>
  <vs-card>
    <template #header>
      <h3>Statistics (Success: {{ network.stat.task_success }})</h3>
    </template>
    <vs-row :key="app" v-for="app of network.apps">
      <vs-col vs-offset="0" vs-w="1.7">
        {{ `P_${app.id}^Ch` }}
        <div class="chart" :id="`supply-${app.id}--1`"></div>
      </vs-col>
      <vs-col vs-offset="0" vs-w="1.7" :key="idx" v-for="(idx, node) in app.critical_nodes">
        {{ `P_${app.id}^${node}` }}
        <div class="chart" :id="`supply-${app.id}-${node}`"></div>
      </vs-col>
    </vs-row>
  </vs-card>
</template>

<style scoped>
.chart {
  width: 100%;
  height: 80px;
}
</style>
