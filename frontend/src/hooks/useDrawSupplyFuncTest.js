import { onMounted, watch } from 'vue'

import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  MarkLineComponent,
  ToolboxComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([
  LineChart,
  GridComponent,
  ToolboxComponent,
  TooltipComponent,
  LegendComponent,
  MarkLineComponent,
  CanvasRenderer
])

export function useDrawSupplyFunc(intf, L, actual, chartDom) {
  let option = {
    legend: {
      width: '100%',
      itemWidth:20,
      itemHeight:10,
      textStyle: { fontSize: 12 }
    },
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      top: '30px',
      bottom: '34px',
      right: '0',
      left: '38px'
    },
    toolbox: {
      feature: {
        saveAsImage: {
          pixelRatio: 8
        }
      }
    },
    xAxis: {
      name: 'Time',
      type: 'category',
      data: [],
      axisLabel: {
        fontSize: 13
      },
      nameLocation: 'center',
      nameGap: 20,
      nameTextStyle: {
        fontSize: 13
      }
    },
    yAxis: {
      nameLocation: 'center',
      name: 'Supply',
      type: 'value',
      nameGap: 25,
      min: 0,
      max: L,
      axisLabel: {
        fontSize: 13
      },
      nameTextStyle: {
        fontSize: 13
      }
    },
    series: [
      {
        name: 'Uniform',
        data: [],
        type: 'line',
        symbolSize: 4,
        lineStyle: {
          type: 'solid',
          width: 1.5
        },

        color: 'grey',
        animation: false
      },
      {
        name: 'S',
        data: [],
        type: 'line',
        lineStyle: {
          type: 'solid',
          width: 1.5
        },
        symbolSize: 4,
        symbol:"emptyCircle",
        color: 'blue',
        animation: false,
        z:10,
      },
      {
        name: 'Sᵤ',
        data: [],
        type: 'line',
        lineStyle: {
          type: 'solid',
          width: 1.5
        },
        symbolSize: 4,
        color: 'red',
        animation: false
      },
      {
        name: 'Sₗ',
        data: [],
        type: 'line',
        symbolSize: 4,
        lineStyle: {
          type: 'solid',
          width: 1.5
        },
        color: 'red',
        animation: false
      },

      // {
      //   name: 'Sᵤ*',
      //   data: [],
      //   type: 'line',
      //   symbolSize: 4,
      //   lineStyle: {
      //     type: 'dashed',
      //     width: 1.5
      //   },
      //   color: 'orange',
      //   animation: false
      // },
      // {
      //   name: 'Sₗ*',
      //   data: [],
      //   type: 'line',
      //   symbolSize: 4,
      //   lineStyle: {
      //     type: 'dashed',
      //     width: 1.5
      //   },
      //   color: 'orange',
      //   animation: false
      // }
    ]
  }

  let chart = {}
  let slots = actual.value.split(',').map(Number)
  function draw() {
    option.xAxis.data = []
    for (let i = 0; i < 4; i++) {
      option.series[i].data = []
    }
    for (let i = 0; i <= L; i++) {
      option.xAxis.data.push(i)
      option.series[0].data.push(Math.round(10e5 * i * intf.alpha) / 10e5)
    }
    option.yAxis.max = L * intf.alpha
    let supplied = 0
    let minIns = 0
    let maxIns = 0
    for (let i = 0; i <= L; i++) {
      // for (let i = 0; i <= slots.length; i++) {
      if (i <= slots[slots.length - 1]) {
        if (slots[supplied] > 0 && i == slots[supplied]) {
          supplied++
        }
        // supplied = slots[i]
        let ins = Number((supplied - option.series[0].data[i]).toFixed(5))
        if (minIns > ins) {
          minIns = ins
        }
        if (maxIns < ins) {
          maxIns = ins
        }
      }
      else {
        break
      }
      option.series[1].data.push(supplied)
    }
    for (let i = 0; i <= L; i++) {
      option.series[2].data.push(Math.round(10e5 * (i * intf.alpha + minIns + intf.r)) / 10e5)
      option.series[3].data.push(Math.round(10e5 * (i * intf.alpha + maxIns - intf.r)) / 10e5)

      // if (maxIns != 0)
      //   option.series[4].data.push(Math.round(10e5 * (i * intf.alpha + maxIns)) / 10e5)
      // if (minIns != 0)
      //   option.series[5].data.push(Math.round(10e5 * (i * intf.alpha + minIns)) / 10e5)
    }
    chart.setOption(option)
  }
  watch(actual, () => {
    if (actual.value[actual.value.length - 1] == ',')
      slots = actual.value.slice(0, -1).split(',').map(Number)
    else slots = actual.value.split(',').map(Number)
    draw()
  })
  watch(intf, () => {
    if (actual.value[actual.value.length - 1] == ',')
      slots = actual.value.slice(0, -1).split(',').map(Number)
    else slots = actual.value.split(',').map(Number)
    draw()
  })
  onMounted(() => {
    chart = echarts.init(chartDom.value)
    draw()
  })
}
