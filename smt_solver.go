package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

const NON_CRITICAL_NODE = -999

type smtInput struct {
	Option struct {
		Verbose     bool  `json:"verbose"`
		NumSlots    int   `json:"num_slots"`
		NumChannels int   `json:"num_channels"`
		NumApps     int   `json:"num_apps"`
		Resources   []int `json:"resources"`
	} `json:"option"`
	Interfaces []map[int][2]float64 `json:"interfaces"`
	DAGs       [][]smtInputDAG      `json:"dags"`
	CGs        []map[int][]int      `json:"cg"`
}

type smtInputDAG struct {
	V [][]int `json:"v"`
	P int     `json:"p"`
}

type smtOutput struct {
	Partitions []map[int][]int `json:"partitions"`
}

func (nm *NetworkManager) callSMT() {
	if nm.Settings.Verbose {
		nm.logger.Println("Calling SMT solver...")
	}
	in := smtInput{}
	in.Option.Verbose = false
	in.Option.NumSlots = nm.Settings.NUM_SLOTS
	in.Option.NumChannels = nm.Settings.NUM_CHANNELS
	in.Option.NumApps = nm.Settings.NUM_APPS
	for n := range nm.CriticalNodes {
		in.Option.Resources = append(in.Option.Resources, n)
	}
	for _, app := range nm.Apps {
		intf := make(map[int][2]float64)
		bandwidth := -1 * (app.ID + 1)
		// if app.ID == 1 {
		// 	bandwidth = -2
		// }
		in.Option.Resources = append(in.Option.Resources, bandwidth)
		intf[bandwidth] = [2]float64{app.BI.AvailabilityFactor, app.BI.Regularity}
		for n, v := range app.NI {
			intf[n] = [2]float64{v.AvailabilityFactor, v.Regularity}
		}
		in.Interfaces = append(in.Interfaces, intf)
		dags := []smtInputDAG{}
		for _, dag := range app.PreDG {
			dags = append(dags, smtInputDAG{
				V: dag.Vertices,
				P: dag.Period,
			})
		}
		in.DAGs = append(in.DAGs, dags)

		in.CGs = append(in.CGs, app.TxDG.Graph)
	}

	inputJSON, _ := json.Marshal(in)
	// fmt.Println(string(inputJSON))
	outputJSON, err := exec.Command("python", "-u", "./pyscripts/smt-solver.py", string(inputJSON)).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(outputJSON))
	out := smtOutput{}
	err = json.Unmarshal(outputJSON, &out)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(out)
	for i, partitions := range out.Partitions {
		for r, p := range partitions {
			// fmt.Println(i, r, p)
			for _, slot := range p {
				if r < 0 {
					ch := 1
					for _, cell := range nm.PartitionBandwidth[slot] {
						if cell.Assigned {
							ch++
						}
					}
					for c := ch; c < ch+nm.Apps[i].NumChannels; c++ {
						nm.PartitionBandwidth[slot][c] = Cell{
							Assigned: true,
							AppID:    i,
						}
					}
				} else {
					nm.PartitionNode[slot][r] = Cell{
						Assigned: true,
						AppID:    i,
					}
				}
			}
		}
	}
}
