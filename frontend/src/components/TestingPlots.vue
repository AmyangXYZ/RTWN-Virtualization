<script setup>
import { ref, reactive } from 'vue'
import { useDrawSupplyFunc } from '../hooks/useDrawSupplyFuncTest.js'
import { useDrawPSR } from '../hooks/useDrawPSR.js'
import { useDrawTSR } from '../hooks/useDrawTSR.js'
import { useDrawLatency } from '../hooks/useDrawLatency.js'
// import { useDrawFeasibleSets } from '../hooks/useDrawFeasibleSets.js'

const L = 20
const intf = reactive({ alpha: 0.2, r: 1.5 })
const actual = ref('0')
const chartDomSupply = ref(null)
useDrawSupplyFunc(intf, L, actual, chartDomSupply)

// const settings = reactive({
//   T: 10,
//   N: 4,
//   c: 2,
//   p: 5,
//   alpha: 0.4,
//   r: 1
// })
// const chartDomFS = ref(null)
// const { sc_size, fs_size, fp_size, unsafe } = useDrawFeasibleSets(settings, chartDomFS)

const chartDomPSR = ref(null)
useDrawPSR(chartDomPSR)

const chartDomTSR = ref(null)
useDrawTSR(chartDomTSR)

const chartDomLat = ref(null)
useDrawLatency(chartDomLat)
</script>

<template>
  <vs-card>
    <vs-row>
      <vs-col vs-w="3">
        <h3>
          Interface: alpha:
          <input type="number" v-model="intf.alpha" />
          regularity:
          <input type="number" v-model="intf.r" />
        </h3>
        <h3>Slot allocation: <input v-model="actual" /></h3>
        <div class="chart-supply" ref="chartDomSupply"></div>
      </vs-col>
    </vs-row>
    <vs-row v-if="false">
      <vs-col style="text-align: center" vs-w="3.5">
        <h3>
          T
          <input type="number" v-model="settings.T" />
          N
          <input type="number" v-model="settings.N" />
        </h3>
        <h3>
          c
          <input type="number" v-model="settings.c" />
          p
          <input type="number" v-model="settings.p" />
        </h3>
        <h3>
          alpha
          <input type="number" v-model="settings.alpha" />
          r
          <input type="number" v-model="settings.r" />
        </h3>
        <br />
        <h3 style="text-align: center">
          SC: {{ sc_size }}, FS: {{ fs_size }}, FP: {{ fp_size }}, Unsafe: {{ unsafe }}
        </h3>
        <div class="chart" ref="chartDomFS"></div>
      </vs-col>
    </vs-row>
    <vs-row>
      <vs-col vs-w="6">
        <div class="chart-sr" ref="chartDomPSR"></div>
        <div class="chart-sr" ref="chartDomTSR"></div>
        <div class="chart-lat" ref="chartDomLat"></div>
      </vs-col>
    </vs-row>
  </vs-card>
</template>

<style scoped>
.chart-supply {
  width: 100%;
  height: 220px;
}

.chart-sr {
  width: 50%;
  height: 270px;
}
.chart-lat {
  width: 35%;
  height: 500px;
}
</style>
