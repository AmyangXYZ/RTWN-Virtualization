<script setup>
import { ref } from 'vue'
import { useDrawNodeDeps } from '../hooks/useDrawNodeDeps.js'
import { selectedAppID } from '../hooks/useSelectedApp.js'

const chartDom = ref(null)
const showCriticalOnly = ref(true)
const { deps } = useDrawNodeDeps(chartDom, showCriticalOnly)
</script>

<template>
  <vs-card>
    <template #header>
      <h3>Node Dep Graph in App {{ selectedAppID == -1 ? 0 : selectedAppID }}</h3>
    </template>
    <!-- <vs-row vs-align="center" vs-type="flex" vs-justify="end">
      <vs-col vs-w="2.7">
        <p>Critical node only</p>
      </vs-col>
      <vs-col vs-w="1">
        <vs-switch v-model="showCriticalOnly" />
      </vs-col>
    </vs-row> -->

    <vs-table :data="deps" stripe v-if="false">
      <template #header>
        <h3>Raw dependencies</h3>
      </template>
      <template #thead>
        <vs-th>
          <p>TaskID</p>
        </vs-th>
        <vs-th>
          <p>V1</p>
        </vs-th>
        <vs-th>
          <p>V2</p>
        </vs-th>
        <vs-th>
          <p>Next</p>
        </vs-th>
        <vs-th>
          <p>Period</p>
        </vs-th>
      </template>

      <template v-slot="{ data }">
        <vs-tr :data="tr" :key="i" v-for="(tr, i) in data">
          <vs-td>
            <p>{{ tr.task_id }}</p>
          </vs-td>
          <vs-td>
            <p>{{ tr.sender }}</p>
          </vs-td>
          <vs-td>
            <p>{{ tr.receiver }}</p>
          </vs-td>
          <vs-td>
            <!-- <p>{{ tr.next }}</p> -->
          </vs-td>
          <vs-td>
            <p>{{ tr.period }}</p>
          </vs-td>
        </vs-tr>
      </template>
    </vs-table>

    <div class="chart" ref="chartDom"></div>
  </vs-card>
</template>

<style scoped>
.chart {
  width: 100%;
  height: 240px;
}
p {
  font-size: 0.85rem;
  font-weight: 600;
}

.vs-table--tbody-table .tr-values td {
  padding-top: 2px;
  padding-bottom: 2px;
}
</style>
