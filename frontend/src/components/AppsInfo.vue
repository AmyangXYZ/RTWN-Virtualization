<script setup>
import { network } from '../hooks/useNetwork.js'
import { selectedAppInTable, selectApp } from '../hooks/useSelectedApp.js'
import { useSeed } from '../hooks/useSeed.js'
const { seed, updateSeed } = useSeed('app')
</script>

<template>
  <vs-card>
    <template #header>
      <vs-row vs-align="center" vs-w="12">
        <vs-col vs-w="2">
          <h3>Applications</h3>
        </vs-col>
        <vs-col vs-offset="6" vs-w="1">
          <span style="font-size: 0.9rem">Seed: </span>
        </vs-col>
        <vs-col vs-w="2">
          <vs-input-number size="small" color="danger" v-model="seed"></vs-input-number>
        </vs-col>
        <vs-col vs-offset="0.2" vs-w="0.8">
          <vs-button
            color="danger"
            type="relief"
            size="small"
            style="padding: 3px; line-height: 1"
            @click="updateSeed"
            ><vs-icon size="1.2rem" icon="refresh"></vs-icon
          ></vs-button>
        </vs-col>
      </vs-row>
    </template>

    <vs-table
      :data="network.apps"
      stripe
      v-model="selectedAppInTable"
      @selected="(item) => selectApp(item.id)"
    >
      <template #thead>
        <vs-th>
          <p>ID</p>
        </vs-th>
        <vs-th>
          <p>Num_Channels</p>
        </vs-th>
        <vs-th>
          <p>Num_Nodes</p>
        </vs-th>
        <vs-th>
          <p>Num_Tasks</p>
        </vs-th>
      </template>

      <template v-slot="{ data }">
        <vs-tr :data="tr" :key="i" v-for="(tr, i) in data">
          <vs-td>
            <p>{{ tr.id }}</p>
          </vs-td>
          <vs-td>
            <p>{{ tr.num_channels }}</p>
          </vs-td>
          <vs-td>
            <p>{{ tr.used_nodes.length }}</p>
          </vs-td>
          <vs-td>
            <p>{{ tr.num_tasks }}</p>
          </vs-td>
        </vs-tr>
      </template>
    </vs-table>

    <!-- <vs-table
      style="margin-top: 5px"
      v-if="selectedAppID > -1"
      :data="network.apps[selectedAppID].tasks"
      stripe
    >
      <template #header>
        <h3>Tasks of App {{ selectedAppID }}</h3>
      </template>
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
    </vs-table> -->
  </vs-card>
</template>

<style scoped>
p {
  font-size: 0.8rem;
  font-weight: 600;
}

.vs-table--tbody-table .tr-values td {
  padding-top: 2px;
  padding-bottom: 2px;
}
</style>
