<script setup>
import Apps from '../components/AppsInfo.vue'
import Tasks from '../components/TaskInfo.vue'
import NodeDependency from '../components/NodeDependency.vue'
import ResourceInterfaces from '../components/ResourceInterfaces.vue'
import ResourcePartitions from '../components/ResourcePartitions.vue'
import AppSchedule from '../components/AppSchedule.vue'
import Settings from '../components/NetworkSettings.vue'
import Statistics from '../components/NetworkStats.vue'
import NetworkTopology from '../components/NetworkTopology.vue'
</script>

<template>
  <vs-row vs-w="12">
    <vs-col vs-offset="0.5" vs-w="4.2">
      <Settings />
      <NetworkTopology />
      <Apps />
      <Tasks />
      <NodeDependency />
    </vs-col>
    <vs-col vs-offset="0.2" vs-w="6.6">
      <ResourceInterfaces />
      <ResourcePartitions />
      
      <Statistics />
    </vs-col>
  </vs-row>
  <AppSchedule v-if="false"/>
</template>
