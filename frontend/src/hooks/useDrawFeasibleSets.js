import { onMounted, watch, ref } from 'vue'

import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  ToolboxComponent,
  TooltipComponent,
  VisualMapComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([
  LineChart,
  TooltipComponent,
  GridComponent,
  ToolboxComponent,
  VisualMapComponent,
  CanvasRenderer
])

import { combinadic, combination } from './combinatorics.js'

export function useDrawFeasibleSets(settings, chartDom) {
  let option = {
    tooltip: {
      trigger: 'item',
      formatter: (item) => {
        return `slots: [${item.data.name}]`
      }
    },
    toolbox: {
      feature: {
        saveAsImage: {
          pixelRatio: 8
        }
      }
    },
    grid: {
      top: '7%',
      left:'0%',
      bottom:'0%',
      right:"0%"
    },
    xAxis: {
      type: 'category',
      data: [],
      axisLabel: {
        show: false
      }
    },
    yAxis: {
      type: 'category',
      data: [],
      axisLabel: {
        show: false
      }
    },
    visualMap: {
      type: 'piecewise',
      pieces: [
        { value: 0, color: 'lightgrey', label: 'SC' },
        { value: 1, color: 'deepskyblue', label: 'FS' },
        { value: 2, color: 'red', label: 'FP' },
        { value: 3, color: 'purple', label: 'FS+FP' }
      ],
      left: 'center',
      top: 'top',
      width:"100%",
      orient:"horizontal",
    },
    series: [
      {
        name: 'Slots',
        data: [],
        type: 'heatmap',
        itemStyle: {
          borderWidth: 1,
          borderColor: 'white'
        },
        color: 'grey',
        animation: false
      }
    ]
  }

  let chart = {}
  let sc_size = ref(0)
  let fs_size = ref(0)
  let fp_size = ref(0)
  let unsafe = ref(0)
  function draw() {
    sc_size.value = 0
    fs_size.value = 0
    fp_size.value = 0
    unsafe.value = 0

    const T = settings.T
    const N = settings.N
    const p = settings.p
    const c = settings.c
    const r = settings.r
    const alpha = settings.alpha

    if (T > 0 && p > 0 && c > 0 && r > 0) {
      sc_size.value = combination(T, N)
      fs_size.value = 0
      fp_size.value = 0

      const X = parseInt(Number(Math.sqrt(parseInt(sc_size.value))))
      let Y = X
      while (Y * X < sc_size.value) Y++
      option.xAxis.data = []
      option.yAxis.data = []
      option.series[0].data = []
      for (let i = 0; i < X; i++) {
        option.xAxis.data.push(i)
      }
      for (let i = 0; i < Y; i++) {
        option.yAxis.data.push(i)
      }

      const comb = combinadic(T, N)
      for (let i = 0; i < sc_size.value; i++) {
        const s = comb(i)
        for (let j = 0; j < s.length; j++) {
          s[j]++
        }

        let minIns = 0
        let maxIns = 0
        let supplied = 0
        let supply_func = [0]
        for (let t = 1; t <= T; t++) {
          if (t > s[s.length - 1]) {
            break
          }
          if (s[supplied] > 0 && t === s[supplied]) {
            supplied++
          }
          supply_func.push(supplied)
          const ins = Number((supplied - alpha * t).toFixed(5))
          if (minIns > ins) {
            minIns = ins
          }
          if (maxIns < ins) {
            maxIns = ins
          }
        }

        const l = supply_func.length
        if (l < T + 1) {
          for (let tt = 0; tt < T + 1 - l; tt++) {
            supply_func.push(supplied)
          }
        }
        let label = 0
        // check feasible schedule
        let feasible_instance = 0
        for (let k = 0; k < Math.floor(T / p); k++) {
          if (Number((supply_func[(k + 1) * p] - supply_func[k * p]).toFixed(5)) >= c) {
            feasible_instance += 1
          }
        }
        if (feasible_instance == T / p) {
          fs_size.value += 1
          label = 1
        }

        // check feasible partition
        if (Math.abs(Number((maxIns - minIns).toFixed(5))) < r) {
          fp_size.value++
          if (label == 1) {
            label = 3
          } else {
            label = 2
            unsafe.value++
          }
        }

        option.series[0].data.push({
          name: s,
          value: [i % X, parseInt(i / X), label]
        })
      }

      chart.setOption(option)
    }
  }
  watch(settings, () => {
    draw()
  })

  onMounted(() => {
    chart = echarts.init(chartDom.value)
    draw()
  })
  return { sc_size, fs_size, fp_size, unsafe }
}
