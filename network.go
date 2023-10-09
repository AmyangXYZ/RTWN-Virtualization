package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
)

type Network struct {
	ID       int             `json:"id"`
	Settings SystemSettings  `json:"settings"`
	RandTopo *rand.Rand      `json:"-"`
	RandApp  *rand.Rand      `json:"-"`
	Nodes    []*Node         `json:"nodes"`
	Apps     []*Application  `json:"apps"`
	Manager  *NetworkManager `json:"manager"`
	Stat     *Statistics     `json:"stat"`

	logger *log.Logger
}

type SystemSettings struct {
	NUM_APPS         int     `json:"num_apps"`
	NUM_NODES        int     `json:"num_nodes"`
	NUM_SLOTS        int     `json:"num_slots"`
	NUM_CHANNELS     int     `json:"num_channels"`
	NUM_TASK_MAX_APP int     `json:"num_task_max_app"`
	RAND_SEED_TOPO   int64   `json:"seed_topo"`
	RAND_SEED_APP    int64   `json:"seed_app"`
	GRID_X           int     `json:"grid_x"`
	GRID_Y           int     `json:"grid_y"`
	TX_RANGE         float64 `json:"tx_range"`
	MAX_HOP          int     `json:"max_hop"`
	METHOD           int     `json:"method"`
	Verbose          bool    `json:"-"`
}

func NewNetwork(id int, settings SystemSettings) *Network {
	nw := Network{
		ID:     id,
		logger: log.New(os.Stdout, fmt.Sprintf("[Network #%d]: ", id), 0),
	}
	nw.Settings = settings

	return &nw
}

func (nw *Network) Run() {
	s, _ := json.Marshal(nw.Settings)
	if nw.Settings.Verbose {
		nw.logger.Println("Settings:", string(s))
	}
	nw.RandTopo = rand.New(rand.NewSource(nw.Settings.RAND_SEED_TOPO))
	nw.RandApp = rand.New(rand.NewSource(nw.Settings.RAND_SEED_APP))
	nw.Stat = nw.NewStat()

	nw.CreateNodes()
	// nw.CreateNodesCaseStudy()
	nw.CreateApps()
	nw.CreateNetworkManger()
	nw.CollectStat()
}

func (nw *Network) CreateNodes() {
	if nw.Settings.Verbose {
		nw.logger.Println("Creating nodes...")
	}
	// init
	nw.Nodes = make([]*Node, nw.Settings.NUM_NODES)
	for i := 0; i < nw.Settings.NUM_NODES; i++ {
		n := nw.NewNode(i)
		nw.Nodes[i] = n
	}

	// form topology
	for _, n := range nw.Nodes {
		for _, nn := range nw.Nodes {
			if n.ID != nn.ID {
				distance := math.Pow(float64(n.Position[0]-nn.Position[0]), 2) + math.Pow(float64(n.Position[1]-nn.Position[1]), 2)
				if distance <= math.Pow(nw.Settings.TX_RANGE, 2) {
					n.Neighbors = append(n.Neighbors, nn.ID)
				}
			}
		}
	}

	// construct routing table
	// for _, n := range nw.Nodes {
	// 	for _, nn := range nw.Nodes {
	// 		if n.ID != nn.ID {
	// 			n.FindAllPaths(nn.ID)
	// 		}
	// 	}
	// }
}

type topo struct {
	Data []struct {
		SensorID int `json:"sensor_id"`
		Parent   int `json:"parent"`
	} `json:"data"`
}

func (nw *Network) CreateNodesCaseStudy() {
	if nw.Settings.Verbose {
		nw.logger.Println("Creating nodes...")
	}
	// init
	nw.Nodes = make([]*Node, nw.Settings.NUM_NODES)
	for i := 0; i < nw.Settings.NUM_NODES; i++ {
		n := nw.NewNode(i)
		nw.Nodes[i] = n
	}

	// read topo.json
	content, err := os.ReadFile("./topology_49.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var t topo
	err = json.Unmarshal(content, &t)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	for _, n := range t.Data {
		nw.Nodes[n.SensorID].Neighbors = append(nw.Nodes[n.SensorID].Neighbors, n.Parent)
		nw.Nodes[n.Parent].Neighbors = append(nw.Nodes[n.Parent].Neighbors, n.SensorID)
	}

}

func (nw *Network) CreateApps() {
	if nw.Settings.Verbose {
		nw.logger.Println("Creating applications...")
	}
	nw.Apps = make([]*Application, nw.Settings.NUM_APPS)

	var hyperPeriodCandidate []int
	for p := 8; p <= nw.Settings.NUM_SLOTS; p++ {
		if nw.Settings.NUM_SLOTS%p == 0 {
			cnt := 0
			for i := 8; i < p; i++ {
				if p%i == 0 {
					cnt++
				}
			}
			if cnt >= 2 {
				hyperPeriodCandidate = append(hyperPeriodCandidate, p)
			}
		}
	}

	for i := 0; i < nw.Settings.NUM_APPS; i++ {
		hp := hyperPeriodCandidate[nw.RandApp.Intn(len(hyperPeriodCandidate))]
		// hp = t.Settings.NUM_SLOTS
		app := nw.NewApplication(i, hp)
		nw.Apps[i] = app
	}

	// separate this to make the highlevel t.Settings, like policy,
	// hyperperiod, num_task be independent with topology
	for _, app := range nw.Apps {
		app.GenTasks()
		// app.GenTasksCaseStudy()
	}
}

func (nw *Network) CreateNetworkManger() {
	nw.Manager = nw.NewManager()
	nw.Manager.Routing()
	nw.Manager.Partitioning()
}
