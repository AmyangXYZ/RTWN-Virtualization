import { ref, watch, onMounted } from 'vue'
import { network } from './useNetwork.js'

import * as echarts from 'echarts/core'
import { GraphChart } from 'echarts/charts'
import { CanvasRenderer } from 'echarts/renderers'

import { chartTheme } from './useConfig.js'

echarts.use([GraphChart, CanvasRenderer])

export function useDrawCriticalNodes(chartDom) {
  const colors = [
    '#5470c6',
    '#91cc75',
    '#fac858',
    '#ee6666',
    '#73c0de',
    '#3ba272',
    '#fc8452',
    '#9a60b4',
    '#ea7ccc'
  ]

  const nodes = ref([])

  let option = {
    series: [
      {
        type: 'graph',
        layout: 'force',
        force: {
          repulsion: 15
        },
        nodes: [],
        edges: [],
        draggable: true,
        symbolSize: 12,
        zoom: 1,
        emphasis: {
          focus: 'none'
        },
        lineStyle: {
          width: 1.2,
          color: 'black',
          opacity: 0.5
        },
        label: {
          show: true,
          fontSize: 10,
          color: 'white',
          fontWeight: 600
        },
        roam: true,
        animation: false
      }
    ]
  }

  let chart = {}
  function draw() {
    nodes.value = network.value.manager.critical_nodes
    option.series[0].nodes = []
    option.series[0].edges = []
    for (let n in nodes.value) {
      let apps = nodes.value[n]
      option.series[0].nodes.push({
        name: n + '',
        symbol: 'rect',
        symbolSize: 14,
        itemStyle: {
          color: 'grey'
        }
      })
      for (let j in apps) {
        let appID = apps[j]
        option.series[0].nodes.push({
          name: n + '-' + appID,
          value: appID,
          itemStyle: {
            color: colors[appID]
          },
          label: {
            formatter: (item) => {
              return item.value
            }
          }
        })

        option.series[0].edges.push({
          source: n + '-' + appID,
          target: n + ''
        })
      }
    }
    chart.setOption(option)
  }

  onMounted(() => {
    chart = echarts.init(chartDom.value, chartTheme)
    draw()
  })

  watch(network, () => {
    draw()
  })
  return { nodes }
}
