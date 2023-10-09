<script setup>
import { ref } from 'vue'
import { useDrawSchedule } from '../hooks/useDrawSchedule.js'

const chartDomChannel = ref(null)
useDrawSchedule('channel', chartDomChannel)
</script>

<template>
  <vs-card>
    <template #header>
      <h3>Schedule</h3>
    </template>
    <div class="chart" ref="chartDomChannel"></div>
    
  </vs-card>
</template>

<style scoped>
.chart {
  width: 100%;
  height: 1px;
}
</style>
