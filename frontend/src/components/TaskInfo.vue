<script setup>
// import {ref} from 'vue'
import { network } from '../hooks/useNetwork.js'
import { selectedAppID } from '../hooks/useSelectedApp.js'
// const selectedAppID = ref(0)
</script>

<template>
  <vs-card v-if="selectedAppID>-1">
    <template #header>
      <vs-row vs-align="center" vs-w="12">
        <h3>Tasks in App {{ selectedAppID }}</h3>
      </vs-row>
    </template>

    <vs-table :data="network.apps[selectedAppID].tasks" stripe>
      <!-- <template #header>
        <h3>Tasks of App {{ selectedAppID }}</h3>
      </template> -->
      <template #thead>
        <vs-th>
          <p>ID</p>
        </vs-th>
        <vs-th>
          <p>Deadline</p>
        </vs-th>
        <vs-th>
          <p>Period</p>
        </vs-th>
        <vs-th>
          <p>Path</p>
        </vs-th>
      </template>

      <template v-slot="{ data }">
        <vs-tr :data="tr" :key="indextr" v-for="(tr, indextr) in data">
          <vs-td :data="tr.id">
            <p>{{ tr.id }}</p>
          </vs-td>

          <vs-td :data="tr.period">
            <p>{{ tr.period }}</p>
          </vs-td>

          <vs-td :data="tr.period">
            <p>{{ tr.period }}</p>
          </vs-td>

          <vs-td :data="tr.path">
            <p>{{ tr.path }}</p>
          </vs-td>
        </vs-tr>
      </template>
    </vs-table>
  </vs-card>
</template>

<style scoped>
p {
  font-size: 0.78rem;
  font-weight: 600;
}

.vs-table--tbody-table .tr-values td {
  padding-top: 2px;
  padding-bottom: 2px;
}
</style>
