package graph

import (
	"math"
	"sort"

	"github.com/vivhuozhe/OpenTrace/backend/internal/models"
)

type Router struct {
	Nodes map[int]models.Node
	Edges map[int][]models.Edge
}

type PathResult struct {
	Nodes []models.Node
	Cost  float64
}

func (r *Router) FindPath(startID, endID int) ([]models.Node, error) {
	paths := r.FindKShortestPaths(startID, endID, 1)
	if len(paths) == 0 {
		return nil, nil
	}
	return paths[0].Nodes, nil
}

func (r *Router) FindKShortestPaths(startID, endID int, k int) []PathResult {
	if k < 1 {
		k = 1
	}
	if k > 5 {
		k = 5
	}

	A := []PathResult{}
	B := []PathResult{}

	firstPath := r.dijkstra(startID, endID)
	if firstPath == nil {
		return A
	}
	A = append(A, PathResult{Nodes: firstPath, Cost: r.calculatePathCost(firstPath)})

	for k := 1; k < len(A[0].Nodes)-1 && len(A) < k; k++ {
		for i := 0; i < len(A[0].Nodes)-1; i++ {
			spurnNode := A[0].Nodes[i]
			rootPath := A[0].Nodes[:i+1]

			removedEdges := []models.Edge{}
			for _, path := range A {
				if len(path.Nodes) > i {
					for j := 0; j < len(path.Nodes)-1; j++ {
						if r.nodesEqual(path.Nodes[j], spurnNode) && r.nodesEqual(path.Nodes[j+1], rootPath[i+1]) {
							edge := r.getEdge(path.Nodes[j].ID, path.Nodes[j+1].ID)
							if edge != nil {
								removedEdges = append(removedEdges, *edge)
								r.removeEdge(edge.SourceNodeID, edge.TargetNodeID)
							}
						}
					}
				}
			}

			for _, edge := range r.Edges[spurnNode.ID] {
				if edge.TargetNodeID == rootPath[i+1].ID {
					removedEdges = append(removedEdges, edge)
					r.removeEdge(edge.SourceNodeID, edge.TargetNodeID)
				}
			}

			candidatePath := r.dijkstra(startID, endID)
			if candidatePath != nil {
				cost := r.calculatePathCost(candidatePath)
				if !r.pathExistsIn(A, candidatePath) && !r.pathExistsIn(B, candidatePath) {
					B = append(B, PathResult{Nodes: candidatePath, Cost: cost})
				}
			}

			for _, edge := range removedEdges {
				r.restoreEdge(edge)
			}
		}

		if len(B) == 0 {
			break
		}

		sort.Slice(B, func(i, j int) bool {
			return B[i].Cost < B[j].Cost
		})

		A = append(A, B[0])
		B = B[1:]
	}

	return A
}

func (r *Router) dijkstra(startID, endID int) []models.Node {
	dist := make(map[int]float64)
	prev := make(map[int]int)
	for id := range r.Nodes {
		dist[id] = math.MaxFloat64
	}
	dist[startID] = 0

	pq := []int{startID}

	for len(pq) > 0 {
		curr := pq[0]
		pq = pq[1:]

		if curr == endID {
			return r.reconstructPath(prev, endID)
		}

		for _, edge := range r.Edges[curr] {
			neighbor := edge.TargetNodeID
			newDist := dist[curr] + edge.DistanceMeters

			if newDist < dist[neighbor] {
				dist[neighbor] = newDist
				prev[neighbor] = curr
				pq = append(pq, neighbor)
			}
		}
	}
	return nil
}

func (r *Router) nodesEqual(a, b models.Node) bool {
	return a.ID == b.ID
}

func (r *Router) getEdge(from, to int) *models.Edge {
	for _, edge := range r.Edges[from] {
		if edge.TargetNodeID == to {
			return &edge
		}
	}
	return nil
}

func (r *Router) removeEdge(from, to int) {
	newEdges := []models.Edge{}
	for _, edge := range r.Edges[from] {
		if edge.TargetNodeID != to {
			newEdges = append(newEdges, edge)
		}
	}
	r.Edges[from] = newEdges
}

func (r *Router) restoreEdge(edge models.Edge) {
	r.Edges[edge.SourceNodeID] = append(r.Edges[edge.SourceNodeID], edge)
}

func (r *Router) calculatePathCost(nodes []models.Node) float64 {
	cost := 0.0
	for i := 0; i < len(nodes)-1; i++ {
		cost += r.GetDistance(nodes[i].ID, nodes[i+1].ID)
	}
	return cost
}

func (r *Router) pathExistsIn(paths []PathResult, nodes []models.Node) bool {
	for _, p := range paths {
		if len(p.Nodes) == len(nodes) {
			match := true
			for i := range p.Nodes {
				if p.Nodes[i].ID != nodes[i].ID {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}
	return false
}

func (r *Router) reconstructPath(prev map[int]int, endID int) []models.Node {
	path := []models.Node{}
	for curr := endID; curr != 0; curr = prev[curr] {
		path = append([]models.Node{r.Nodes[curr]}, path...)
	}
	return path
}

func (r *Router) GetDistance(fromID, toID int) float64 {
	for _, edge := range r.Edges[fromID] {
		if edge.TargetNodeID == toID {
			return edge.DistanceMeters
		}
	}
	return 0
}
