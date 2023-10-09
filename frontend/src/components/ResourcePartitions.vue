<script setup>
import { ref } from 'vue'
import { useDrawPartition } from '../hooks/useDrawPartition.js'

const chartDomNode = ref(null)
const chartDomChannel = ref(null)
useDrawPartition('node', chartDomNode)
useDrawPartition('channel', chartDomChannel)
</script>

<template>
  <vs-card>
    <template #header>
      <h3>Partitions</h3>
    </template>
    <div class="chart" ref="chartDomChannel" style="margin-bottom: 8px"></div>
    <div class="chart" ref="chartDomNode"></div>
    
  </vs-card>
</template>

<style scoped>
.chart {
  width: 100%;
  height: 1px;
}
</style>
