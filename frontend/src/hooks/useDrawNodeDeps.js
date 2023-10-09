import { ref, onMounted, watch } from 'vue'
import { network } from '../hooks/useNetwork.js'
import { selectedAppID } from './useSelectedApp.js'

import * as echarts from 'echarts/core'
import { GraphChart } from 'echarts/charts'
import { CanvasRenderer } from 'echarts/renderers'

import { chartTheme } from './useConfig.js'

echarts.use([GraphChart, CanvasRenderer])

export function useDrawNodeDeps(chartDom, showCriticalOnly) {
  const deps = ref([])

  let option = {
    series: [
      {
        type: 'graph',
        layout: 'force',
        force: {
          initialLayout: 'circular',
          repulsion: 80,
          edgeLength: 50
        },
        nodes: [],
        edges: [],
        symbolSize: 24,
        label: {
          show: true,
          fontSize: 11,
          color: 'black',
          fontWeight: 600
        },
        itemStyle: {
          color: 'white',
          borderColor:"black"
        },
        lineStyle: {
          width: 1.5,
          color: 'black',
          opacity: 0.8
        },
        symbol: 'rect',
        edgeSymbol: ['none', 'arrow'],
        edgeSymbolSize: 5,
        draggable: true,
        animation: false
      }
    ]
  }

  let chart = {}

  function draw() {
    option.series[0].nodes = []
    option.series[0].edges = []
    const id = selectedAppID.value == -1 ? 0 : selectedAppID.value
    const app = network.value.apps[id]
    const dags = app.pre_dep_graph
    deps.value = dags

    for (let i in dags) {
      const dag = dags[i]
      for (let j in dag.vertices) {
        let v = dag.vertices[j]
        let node = {
          name: `${i}-${j}`, // unique name
          value: `${v[0]}\n${v[1]}`,
          label: {
            formatter: (item) => {
              return item.value
            }
          }
        }
        if (showCriticalOnly.value) {
          node.value = `${app.critical_nodes[v[0]] ? v[0] : '-'}\n${
            app.critical_nodes[v[1]] ? v[1] : '-'
          }`
        }
        option.series[0].nodes.push(node)
        if (j < dag.vertices.length - 1) {
          option.series[0].edges.push({
            source: `${i}-${j}`,
            target: `${i}-${parseInt(j) + 1}`
          })
        }
      }
    }
    chart.setOption(option)
  }

  onMounted(() => {
    chart = echarts.init(chartDom.value, chartTheme)
    draw()
  })

  watch([network, selectedAppID, showCriticalOnly], () => {
    draw()
  })

  return { deps }
}
