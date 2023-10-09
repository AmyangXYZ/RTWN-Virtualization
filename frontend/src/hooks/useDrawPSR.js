import { onMounted } from 'vue'

import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
  ToolboxComponent,
  GridComponent,
  TooltipComponent,
  LegendComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([
  LineChart,
  ToolboxComponent,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  CanvasRenderer
])

import result from './psr.json'

export function useDrawPSR(chartDom) {
  let colors = {
    SMT: '#5470c6',
    SGP_A: '#91cc75',
    SGP_B: 'red',
    AAF: 'grey',
    Magic7: '#9a60b4',
    MulZ: '#fc8452'
  }
  let symbols = {
    SMT: 'emptycircle',
    SGP_A: 'emptyrect',
    SGP_B: 'emptyrect',
    AAF: 'emptytriangle',
    Magic7: 'emptytriangle',
    MulZ:'emptytriangle'
  }
  let option = {
    toolbox: {
      feature: {
        saveAsImage: {
          pixelRatio: 8
        }
      }
    },
    grid:{
      top:"55px",
      right:"0",
      left:"35px",
      bottom:"45px"
    },
    legend: {
      width:"70%",
      // align:"center"
      // right:"10px"
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      name: 'Utilization',
      type: 'category',
      data: [],
      axisLabel: {
        fontSize: 14
      },
      nameLocation: 'center',
      nameGap: 30,
      nameTextStyle: {
        fontSize: 15
      }
    },
    yAxis: {
      name: 'PSR (%)',
      type: 'value',
      min: 0,
      axisLabel: {
        fontSize: 14
      },
      nameTextStyle: {
        fontSize: 15
      }
    },
    series: []
  }

  let chart = {}
  function draw() {
    option.xAxis.data = result.xAxis
    option.series = []
    for (let s in result.data) {
      option.series.push({
        name: s,
        data: result.data[s],
        type: 'line',
        symbol: symbols[s],
        symbolSize: 7.5,
        lineStyle: { width: 3 },
        animation: false,
        color: colors[s]
      })
    }

    chart.setOption(option)
  }

  onMounted(() => {
    chart = echarts.init(chartDom.value)
    draw()
  })
}
