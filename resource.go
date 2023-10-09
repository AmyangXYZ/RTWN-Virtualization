package main

type BandwidthInterface struct {
	NumChannels        int     `json:"num_channels"`
	AppID              int     `json:"app_id"`
	AvailabilityFactor float64 `json:"alpha"`
	Regularity         float64 `json:"r"`
}

type BandwidthPartition struct {
	NumChannels int   `json:"num_channels"`
	AppID       int   `json:"app_id"`
	Slots       []int `json:"slots"`
}

type NodeAccessInterface struct {
	Node               int     `json:"node"`
	AppID              int     `json:"app_id"`
	AvailabilityFactor float64 `json:"alpha"`
	Regularity         float64 `json:"r"`
}

type NodeAccessPartition struct {
	Node  int   `json:"node"`
	AppID int   `json:"app_id"`
	Slots []int `json:"slots"`
}

type AccessRequest struct {
	Node     int     `json:"node"`
	AppID    int     `json:"app_id"`
	Weight   float64 `json:"weight"`
	DAGID    int     `json:"dag_id"`
	HopInDAG int     `json:"hop_in_dag"`
	Peer     int     `json:"peer"`
}

type Link struct {
	ID int `json:"id"`
	N1 int `json:"n1"`
	N2 int `json:"n2"`
}

// l1 -> l2
type DepEdge struct {
	Type int `json:"type"`
	L1   int `json:"l1"`
	L2   int `json:"l2"`
}

type PrecedenceDAG struct {
	ID            int          `json:"id"`
	InvolvedNodes map[int]bool `json:"nodes"`
	Vertices      [][]int      `json:"vertices"`
	Period        int          `json:"period"`

	curHop int
}

// a simple undirected graph
type TransmissionDepGraph struct {
	Graph map[int][]int `json:"graph"`
}

type SupplyGraph struct {
	ID                 int       `json:"id"` // resource id
	AvailabilityFactor float64   `json:"alpha"`
	Regularity         float64   `json:"r"`
	SupplyFunc         []float64 `json:"supply_func"`
	Uniform            []float64 `json:"uniform"`
	MinInsReg          float64   `json:"min_ins_reg"`
	MaxInsReg          float64   `json:"max_ins_reg"`
	UpperBound         []float64 `json:"upper_bound"`
	LowerBound         []float64 `json:"lower_bound"`
	SoftUpperBound     []float64 `json:"soft_upper_bound"`
	SoftLowerBound     []float64 `json:"soft_lower_bound"`
}
