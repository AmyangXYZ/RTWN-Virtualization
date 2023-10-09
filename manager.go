package main

import (
	"sort"
)

type NetworkManager struct {
	Network       `json:"-"`
	CriticalNodes map[int][]int         `json:"critical_nodes"` // node: []apps
	SupplyGraphs  []map[int]SupplyGraph `json:"supply_graphs"`

	PartitionNode      [][]Cell `json:"partition_node"`
	PartitionBandwidth [][]Cell `json:"partition_bandwidth"`
}

func (nw *Network) NewManager() *NetworkManager {
	if nw.Settings.Verbose {
		nw.logger.Println("Creating manager...")
	}
	pb := make([][]Cell, nw.Settings.NUM_SLOTS+1)
	pn := make([][]Cell, nw.Settings.NUM_SLOTS+1)
	for c := range pb {
		pb[c] = make([]Cell, 1+nw.Settings.NUM_CHANNELS)
		pn[c] = make([]Cell, nw.Settings.NUM_NODES)
	}
	sg := make([]map[int]SupplyGraph, nw.Settings.NUM_APPS)
	for i := range sg {
		sg[i] = make(map[int]SupplyGraph)
	}
	return &NetworkManager{
		Network:            *nw,
		CriticalNodes:      make(map[int][]int),
		SupplyGraphs:       sg,
		PartitionBandwidth: pb,
		PartitionNode:      pn,
	}
}

// should find disjoint paths
func (nm *NetworkManager) Routing() {
	// um := make(map[int]bool)
	for i := 0; i < nm.Settings.NUM_APPS; i++ {
		app := nm.Apps[i]
		usedNodesMap := make(map[int]bool)
		for _, t := range app.Tasks {
			// shortest path
			// t.Path = t.Paths[0]
			// for _, n := range t.Path {
			// 	usedNodesMap[n] = true
			// }
			// sort.Ints(app.UsedNodes)
			// fmt.Println(t.Path)
			// continue
			// disjoint path among apps
		L:
			for _, p := range t.Paths {
				disjoint := true
				for _, n := range p {
					for _, appOther := range nm.Apps {
						if appOther.ID != app.ID {
							if IntSliceIndexOf(appOther.UsedNodes, n) > -1 {
								disjoint = false
								continue L
							}
						}
					}
				}
				if disjoint {
					t.Path = p
					break
				}
			}
			// can't find disjoint path
			if len(t.Path) == 0 {
				t.Path = t.Paths[0]
			}
			for _, n := range t.Path {
				usedNodesMap[n] = true
				// um[n] = true
			}
		}

		for n := range usedNodesMap {
			app.UsedNodes = append(app.UsedNodes, n)
		}
		sort.Ints(app.UsedNodes)
	}
	// tmp := []int{}
	// for n := range um {
	// 	tmp = append(tmp, n)
	// }
	// sort.Ints(tmp)
	// s, _ := json.Marshal(tmp)
	// fmt.Println(string(s))
}

func (nm *NetworkManager) Partitioning() {
	nm.FindCriticalNodes()
	for _, app := range nm.Apps {
		app.AbstractInterfaces()
		app.ConstructDepGraphs()
	}
	switch nm.Settings.METHOD {
	case METHOD_SGP:
		nm.AllocateBandwidthPartition()
		nm.AllocateNodePartitions()
	case METHOD_SMT:
		nm.callSMT()
	case METHOD_RRP:
		nm.AllocateBandwidthPartitionRRP()
		nm.AllocateNodePartitions()
	case METHOD_RR:
		nm.AllocateBandwidthPartitionRR()
	}
	nm.DistributePartitions()
}

func (nm *NetworkManager) FindCriticalNodes() {
	// for _, app := range nm.Apps {
	// 	for _, v := range app.UsedNodes {
	// 		app.CriticalNodes[v] = true
	// 		nm.CriticalNodes[v] = []int{app.ID}
	// 	}
	// }
	// return

	for _, app1 := range nm.Apps {
		for _, app2 := range nm.Apps {
			if app1.ID != app2.ID {
				common := IntSliceIntersect(app1.UsedNodes, app2.UsedNodes)
				for _, n := range common {
					app1.CriticalNodes[n] = true
					if _, exist := nm.CriticalNodes[n]; exist {
						if IntSliceIndexOf(nm.CriticalNodes[n], app1.ID) == -1 {
							nm.CriticalNodes[n] = append(nm.CriticalNodes[n], app1.ID)
						}
						if IntSliceIndexOf(nm.CriticalNodes[n], app2.ID) == -1 {
							nm.CriticalNodes[n] = append(nm.CriticalNodes[n], app2.ID)
						}
					} else {
						nm.CriticalNodes[n] = []int{app1.ID, app2.ID}
					}
				}
			}
		}
	}
}

func (nm *NetworkManager) AllocateBandwidthPartitionRRP() {
	for _, app := range nm.Apps {
		app.BI.Regularity = 1
		for n, intf := range app.NI {
			intf.Regularity = 1
			app.NI[n] = intf
		}
	}
	nm.AllocateBandwidthPartition()
}

func (nm *NetworkManager) AllocateBandwidthPartitionRR() {
	for slot := 1; slot <= nm.Settings.NUM_SLOTS; slot++ {
		for ch := range nm.PartitionBandwidth[slot] {
			nm.PartitionBandwidth[slot][ch] = Cell{
				Assigned: true,
				AppID:    (slot + ch) % len(nm.Apps),
			}
		}
	}
}

func (nm *NetworkManager) AllocateBandwidthPartition() {
	nm.InitBandwidthSupplyGraph()

	for slot := 1; slot <= nm.Settings.NUM_SLOTS; slot++ {
		var requests []AccessRequest
		for i := range nm.Apps {
			requests = append(requests, nm.MakeBandwidthRequests(slot, i)...)
		}
		sort.SliceStable(requests, func(i, j int) bool {
			return requests[i].Weight > requests[j].Weight
		})

		// fmt.Printf("[Slot %d] Allocat node %d to app %d\n", slot, n, selectedApp)
		c := 1

		for _, req := range requests {
			if req.Weight > -1 {
				selectedApp := nm.Apps[req.AppID]
				if c+selectedApp.NumChannels <= 1+nm.Settings.NUM_CHANNELS {
					nm.SupplyGraphs[selectedApp.ID][-1].SupplyFunc[slot]++
					for cc := c; cc < c+selectedApp.NumChannels; cc++ {
						nm.PartitionBandwidth[slot][cc] = Cell{
							Assigned: true,
							AppID:    selectedApp.ID,
						}
					}
					c += selectedApp.NumChannels
				}
			}
		}
		nm.UpdateBandwidthSupplyGraph(slot)
	}
}

func (nm *NetworkManager) InitBandwidthSupplyGraph() {
	for i, app := range nm.Apps {
		intf := app.BI
		n := -1 // use -1 to represent bandwidth resource
		g := SupplyGraph{
			ID:                 n,
			AvailabilityFactor: intf.AvailabilityFactor,
			Regularity:         intf.Regularity,
			SupplyFunc:         make([]float64, nm.Settings.NUM_SLOTS+1),
			Uniform:            make([]float64, nm.Settings.NUM_SLOTS+1),
			UpperBound:         make([]float64, nm.Settings.NUM_SLOTS+1),
			LowerBound:         make([]float64, nm.Settings.NUM_SLOTS+1),
			SoftUpperBound:     make([]float64, nm.Settings.NUM_SLOTS+1),
			SoftLowerBound:     make([]float64, nm.Settings.NUM_SLOTS+1),
		}
		for slot := 0; slot <= nm.Settings.NUM_SLOTS; slot++ {
			g.Uniform[slot] = Round(g.AvailabilityFactor * float64(slot))
			g.UpperBound[slot] = Round(g.Uniform[slot] + g.Regularity)
			g.LowerBound[slot] = Round(g.Uniform[slot] - g.Regularity)
			copy(g.SoftLowerBound, g.LowerBound)
			copy(g.SoftUpperBound, g.UpperBound)
		}
		nm.SupplyGraphs[i][n] = g
	}
}

// upon response from manager, update supply func and the bounds
func (nm *NetworkManager) UpdateBandwidthSupplyGraph(slot int) {
	for i := range nm.Apps {
		n := -1
		g := nm.SupplyGraphs[i][n]

		// next slot
		if slot < nm.Settings.NUM_SLOTS {
			g.SupplyFunc[slot+1] = g.SupplyFunc[slot]
		}

		// update upper and lower bound
		insReg := Round(g.SupplyFunc[slot] - g.Uniform[slot])
		if g.MinInsReg > insReg {
			g.MinInsReg = insReg
			for t := 0; t <= nm.Settings.NUM_SLOTS; t++ {
				g.UpperBound[t] = Round(g.MinInsReg + g.Uniform[t] + float64(g.Regularity))
				g.SoftLowerBound[t] = Round(g.MinInsReg + g.Uniform[t])
			}
		}
		if g.MaxInsReg < insReg {
			g.MaxInsReg = insReg
			for t := 0; t <= nm.Settings.NUM_SLOTS; t++ {
				g.LowerBound[t] = Round(g.MaxInsReg + g.Uniform[t] - float64(g.Regularity))
				g.SoftUpperBound[t] = Round(g.MaxInsReg + g.Uniform[t])
			}
		}
		if g.SoftUpperBound[slot] > g.UpperBound[slot] || g.MaxInsReg == 0 {
			for t := 0; t <= nm.Settings.NUM_SLOTS; t++ {
				g.SoftUpperBound[t] = g.UpperBound[t]
			}
		}
		if g.SoftLowerBound[slot] < g.LowerBound[slot] || g.MinInsReg == 0 {
			for t := 0; t <= nm.Settings.NUM_SLOTS; t++ {
				g.SoftLowerBound[t] = g.LowerBound[t]
			}
		}
		nm.SupplyGraphs[i][n] = g
	}
}

func (nm *NetworkManager) MakeBandwidthRequests(slot, appID int) []AccessRequest {
	requests := []AccessRequest{}
	n := -1

	weight := nm.AccessRequestWeight(slot, appID, n)
	requests = append(requests, AccessRequest{
		Node:   n,
		AppID:  appID,
		Weight: Round(weight),
	})
	// fmt.Println(slot, requests)
	return requests
}

func (nm *NetworkManager) AllocateNodePartitions() {
	nm.InitNodeSupplyGraph()

	for slot := 1; slot <= nm.Settings.NUM_SLOTS; slot++ {
		requests := make(map[int][]AccessRequest)
		for i := range nm.Apps {
			for _, req := range nm.MakeNodeRequests(slot, i) {
				requests[req.Node] = append(requests[req.Node], req)
			}
		}

		candidates := []AccessRequest{}
		for _, reqs := range requests {
			if len(reqs) == 0 {
				continue
			}
			sort.SliceStable(reqs, func(i, j int) bool {
				return reqs[i].Weight > reqs[j].Weight
			})
			if reqs[0].Weight != -1 {
				candidates = append(candidates, reqs[0])
			}

		}

		for _, req := range candidates {
			// don't allocate this node without peer
			if req.Peer > -1 {
				found := false
				for _, req2 := range candidates {
					if req2.Node == req.Peer && req2.AppID == req.AppID && req2.DAGID == req.DAGID {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			// fmt.Printf("[Slot %d] Allocat node %d to app %d\n", slot, req.Node, req.AppID)
			if req.HopInDAG == nm.Apps[req.AppID].PreDG[req.DAGID].curHop {
				nm.Apps[req.AppID].PreDG[req.DAGID].curHop++
			}
			nm.SupplyGraphs[req.AppID][req.Node].SupplyFunc[slot]++
			nm.PartitionNode[slot][req.Node] = Cell{
				Assigned: true,
				AppID:    req.AppID,
			}
		}

		nm.UpdateNodeSupplyGraph(slot)
	}
}

func (nm *NetworkManager) InitNodeSupplyGraph() {
	for i, app := range nm.Apps {
		for n := range app.CriticalNodes {
			intf := app.NI[n]
			g := SupplyGraph{
				ID:                 n,
				AvailabilityFactor: intf.AvailabilityFactor,
				Regularity:         intf.Regularity,
				SupplyFunc:         make([]float64, nm.Settings.NUM_SLOTS+1),
				Uniform:            make([]float64, nm.Settings.NUM_SLOTS+1),
				UpperBound:         make([]float64, nm.Settings.NUM_SLOTS+1),
				LowerBound:         make([]float64, nm.Settings.NUM_SLOTS+1),
				SoftUpperBound:     make([]float64, nm.Settings.NUM_SLOTS+1),
				SoftLowerBound:     make([]float64, nm.Settings.NUM_SLOTS+1),
			}
			for slot := 0; slot <= nm.Settings.NUM_SLOTS; slot++ {
				g.Uniform[slot] = Round(g.AvailabilityFactor * float64(slot))
				g.UpperBound[slot] = Round(g.Uniform[slot] + g.Regularity)
				g.LowerBound[slot] = Round(g.Uniform[slot] - g.Regularity)
				copy(g.SoftLowerBound, g.LowerBound)
				copy(g.SoftUpperBound, g.UpperBound)
			}
			nm.SupplyGraphs[i][n] = g
		}
	}
}

// upon response from manager, update supply func and the bounds
func (nm *NetworkManager) UpdateNodeSupplyGraph(slot int) {
	for i := range nm.Apps {
		for n, g := range nm.SupplyGraphs[i] {
			if n == -1 {
				continue
			}
			// next slot
			if slot < nm.Settings.NUM_SLOTS {
				g.SupplyFunc[slot+1] = g.SupplyFunc[slot]
			}

			// update upper and lower bound
			insReg := Round(g.SupplyFunc[slot] - g.Uniform[slot])
			if g.MinInsReg > insReg {
				g.MinInsReg = insReg
				for t := 0; t <= nm.Settings.NUM_SLOTS; t++ {
					g.UpperBound[t] = Round(g.MinInsReg + g.Uniform[t] + float64(g.Regularity))
					g.SoftLowerBound[t] = Round(g.MinInsReg + g.Uniform[t])
				}
			}
			if g.MaxInsReg < insReg {
				g.MaxInsReg = insReg
				for t := 0; t <= nm.Settings.NUM_SLOTS; t++ {
					g.LowerBound[t] = Round(g.MaxInsReg + g.Uniform[t] - float64(g.Regularity))
					g.SoftUpperBound[t] = Round(g.MaxInsReg + g.Uniform[t])
				}
			}
			if g.SoftUpperBound[slot] > g.UpperBound[slot] || g.MaxInsReg == 0 {
				for t := 0; t <= nm.Settings.NUM_SLOTS; t++ {
					g.SoftUpperBound[t] = g.UpperBound[t]
				}
			}
			if g.SoftLowerBound[slot] < g.LowerBound[slot] || g.MinInsReg == 0 {
				for t := 0; t <= nm.Settings.NUM_SLOTS; t++ {
					g.SoftLowerBound[t] = g.LowerBound[t]
				}
			}
			nm.SupplyGraphs[i][n] = g
		}
	}
}

func (nm *NetworkManager) MakeNodeRequests(slot, appID int) []AccessRequest {
	requests := []AccessRequest{}
	app := nm.Apps[appID]

	for _, dag := range app.PreDG {
		if slot%dag.Period == 1 {
			dag.curHop = 0
		}
		if dag.curHop == len(dag.Vertices) {
			continue
		}

		v := dag.Vertices[dag.curHop]
		tmp := []AccessRequest{}

		for _, n := range v {
			if !app.CriticalNodes[n] {
				continue
			}
			weight := nm.AccessRequestWeight(slot, appID, n)
			tmp = append(tmp, AccessRequest{
				Node:     n,
				AppID:    appID,
				Weight:   Round(weight),
				DAGID:    dag.ID,
				HopInDAG: dag.curHop,
				Peer:     -1,
			})
		}
		if len(tmp) == 2 {
			tmp[0].Weight = Float64Max(tmp[0].Weight, tmp[1].Weight)
			tmp[1].Weight = Float64Max(tmp[0].Weight, tmp[1].Weight)
			tmp[0].Peer = tmp[1].Node
			tmp[1].Peer = tmp[0].Node
		} else if len(tmp) == 0 {
			// check if non-critical node got enough slot
			cnt := 0
			for ss := slot/dag.Period*dag.Period + 1; ss <= slot; ss++ {
				for _, cell := range nm.PartitionBandwidth[ss] {
					if cell.Assigned && cell.AppID == appID {
						cnt++
						break
					}
				}
			}
			if cnt > dag.curHop {
				dag.curHop++
			}
		}
		requests = append(requests, tmp...)
	}
	// third round search - if channel constraint sat
	for i1, req1 := range requests {
		for i2, req2 := range requests {
			if req1 != req2 && IntSliceIndexOf(app.TxDG.Graph[req1.Node], req2.Node) > -1 {
				if req1.Weight > req2.Weight {
					requests[i2].Weight = -1
					for j, req3 := range requests {
						if req3.DAGID == req2.DAGID && req3.HopInDAG == req2.HopInDAG {
							requests[j].Weight = -1
						}
					}
				} else {
					requests[i1].Weight = -1
					for j, req3 := range requests {
						if req3.DAGID == req1.DAGID && req3.HopInDAG == req1.HopInDAG {
							requests[j].Weight = -1
						}
					}
				}
			}
		}
	}
	return requests
}

// Look up the supply graph
func (nm *NetworkManager) AccessRequestWeight(slot, i, n int) float64 {
	g := nm.SupplyGraphs[i][n]
	var weight float64 = -1

	if n != -1 {
		hasChannelAccess := false
		for _, cell := range nm.PartitionBandwidth[slot] {
			if cell.Assigned && cell.AppID == i {
				hasChannelAccess = true
				break
			}
		}
		if !hasChannelAccess {
			return weight
		}
	}

	if g.SupplyFunc[slot] < g.AvailabilityFactor*float64(nm.Settings.NUM_SLOTS) &&
		g.UpperBound[slot] > g.SupplyFunc[slot]+1 {
		tStartSoft, tEnd, tEndSoft := slot, nm.Settings.NUM_SLOTS, nm.Settings.NUM_SLOTS

		for t := slot; t <= nm.Settings.NUM_SLOTS; t++ {
			if g.LowerBound[t] >= g.SupplyFunc[slot] {
				tEnd = t
				break
			}
		}
		for t := slot; t <= nm.Settings.NUM_SLOTS; t++ {
			if g.SoftUpperBound[t] >= g.SupplyFunc[slot]+1 {
				tStartSoft = t
				break
			}
		}
		for t := slot; t <= nm.Settings.NUM_SLOTS; t++ {
			if g.SoftLowerBound[t] > g.SupplyFunc[slot] {
				tEndSoft = t
				break
			}
		}

		weight = 10 / float64(tEnd)
		if slot >= tStartSoft && slot <= tEndSoft {
			weight += 10
		}
	}
	return Round(weight)
}

func (nm *NetworkManager) DistributePartitions() {
	for _, app := range nm.Apps {
		app.Manager = nm
	}
}

func (nm *NetworkManager) Report() {
	for appID := range nm.Apps {
		for n, g := range nm.SupplyGraphs[appID] {
			if g.SupplyFunc[nm.Settings.NUM_SLOTS] < g.Uniform[nm.Settings.NUM_SLOTS] {
				if nm.Settings.Verbose {
					nm.logger.Printf("! App %v's partition %v failed\n", appID, n)
				}
			}
		}
	}
}
