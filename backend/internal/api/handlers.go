package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/vivhuozhe/OpenTrace/backend/internal/graph"
	"github.com/vivhuozhe/OpenTrace/backend/internal/models"
	"github.com/vivhuozhe/OpenTrace/backend/internal/spatial"
	"net/http"
	"strconv"
)

type Level struct {
	ID           int    `json:"id"`
	BuildingName string `json:"building_name"`
	LevelNumber  int    `json:"level_number"`
	Label        string `json:"label"`
}

type PathNode struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Floor    string  `json:"floor"`
	Category string  `json:"category"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
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

type AlternateRoutesResponse struct {
	Routes []RouteResponse `json:"routes"`
}

type Handler struct {
	Repo   *spatial.MapRepo
	Router *graph.Router
	DB     *sqlx.DB
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

func (h *Handler) GetAlternateRoutes(c *gin.Context) {
	startID, _ := strconv.Atoi(c.Query("start"))
	endID, _ := strconv.Atoi(c.Query("end"))
	count, _ := strconv.Atoi(c.Query("count"))
	if count < 1 {
		count = 1
	}
	if count > 5 {
		count = 5
	}

	paths := h.Router.FindKShortestPaths(startID, endID, count)

	var response AlternateRoutesResponse
	for _, pathResult := range paths {
		var route RouteResponse
		var totalDist float64

		for i, node := range pathResult.Nodes {
			route.Path = append(route.Path, PathNode{
				ID:       node.ID,
				Name:     node.Name,
				Floor:    node.FloorLabel,
				Category: node.Category,
				Lat:      node.Lat,
				Lon:      node.Lon,
			})

			if i < len(pathResult.Nodes)-1 {
				totalDist += h.Router.GetDistance(node.ID, pathResult.Nodes[i+1].ID)
			}
		}

		route.Summary.TotalDistance = totalDist
		route.Summary.NodeCount = len(pathResult.Nodes)
		route.Summary.StartPoint = pathResult.Nodes[0].Name
		route.Summary.EndPoint = pathResult.Nodes[len(pathResult.Nodes)-1].Name

		response.Routes = append(response.Routes, route)
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetLevels(c *gin.Context) {
	query := `SELECT id, building_name, level_number, label FROM levels ORDER BY building_name, level_number`
	rows, err := h.Repo.DB.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var levels []Level
	for rows.Next() {
		var l Level
		err := rows.Scan(&l.ID, &l.BuildingName, &l.LevelNumber, &l.Label)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		levels = append(levels, l)
	}
	c.JSON(200, levels)
}

func (h *Handler) GetAllNodes(c *gin.Context) {
	nodes, _, err := h.Repo.GetGraphData()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, nodes)
}

func (h *Handler) GetAllEdges(c *gin.Context) {
	_, edges, err := h.Repo.GetGraphData()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, edges)
}

func (h *Handler) GetMapData(c *gin.Context) {
	levelID, _ := strconv.Atoi(c.Query("level"))
	geojson, err := h.Repo.GetLevelGeoJSON(levelID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Data(200, "application/json", []byte(geojson))
}

type DashboardResponse struct {
	Levels []Level       `json:"levels"`
	Nodes  []models.Node `json:"nodes"`
	Edges  []models.Edge `json:"edges"`
}

func (h *Handler) GetDashboardData(c *gin.Context) {
	levelsQuery := `SELECT id, building_name, level_number, label FROM levels ORDER BY building_name, level_number`
	levelRows, err := h.Repo.DB.Query(levelsQuery)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer levelRows.Close()

	var levels []Level
	for levelRows.Next() {
		var l Level
		err := levelRows.Scan(&l.ID, &l.BuildingName, &l.LevelNumber, &l.Label)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		levels = append(levels, l)
	}

	nodes, edges, err := h.Repo.GetGraphData()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, DashboardResponse{
		Levels: levels,
		Nodes:  nodes,
		Edges:  edges,
	})
}
