package api

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/vivhuozhe/OpenTrace/backend/internal/graph"
	"github.com/vivhuozhe/OpenTrace/backend/internal/spatial"
)

type PathNode struct {
    ID         int     `json:"id"`
    Name       string  `json:"name"`
    Floor      string  `json:"floor"`
    Category   string  `json:"category"`
    Lat        float64 `json:"lat"`
    Lon        float64 `json:"lon"`
}

type RouteResponse struct {
    Summary struct {
        TotalDistance float64 `json:"total_distance_meters"`
        StartPoint    string  `json:"start_point"`
        EndPoint      string  `json:"end_point"`
        NodeCount     int     `json:"node_count"`
    } `json:"summary"`
    Path []PathNode `json:"path"`
}

type Handler struct {
	Repo   *spatial.MapRepo
	Router *graph.Router
}

func (h *Handler) GetMap(c *gin.Context) {
	levelID, _ := strconv.Atoi(c.Param("id"))
	data, err := h.Repo.GetLevelGeoJSON(levelID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Data(200, "application/json", []byte(data))
}

func (h *Handler) GetRoute(c *gin.Context) {
    startID, _ := strconv.Atoi(c.Query("start"))
    endID, _ := strconv.Atoi(c.Query("end"))

    pathNodes, err := h.Router.FindPath(startID, endID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Path not found"})
        return
    }

    var response RouteResponse
    var totalDist float64

    for i, node := range pathNodes {
        response.Path = append(response.Path, PathNode{
            ID:       node.ID,
            Name:     node.Name,
            Floor:    node.FloorLabel,
            Category: node.Category,
            Lat:      node.Lat,
            Lon:      node.Lon,
        })

        if i < len(pathNodes)-1 {
            totalDist += h.Router.GetDistance(node.ID, pathNodes[i+1].ID)
        }
    }

    response.Summary.TotalDistance = totalDist
    response.Summary.NodeCount = len(pathNodes)
    response.Summary.StartPoint = pathNodes[0].Name
    response.Summary.EndPoint = pathNodes[len(pathNodes)-1].Name

    c.JSON(http.StatusOK, response)
}
