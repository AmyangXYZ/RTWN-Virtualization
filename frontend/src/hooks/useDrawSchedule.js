import { ref, watch, onMounted } from 'vue'
import { network } from '../hooks/useNetwork.js'

import * as echarts from 'echarts/core'
import { HeatmapChart } from 'echarts/charts'
import { GridComponent, ToolboxComponent, VisualMapComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

import { chartTheme } from './useConfig.js'

echarts.use([HeatmapChart, ToolboxComponent, VisualMapComponent, GridComponent, CanvasRenderer])

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

const CELL_HEIGHT = 30

export function useDrawSchedule(resource, chartDom) {
  const partition = ref(null)

  let option = {
    toolbox: {
      feature: {
        saveAsImage: {
          pixelRatio: 8
        }
      }
    },
    tooltip: {
      trigger: "axis"
    },
    grid: [
      {
        top: '40px',
        bottom: '0px',
        // height:"60%",
        left: '42px',
        right: '8px'
      }
    ],
    xAxis: [
      {
        name: 'Time',
        nameLocation: 'center',
        nameGap: 24,
        nameTextStyle: {
          fontSize: 13,
          fontWeight: '600'
        },
        position: 'top',
        type: 'category',
        data: [],
        splitArea: {
          show: true
        },
        axisLabel: {
          show: true,
          fontSize: 12,
          interval: 0,
        },
        gridIndex: 0
      }
    ],
    yAxis: [
      {
        name: resource == 'node' ? 'Node' : 'Channel',
        nameLocation: 'center',
        nameGap: 25,
        nameTextStyle: {
          fontSize: 13,
          fontWeight: '600'
        },
        type: 'category',
        data: [],
        splitArea: {
          show: true
        },
        axisLabel: {
          interval: 0,
          fontSize: 12
        },
        inverse: true,
        gridIndex: 0
      }
    ],
    visualMap: {
      show: true,
      type: 'piecewise',
      pieces: [],
      min: 0,
      max: 10,
      calculable: true,
      orient: 'horizontal',
      right: '1%',
      top: '0px',
      itemWidth: 15,
      itemHeight: 12,
      width: "50%",
      textStyle: {
        fontSize: 12
      },
      formatter: (item) => {
        return `App ${item + 1}`
      }
    },
    series: []
  }

  let chart = {}

  function draw() {
    option.xAxis[0].data = []
    option.yAxis[0].data = []
    option.series = []

    const settings = network.value.settings

    for (let a = 0; a < settings.num_apps; a++) {
      option.visualMap.pieces.push({ value: a, color: colors[a] })
      option.series.push({
        type: 'heatmap',
        data: [],
        label: {
          show: false,
          fontSize: 13,
          formatter: (item) => {
            return item.name
          }
        },
        itemStyle: {
          borderColor: 'whitesmoke',
          borderWidth: 1
        },
        emphasis: {
          itemStyle: {
            borderWidth: 0,
            shadowBlur: 10,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        animation: false
      })
    }

    partition.value = network.value.manager.partition_bandwidth
    chartDom.value.style.height = settings.num_channels * CELL_HEIGHT + 40 + 'px'

    for (let x = 0; x < partition.value.length; x++) {
      option.xAxis[0].data.push(x)
    }
    for (let y = 1; y < partition.value[0].length; y++) {
      option.yAxis[0].data.push(y)
    }
    var cells = []
    for (let slot = 1; slot <= settings.num_slots; slot++) {
      let ch_app = { 0: 0, 1: 0, 2: 0 }
      for (let ch = 0; ch < partition.value[slot].length; ch++) {
        let cell = partition.value[slot][ch]
        if (cell.assigned) {
          // find start ch
          ch_app[cell.app_id]++
          const sch_cell = network.value.apps[cell.app_id].schedule[1][slot][ch_app[cell.app_id]]
          if (sch_cell.assigned) {
            cells.push({
              slot: {
                slot_offset: slot + 5,
                channel_offset: ch ,
              },
              subslot: {
                offset: 0, period: 1,
              },
              sender: sch_cell.sender,
              receiver: sch_cell.receiver,
            })
            option.series[cell.app_id].data.push({
              value: [slot + 5, ch - 1, cell.app_id],
              name: `${sch_cell.sender}\n${sch_cell.receiver}`,
              label: {
                show: sch_cell.assigned,
                fontSize: 11,
                // show:hasTx,
                // formatter: (item) => {
                //   return item.name[2]
                // }
              }
            })
          }
        }
      }
    }
    console.log(cells)



    chart.setOption(option)
    chart.resize()
  }

  onMounted(() => {
    chart = echarts.init(chartDom.value, chartTheme)
    draw()
  })

  watch(network, () => {
    draw()
  })

  return { partition }
}
