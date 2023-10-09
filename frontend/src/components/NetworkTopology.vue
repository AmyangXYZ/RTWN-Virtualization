<script setup>
import { ref } from 'vue'
import { useDrawTopology } from '../hooks/useDrawTopology.js'
import { useSeed } from '../hooks/useSeed.js'

const { seed, updateSeed } = useSeed('topo')
const chartDom = ref(null)
useDrawTopology(chartDom)
</script>

<template>
  <vs-card>
    <template #header>
      <vs-row vs-align="center" vs-w="12">
        <vs-col vs-w="2">
          <h3>Topology</h3>
        </vs-col>
        <vs-col vs-offset="6" vs-w="1">
          <span style="font-size: 0.9rem">Seed: </span>
        </vs-col>
        <vs-col vs-w="2">
          <vs-input-number size="small" color="primary" v-model="seed"></vs-input-number>
        </vs-col>
        <vs-col vs-offset="0.2" vs-w="0.8">
          <vs-button
            color="primary"
            type="relief"
            size="small"
            style="padding: 3px; line-height: 1"
            @click="updateSeed"
            ><vs-icon size="1.2rem" icon="refresh"></vs-icon
          ></vs-button>
        </vs-col>
      </vs-row>
    </template>
    <div class="chart" ref="chartDom"></div>
  </vs-card>
</template>

<style scoped>
.chart {
  width: 100%;
  height: 350px;
}
</style>
