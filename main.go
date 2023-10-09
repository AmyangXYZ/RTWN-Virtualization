package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	METHOD_SGP = 0
	METHOD_SMT = 1
	METHOD_RRP = 2
	METHOD_RR  = 3
)

var (
	Networks    []*Network
	NumNetworks = 1
)

func main() {
	go RunWeb()
	<-time.After(100 * time.Millisecond)

	X := []float64{}
	sgp := []float64{}
	smt := []float64{}
	limit := make(chan struct{}, 3)
	needSMT := make(map[int][]int64)
	for x := 1; x <= 1; x++ {
		x := x * 50
		X = append(X, float64(x))
		fmt.Printf("x = %v, ", x)
		var sr float64
		var wg sync.WaitGroup
		for i := 0; i < NumNetworks; i++ {
			wg.Add(1)
			i := i
			limit <- struct{}{}
			go func(sr *float64) {
				defer func() {
					<-limit
					wg.Done()
				}()
				nw := NewNetwork(i, SystemSettings{
					NUM_APPS:         4,
					NUM_NODES:        x,
					NUM_SLOTS:        50,
					NUM_CHANNELS:     16,
					NUM_TASK_MAX_APP: 3,
					RAND_SEED_TOPO:   79,
					RAND_SEED_APP:    165 + int64(i),
					GRID_X:           32,
					GRID_Y:           24,
					TX_RANGE:         10,
					MAX_HOP:          4,
					METHOD:           METHOD_SGP,
					Verbose:          false,
				})
				Networks = append(Networks, nw)
				nw.Run()
				// nw.logger.Println("success:", nw.Stat.TaskSuccess)
				if nw.Stat.TaskSuccess {
					*sr++
				} else {
					needSMT[x] = append(needSMT[x], 165+int64(i))
				}
			}(&sr)
		}
		wg.Wait()
		fmt.Printf("sr: %v%%\n", Round(sr/float64(NumNetworks)*100))
		sgp = append(sgp, Round(sr/float64(NumNetworks)*100))
	}

	if len(needSMT) < 0 {
		fmt.Println("call SMT")

		tmp := make(map[int]float64)
		limit := make(chan struct{}, 8)
		for x, v := range needSMT {
			var wg sync.WaitGroup
			sr := sgp[x/2-1] * float64(NumNetworks) / 100

			for _, i := range v {
				// limit # of concurrent smt caller

				wg.Add(1)
				i := i
				// go func() {
				limit <- struct{}{}
				// }()
				go func(sr *float64) {
					defer func() {
						<-limit
						wg.Done()
					}()
					nw := NewNetwork(int(i), SystemSettings{
						NUM_APPS:         x,
						NUM_NODES:        120,
						NUM_SLOTS:        50,
						NUM_CHANNELS:     16,
						NUM_TASK_MAX_APP: 3,
						RAND_SEED_TOPO:   79,
						RAND_SEED_APP:    i,
						GRID_X:           32,
						GRID_Y:           24,
						TX_RANGE:         20,
						MAX_HOP:          3,
						METHOD:           METHOD_SMT,
						Verbose:          false,
					})
					Networks = append(Networks, nw)
					nw.Run()
					// nw.logger.Println("success:", nw.Stat.TaskSuccess)
					if nw.Stat.TaskSuccess {
						*sr++
					}

				}(&sr)
			}

			wg.Wait()
			fmt.Printf("smt: x = %v, sr: %v%%\n", x, sr/float64(NumNetworks)*100)
			tmp[x] = sr / float64(NumNetworks) * 100
		}

		for _, x := range X {
			if sr, ok := tmp[int(x)]; ok {
				smt = append(smt, sr)
			} else {
				smt = append(smt, 100)
			}
		}
	}
	dump := map[string][]float64{
		"x":   X,
		"sgp": sgp,
		"smt": smt,
	}
	j, _ := json.Marshal(dump)
	fmt.Println(string(j))
	select {}
}
