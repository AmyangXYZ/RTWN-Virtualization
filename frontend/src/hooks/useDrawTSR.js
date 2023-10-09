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

import result from './tsr-a.json'

export function useDrawTSR(chartDom) {
  let option = {
    toolbox: {
      feature: {
        saveAsImage: {
          pixelRatio: 8
        }
      }
    },
    legend: {},
    grid:{
      top:"50px",
      right:"0",
      left:"35px",
      bottom:"45px"
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      name: 'Number of Applications - N',
      type: 'category',
      data: result.x,
      axisLabel: {
        fontSize: 14,
        // interval:0,
      },
      nameLocation: 'center',
      nameGap: 30,
      nameTextStyle: {
        fontSize: 15
      }
    },
    yAxis: {
      name: 'TSR (%)',
      type: 'value',
      min: 0,
      axisLabel: {
        fontSize: 14
      },
      nameTextStyle: {
        fontSize: 15
      }
    },
    series: [
      {
        name: "SMT",
        data: result.smt,
        type: 'line',
        // symbol: symbols[s],
        symbolSize: 7.5,
        lineStyle: { width: 3 },
        animation: false,
        // color: colors[s]
      },
      {
        name: "SGP",
        data: result.sgp,
        type: 'line',
        symbol:"emptyRect",
        // symbol: symbols[s],
        symbolSize: 7.5,
        lineStyle: { width: 3 },
        animation: false,
        // color: colors[s]
      },   
      {
        name: "RRP",
        data: result.rrp,
        type: 'line',
        symbol:"emptyTriangle",
        // symbol: symbols[s],
        symbolSize: 7.5,
        lineStyle: { width: 3 },
        animation: false,
        // color: colors[s]
      },  
      {
        name: "RR",
        data: result.rr,
        symbol:"emptyDiamond",
        type: 'line',
        // symbol: symbols[s],
        symbolSize: 7.5,
        lineStyle: { width: 3 },
        animation: false,
        // color: colors[s]
      },  
      {
        name: "EDP",
        data: result.edp,
        type: 'line',
        symbol:"emptyRoundRect",
        // symbol: symbols[s],
        symbolSize: 7.5,
        lineStyle: { width: 3 },
        animation: false,
        // color: colors[s]
      },     
    ]
  }

  let chart = {}
  function draw() {
    chart.setOption(option)
  }

  onMounted(() => {
    chart = echarts.init(chartDom.value)
    draw()
  })
}
