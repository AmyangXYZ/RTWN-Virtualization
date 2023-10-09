<script setup>
import { ref } from 'vue'
import { useDrawCriticalNodes } from '../hooks/useDrawCriticalNodes.js'
import { network } from '../hooks/useNetwork.js'

const chartDom = ref(null)
const { nodes } = useDrawCriticalNodes(chartDom)
</script>

<template>
  <vs-card>
    <template #header>
      <h3>Interfaces</h3>
    </template>

    <vs-row vs-align="center">
      <vs-col vs-w="3.2">
        <div class="chart" ref="chartDom"></div>
      </vs-col>
      <vs-col vs-w="8.8">
        <vs-table :data="Object.keys(nodes)" stripe>
          <template #thead>
            <vs-th>
              <!-- <p>Node</p> -->
            </vs-th>
            <vs-th v-for="(app, i) in network.apps" :key="i">
              <p>App {{ app.id }}</p>
            </vs-th>
          </template>

          <template v-slot="{ data }">
            <vs-tr>
              <vs-td>
                <p style="color: #bf0000">Channel</p>
              </vs-td>
              <vs-td v-for="(app, ai) in network.apps" :key="ai">
                <p style="color: #bf0000">
                  ({{ app.bandwidth_interface.alpha }},{{ app.bandwidth_interface.r.toFixed(2) }})
                </p>
              </vs-td>
            </vs-tr>

            <vs-tr :data="node" :key="node" v-for="node in data">
              <vs-td>
                <p>{{ 'v_' + node }}</p>
              </vs-td>
              <vs-td v-for="(app, ai) in network.apps" :key="ai">
                <p v-if="app.interfaces[node] != null">
                  ({{ app.interfaces[node].alpha }},{{ app.interfaces[node].r.toFixed(2) }})
                </p>
                <p v-else>-</p>
              </vs-td>
            </vs-tr>
          </template>
        </vs-table>
      </vs-col>
    </vs-row>
  </vs-card>
</template>

<style scoped>
.chart {
  width: 100%;
  height: 220px;
}

.vs-table--tbody-table .tr-values td {
  padding-top: 2px;
  padding-bottom: 2px;
}

p {
  font-size: 0.78rem;
  font-weight: 600;
}
</style>
