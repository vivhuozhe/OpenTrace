package graph

import (
	"math"
	"github.com/vivhuozhe/OpenTrace/backend/internal/models"
)

type Router struct {
	Nodes map[int]models.Node
	Edges map[int][]models.Edge
}

func (r *Router) FindPath(startID, endID int) ([]models.Node, error) {
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
			return r.reconstructPath(prev, endID), nil
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
	return nil, nil
}

func (r *Router) reconstructPath(prev map[int]int, endID int) []models.Node {
	path := []models.Node{}
	for curr := endID; curr != 0; curr = prev[curr] {
		path = append([]models.Node{r.Nodes[curr]}, path...)
	}
	return path
}
