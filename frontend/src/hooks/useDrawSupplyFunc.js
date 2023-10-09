import { network } from '../hooks/useNetwork.js'

import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  MarkLineComponent,
  TitleComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

import { chartTheme } from './useConfig.js'

echarts.use([
  LineChart,
  TitleComponent,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  MarkLineComponent,
  CanvasRenderer
])

export function useDrawSupplyFunc(app, resource) {
  let option = {
    // legend: {},
    title: {
      // text: `A(${app})R(${resource})`
    },
    grid: {
      top: '30px',
      left: '30px',
      right: '30px',
      bottom: '30px'
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: [],
      axisLabel: {
        // fontSize: 13
        interval:19
      },
      // nameTextStyle: {
      //   fontSize: 13
      // }
    },
    yAxis: {
      name: 'Supply',
      type: 'value',
      min: 0,
      max:1,
      minInterval: "20",
      // axisLabel: {
      //   fontSize: 9
      // },
      // nameTextStyle: {
      //   fontSize: 13
      // }
    },
    series: [
      {
        name: 'Uniform',
        data: [],
        type: 'line',
        lineStyle: {
          type: 'dashed',
          width:1.2,
        },
        symbolSize: 0,
        color: 'grey',
        animation: false
      },
      {
        name: 'Up',
        data: [],
        type: 'line',
        symbolSize: 0,
        lineStyle: {
          type: 'dashed',
          width:1.2,
        },
        color: 'red',
        animation: false
      },
      {
        name: 'Lo',
        data: [],
        type: 'line',
        symbolSize: 0,
        lineStyle: {
          type: 'dashed',
          width:1.2,
        },
        color: 'red',
        animation: false
      },
      {
        name: 'Supply',
        data: [],
        type: 'line',
        color: 'blue',
        symbolSize: 0,
        animation: false,
        lineStyle:{
          
          width:1.2,
        }
      },
      {
        name: 'Soft Up',
        data: [],
        type: 'line',
        lineStyle: {
          type: 'dashed',
          width:1.2,
        },
        color: 'orange',
        symbolSize: 0,
        animation: false
      },
      {
        name: 'Soft Lo',
        data: [],
        type: 'line',
        lineStyle: {
          type: 'dashed',
          width:1.2,
        },
        color: 'orange',
        symbolSize: 0,
        animation: false
      }
    ]
  }

  let chart = {}
  function draw() {
    option.xAxis.data = []
    for (let i = 0; i < option.series.length; i++) {
      option.series[i].data = []
    }
    // if (Object.keys(network.value.manager.supply_graphs[0]).length == 0) return
    if (network.value.manager.supply_graphs[app][resource] == null) return
    // const graph =
    //   network.value.manager.supply_graphs[0][Object.keys(network.value.manager.supply_graphs[0])[0]]
    const graph = network.value.manager.supply_graphs[app][resource]
    for (let i = 0; i <= network.value.settings.num_slots; i++) {
      option.xAxis.data.push(i)
    }
    option.yAxis.max = graph.uniform[network.value.settings.num_slots]

    option.series[0].data = graph.uniform
    option.series[1].data = graph.upper_bound
    option.series[2].data = graph.lower_bound
    option.series[3].data = graph.supply_func
    option.series[4].data = graph.soft_upper_bound
    option.series[5].data = graph.soft_lower_bound
    chart.setOption(option)
  }

  // use getID since ref's array doesn't guarantee the same order
  // https://vuejs.org/guide/essentials/template-refs.html#refs-inside-v-for
  chart = echarts.init(document.getElementById(`supply-${app}-${resource}`), chartTheme)
  draw()
}
