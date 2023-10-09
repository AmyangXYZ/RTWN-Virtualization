<script setup>
import { ref, watchEffect } from 'vue'
import { network } from '../hooks/useNetwork.js'

const settingItems = ref([])
watchEffect(() => {
  settingItems.value = [
    {
      value: network.value.settings.num_apps,
      icon: 'group',
      color: 'danger'
    },
    {
      value: network.value.settings.num_nodes,
      icon: 'memory',
      color: 'primary'
    },
    {
      value: network.value.settings.num_slots,
      icon: 'schedule',
      color: 'success'
    },
    {
      value: network.value.settings.num_channels,
      icon: 'wifi_channel',
      color: 'grey'
    }
  ]
})
</script>

<template>
  <vs-row vs-w="12" vs-type="flex" vs-justify="space-between">
    <vs-col :key="item" v-for="item in settingItems" vs-w="2.5">
      <vs-card class="it">
        <vs-row vs-align="center" vs-type="flex" vs-justify="space-between">
          <vs-col vs-offset="0.5" vs-w="5">
            <vs-icon class="icon" :color="item.color" :icon="item.icon" />
          </vs-col>
          <vs-col vs-w="4">
            <h3>
              {{ item.value }}
            </h3>
          </vs-col>
        </vs-row>
      </vs-card>
    </vs-col>
  </vs-row>
</template>

<style scoped>
.it {
  line-height: 0.2rem;
}
.icon {
  font-size: 2rem;
}
h3 {
  font-size: 1.1rem;
}
</style>
