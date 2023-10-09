package main

import (
	"sort"
)

type Node struct {
	Network      `json:"-"`
	ID           int             `json:"id"`
	Position     [2]int          `json:"pos"`
	Neighbors    []int           `json:"neighbors"`
	RoutingTable map[int][][]int `json:"-"`
}

func (nw *Network) NewNode(id int) *Node {
	n := &Node{
		Network:      *nw,
		ID:           id,
		RoutingTable: make(map[int][][]int),
	}
	var x, y int

L1:
	for {
		x, y = nw.RandTopo.Intn(nw.Settings.GRID_X-1)+1, nw.RandTopo.Intn(nw.Settings.GRID_Y-1)+1
		overlapping := false
		for _, nn := range nw.Nodes {
			if nn == nil {
				continue
			}
			if x == nn.Position[0] && y == nn.Position[1] {
				overlapping = true
				goto L1
			}
		}
		if !overlapping {
			break
		}
	}

	n.Position = [2]int{x, y}
	return n
}

func (n *Node) FindAllPaths(dst int) []int {
	visited := make(map[int]bool)
	path := []int{}
	n.findPath(n.ID, dst, visited, path)
	sort.SliceStable(n.RoutingTable[dst], func(i, j int) bool {
		if len(n.RoutingTable[dst][i]) == len(n.RoutingTable[dst][j]) {
			return IntSum(n.RoutingTable[dst][i]) < IntSum(n.RoutingTable[dst][j])
		}
		return len(n.RoutingTable[dst][i]) < len(n.RoutingTable[dst][j])
	})
	return path
}

func (n *Node) findPath(cur, dst int, visited map[int]bool, path []int) {
	visited[cur] = true
	path = append(path, cur)

	if cur == dst {
		n.RoutingTable[dst] = append(n.RoutingTable[dst], path)
	} else if len(path) < n.Settings.MAX_HOP {
		for _, nn := range n.Nodes[cur].Neighbors {
			if !visited[nn] {
				pathCpy := make([]int, len(path))
				copy(pathCpy, path)
				n.findPath(nn, dst, visited, pathCpy)
			}
		}
	}
	visited[cur] = false
}
