import { onMounted } from 'vue'

import * as echarts from 'echarts/core'
import { BoxplotChart, ScatterChart } from 'echarts/charts'
import {
  ToolboxComponent,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  DatasetComponent,
  TransformComponent,
  DataZoomComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([
  ToolboxComponent,
  DatasetComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  TransformComponent,
  BoxplotChart,
  DataZoomComponent,
  ScatterChart,
  CanvasRenderer,
])

import res from './latency.json'

export function useDrawLatency(chartDom) {
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
    grid: {
      top: "25px",
      right: "10px",
      left: "45px",
      bottom: "45px"
    },
    legend: {
      // width: "80%",
      itemWidth:18,
      itemHeight:10
      // right: "10px"
    },
    // dataZoom: {},
    dataset: [
      {
        source: [
          res["15"],
          res["29"],
          res["37"],
          res["22"],
          res["25"],
          res["45"],

          res["3"],
          res["4"],

          res["9"],
          res["20"],
          res["32"],
          res["47"],
        ],
      },
      {
        transform: {
          type: "boxplot",
          config: {},
        },
      },
      {
        fromDatasetIndex: 1,
        fromTransformResult: 1,
      },
    ],
    xAxis: {
      name: 'Latency',
      type: 'value',
      // data: [],
      axisLabel: {
        fontSize: 12
      },
      // min: 0,
      nameLocation: 'center',
      nameGap: 25,
      nameTextStyle: {
        fontSize: 13
      }
    },
    yAxis: {
      name: 'Sensor',
      type: 'category',
      // data: [],
      nameGap: 30,
      axisLabel: {
        fontSize: 12,
        formatter: (item) => {
          return ["15", "29", "37", "22", "25", "45", "3", "4", "9", "20", "32", "47"][item]
        }
      },
      inverse: true,
      nameTextStyle: {
        fontSize: 13
      },
      nameLocation: 'center',
    },
    series: [
      {
        name: "Latency",
        type: "boxplot",
        boxWidth: 20,
        datasetIndex: 1,
        itemStyle: {
          borderWidth: 2,
        }
      },
      {
        name: "Outlier",
        type: "scatter",
        datasetIndex: 2,
        color: "red",
        symbolSize: 5,
        symbol: "emptyCircle"
      },
    ]
  }

  let chart = {}
  function draw() {

    console.log(option.dataset[0])
    chart.setOption(option)
  }

  onMounted(() => {
    chart = echarts.init(chartDom.value)
    draw()
  })
}
