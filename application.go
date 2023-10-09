package main

import (
	"sort"
)

type Application struct {
	Network     `json:"-"`
	ID          int     `json:"id"`
	NumChannels int     `json:"num_channels"`
	HyperPeriod int     `json:"hyperperiod"`
	NumTasks    int     `json:"num_tasks"`
	Tasks       []*Task `json:"tasks"`

	UsedNodes     []int                `json:"used_nodes"`
	CriticalNodes map[int]bool         `json:"critical_nodes"`
	PreDG         []*PrecedenceDAG     `json:"pre_dep_graph"`
	TxDG          TransmissionDepGraph `json:"tx_dep_graph"`
	Schedule      [2][][]Cell          `json:"schedule"` // ideal and actual

	NI map[int]NodeAccessInterface `json:"interfaces"`
	BI BandwidthInterface          `json:"bandwidth_interface"`
}

type Task struct {
	ID                int     `json:"id"`
	AppID             int     `json:"app_id"`
	Source            int     `json:"src"`
	Destination       int     `json:"dst"`
	Paths             [][]int `json:"-"`
	Path              []int   `json:"path"`
	Release           int     `json:"release"`
	Period            int     `json:"period"`
	Deadline          int     `json:"deadline"`
	AbsDeadline       int     `json:"-"`
	CurInstance       int     `json:"-"`
	ScheduledHops     int     `json:"-"`
	FinishTime        []int   `json:"-"`
	ScheduledInstance int     `json:"scheduled_instance"`
}

type Cell struct {
	Assigned       bool `json:"assigned"`
	AppID          int  `json:"app_id"`
	TaskID         int  `json:"-"`
	TaskHopID      int  `json:"-"`
	TaskInstanceID int  `json:"-"`
	Sender         int  `json:"sender"`
	Receiver       int  `json:"receiver"`
}

func (nw *Network) NewApplication(id, hyperPeriod int) *Application {
	schIdeal := make([][]Cell, 1+nw.Settings.NUM_SLOTS)
	schActual := make([][]Cell, 1+nw.Settings.NUM_SLOTS)
	for i := range schIdeal {
		schIdeal[i] = make([]Cell, 1+nw.Settings.NUM_CHANNELS)
		schActual[i] = make([]Cell, 1+nw.Settings.NUM_CHANNELS)
	}
	app := &Application{
		Network:       *nw,
		ID:            id,
		NumChannels:   4,
		NumTasks:      nw.RandApp.Intn(nw.Settings.NUM_TASK_MAX_APP) + 1,
		HyperPeriod:   hyperPeriod,
		CriticalNodes: make(map[int]bool),
		TxDG:          TransmissionDepGraph{make(map[int][]int)},
		NI:            make(map[int]NodeAccessInterface),
		Schedule:      [2][][]Cell{schIdeal, schActual},
	}
	// case study
	// app := &Application{
	// 	Network:     *nw,
	// 	ID:          id,
	// 	NumChannels: 5,
	// 	NumTasks:    0,
	// 	// HyperPeriod:   hyperPeriod,
	// 	CriticalNodes: make(map[int]bool),
	// 	TxDG:          TransmissionDepGraph{make(map[int][]int)},
	// 	NI:            make(map[int]NodeAccessInterface),
	// 	Schedule:      [2][][]Cell{schIdeal, schActual},
	// }
	// if id == 0 {
	// 	app.NumChannels = 6
	// } else if id == 1 {
	// 	app.NumChannels = 1
	// } else if id == 2 {
	// 	app.NumChannels = 4
	// }

	return app
}

func (app *Application) GenTasks() {
	var taskPeriodCandidate []int

	for pp := 8; pp <= app.HyperPeriod; pp++ {
		if app.HyperPeriod%pp == 0 {
			taskPeriodCandidate = append(taskPeriodCandidate, pp)
		}
	}

	for i := 0; i < app.NumTasks; i++ {
		// p := app.HyperPeriod
		// if i > 0 {
		p := taskPeriodCandidate[app.RandApp.Intn(len(taskPeriodCandidate))]
		// }

	L:
		task := Task{
			ID:          i,
			AppID:       app.ID,
			Source:      app.RandApp.Intn(app.Settings.NUM_NODES),
			Destination: app.RandApp.Intn(app.Settings.NUM_NODES),
			Period:      p,
			Deadline:    p,
		}
		valid := false
		if task.Source != task.Destination {
			app.Nodes[task.Source].FindAllPaths(task.Destination)
			task.Paths = app.Nodes[task.Source].RoutingTable[task.Destination]

			if len(task.Paths) > 0 {
				if len(task.Paths[0]) < p {
					valid = true
					// t.Path = t.Paths[0]
					app.Tasks = append(app.Tasks, &task)
				}
			}
		}
		if !valid {
			goto L
		}
	}
}

func (app *Application) GenTasksCaseStudy() {
	switch app.ID {
	case 0:
		for _, s := range []int{15, 29, 37} {
			app.Nodes[s].FindAllPaths(1)
			app.Tasks = append(app.Tasks, &Task{
				ID:          app.NumTasks,
				AppID:       app.ID,
				Source:      s,
				Destination: 1,
				Period:      50,
				Deadline:    50,
				Paths:       app.Nodes[s].RoutingTable[1],
			})
			app.NumTasks++
		}
		for _, s := range []int{22, 25} {
			app.Nodes[s].FindAllPaths(1)
			app.Tasks = append(app.Tasks, &Task{
				ID:          app.NumTasks,
				AppID:       app.ID,
				Source:      s,
				Destination: 1,
				Period:      100,
				Deadline:    100,
				Paths:       app.Nodes[s].RoutingTable[1],
			})
			app.NumTasks++
		}
		for _, s := range []int{45} {
			app.Nodes[s].FindAllPaths(1)
			app.Tasks = append(app.Tasks, &Task{
				ID:          app.NumTasks,
				AppID:       app.ID,
				Source:      s,
				Destination: 1,
				Period:      100,
				Deadline:    100,
				Paths:       app.Nodes[s].RoutingTable[1],
			})
			app.NumTasks++
		}

	case 1:
		for _, s := range []int{3, 4} {
			app.Nodes[s].FindAllPaths(1)
			app.Tasks = append(app.Tasks, &Task{
				ID:          app.NumTasks,
				AppID:       app.ID,
				Source:      s,
				Destination: 1,
				Period:      10,
				Deadline:    10,
				Paths:       app.Nodes[s].RoutingTable[1],
			})
			app.NumTasks++
		}

	case 2:
		for _, s := range []int{9, 20, 32, 47} {
			app.Nodes[s].FindAllPaths(1)
			app.Tasks = append(app.Tasks, &Task{
				ID:          app.NumTasks,
				AppID:       app.ID,
				Source:      s,
				Destination: 1,
				Period:      100,
				Deadline:    100,
				Paths:       app.Nodes[s].RoutingTable[1],
			})
			app.NumTasks++
		}
	}

}

func (app *Application) ConstructDepGraphs() {
	for _, t := range app.Tasks {
		nodes := make(map[int]bool)
		var vertices [][]int
		for i := 0; i < len(t.Path)-1; i++ {
			nodes[t.Path[i]] = true
			nodes[t.Path[i+1]] = true
			vertices = append(vertices, []int{t.Path[i], t.Path[i+1]})
		}
		app.PreDG = append(app.PreDG, &PrecedenceDAG{
			ID:            t.ID,
			InvolvedNodes: nodes,
			Vertices:      vertices,
			Period:        t.Period,
		})
	}

	// find channel dependency
	for n := range app.CriticalNodes {
		app.TxDG.Graph[n] = []int{}
	}
	peersOfNode := make(map[int][]int)
	for n := range app.CriticalNodes {
		for _, t := range app.Tasks {
			pos := IntSliceIndexOf(t.Path, n)
			if pos > -1 {
				if pos-1 >= 0 {
					peersOfNode[n] = append(peersOfNode[n], t.Path[pos-1])
				}
				if pos+1 < len(t.Path) {
					peersOfNode[n] = append(peersOfNode[n], t.Path[pos+1])
				}
			}
		}
	}
	for n1 := range app.CriticalNodes {
		for n2 := range app.CriticalNodes {
			if n1 != n2 && len(IntSliceIntersect(peersOfNode[n1], peersOfNode[n2])) > 0 {
				if IntSliceIndexOf(app.TxDG.Graph[n1], n2) == -1 {
					app.TxDG.Graph[n1] = append(app.TxDG.Graph[n1], n2)
				}
				if IntSliceIndexOf(app.TxDG.Graph[n2], n1) == -1 {
					app.TxDG.Graph[n2] = append(app.TxDG.Graph[n2], n1)
				}
			}
		}
	}
}

// Abstract interfaces for resource in the resource pool
func (app *Application) AbstractInterfaces() {
	app.AbstractBandwidthInterface()
	for n := range app.CriticalNodes {
		app.AbstractNodeInterface(n)
	}
}

func (app *Application) AbstractBandwidthInterface() {
	slots := []int{}
	for slot := 1; slot <= app.Settings.NUM_SLOTS; slot++ {
		slots = append(slots, slot)
	}
	channels := []int{}
	for ch := 1; ch <= app.NumChannels; ch++ {
		channels = append(channels, ch)
	}
	app.ComputeSchedule(0, slots, channels, nil)
	N := 0
	for slot := 1; slot <= app.Settings.NUM_SLOTS; slot++ {
		for _, cell := range app.Schedule[0][slot] {
			if cell.Assigned {
				N++
				break
			}
		}
	}
	if N%2 == 1 {
		N += 1
	}
	// N += 8
	alpha := Round(float64(N*2) / float64(app.Settings.NUM_SLOTS))
	app.BI = BandwidthInterface{
		NumChannels:        app.NumChannels,
		AppID:              app.ID,
		AvailabilityFactor: alpha,
		Regularity:         1,
	}
}

func (app *Application) AbstractNodeInterface(node int) {
	txTasks := [][2]int{}

	for _, t := range app.Tasks {
		if idx := IntSliceIndexOf(t.Path, node); idx != -1 {
			c := 2
			if idx == 0 || idx == len(t.Path)-1 {
				c = 1
			}
			txTasks = append(txTasks, [2]int{c, t.Period})
		}
	}

	var N, K, u float64 = 1, 0, 0
	var P = float64(app.Settings.NUM_SLOTS)
	minPeriod := app.Settings.NUM_SLOTS
	for _, t := range txTasks {
		u += Round(float64(t[0]) / float64(t[1]))
		if minPeriod > t[1] {
			minPeriod = t[1]
		}
	}
	N = Round(u * P)
	alpha := N / P

	if minPeriod == app.Settings.NUM_SLOTS {
		R := Round(float64(app.Settings.NUM_SLOTS)*(alpha-alpha*alpha) + 0.001)
		if R < 1 {
			R = 1
		}
		app.NI[node] = NodeAccessInterface{
			Node:               node,
			AppID:              app.ID,
			AvailabilityFactor: alpha,
			Regularity:         R,
		}
		return
	}

	u = 0

L:
	for k := float64(0); k < float64(app.Settings.NUM_SLOTS); k++ {
		for _, t := range txTasks {
			if float64(t[1])-k == 0 {
				break L
			}
			if t[1] < app.Settings.NUM_SLOTS {
				u += Round(float64(t[0]) / (float64(t[1]) - k))
			}
		}
		if u <= N/P {
			K = k
		} else {
			break
		}
	}
	R := 1 + K*N/P
	app.NI[node] = NodeAccessInterface{
		Node:               node,
		AppID:              app.ID,
		AvailabilityFactor: alpha,
		Regularity:         R + 0.5,
	}
}

func (app *Application) ComputeSchedule(flag uint8, slots []int, channels []int, availableNodes []map[int]bool) {
	// reset tasks
	for _, t := range app.Tasks {
		t.CurInstance = 0
		t.ScheduledHops = 0
		t.AbsDeadline = 0
		t.FinishTime = make([]int, 1+app.Settings.NUM_SLOTS/t.Period)
		t.ScheduledInstance = 0
	}

	for _, slot := range slots {
		var candidateTasks []*Task
		for _, t := range app.Tasks {
			curInstance := (slot-1)/t.Period + 1
			if t.CurInstance < curInstance {
				t.CurInstance = curInstance
				t.ScheduledHops = 0
				t.AbsDeadline = (curInstance-1)*t.Period + t.Deadline
				// fmt.Println("release", slot, curInstance)
			}
			if t.ScheduledHops < len(t.Path)-1 {
				candidateTasks = append(candidateTasks, t)
			}
		}

		// EDF
		sort.SliceStable(candidateTasks, func(i, j int) bool {
			if candidateTasks[i].AbsDeadline == candidateTasks[j].AbsDeadline {
				if len(candidateTasks[i].Path) == len(candidateTasks[j].Path) {
					return candidateTasks[i].ID < candidateTasks[j].ID
				}
				return len(candidateTasks[i].Path) > len(candidateTasks[j].Path)
			}
			return candidateTasks[i].AbsDeadline < candidateTasks[j].AbsDeadline
		})

		ch := 0
		for _, task := range candidateTasks {
			sender := task.Path[task.ScheduledHops]
			receiver := task.Path[task.ScheduledHops+1]
			conflict := false
			for _, c := range channels {
				cell := app.Schedule[flag][slot][c]
				if cell.Assigned &&
					(cell.Sender == sender || cell.Sender == receiver || cell.Receiver == sender || cell.Receiver == receiver) {
					conflict = true
					break
				}
			}
			permit := true
			if len(availableNodes) == 0 && conflict {
				permit = false
			} else if len(availableNodes) > 0 {
				if app.CriticalNodes[sender] && !availableNodes[slot][sender] {
					permit = false
				}
				if app.CriticalNodes[receiver] && !availableNodes[slot][receiver] {
					permit = false
				}
			}
			if permit {
				task.ScheduledHops++

				app.Schedule[flag][slot][channels[ch]] = Cell{
					Assigned:       true,
					AppID:          app.ID,
					TaskID:         task.ID,
					TaskHopID:      task.ScheduledHops,
					TaskInstanceID: task.CurInstance,
					Sender:         sender,
					Receiver:       receiver,
				}
				ch++
				if task.ScheduledHops == len(task.Path)-1 {
					task.FinishTime[task.CurInstance] = slot - (task.CurInstance-1)*task.Period
					task.ScheduledInstance++
				}
				if ch == len(channels) {
					break
				}
			}
		}
	}
}

func (app *Application) Report() bool {
	slots := []int{}
	for slot := 1; slot <= app.Settings.NUM_SLOTS; slot++ {
		for _, cell := range app.Manager.PartitionBandwidth[slot] {
			if cell.Assigned && cell.AppID == app.ID {
				slots = append(slots, slot)
				break
			}
		}
	}
	channels := []int{}
	for ch := 1; ch <= app.NumChannels; ch++ {
		channels = append(channels, ch)
	}
	availableNodes := make([]map[int]bool, app.Settings.NUM_SLOTS+1)
	for slot := 1; slot <= app.Settings.NUM_SLOTS; slot++ {
		availableNodes[slot] = make(map[int]bool)
		for node, cell := range app.Manager.PartitionNode[slot] {
			if cell.Assigned && cell.AppID == app.ID {
				availableNodes[slot][node] = true
			}
		}

	}
	// fmt.Println(slots, channels, availableNodes)
	app.ComputeSchedule(1, slots, channels, availableNodes)
	scheduledCnt := 0
	for _, t := range app.Tasks {
		if t.ScheduledInstance == app.Settings.NUM_SLOTS/t.Period {
			scheduledCnt++
		}
	}
	if app.Settings.Verbose {
		app.logger.Printf("[App #%d] schedulable:%v\n", app.ID, scheduledCnt == app.NumTasks)
	}
	return scheduledCnt == app.NumTasks
}
