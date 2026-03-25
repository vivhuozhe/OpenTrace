package api

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/vivhuozhe/OpenTrace/backend/internal/graph"
	"github.com/vivhuozhe/OpenTrace/backend/internal/spatial"
)

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
	start, _ := strconv.Atoi(c.Query("start"))
	end, _ := strconv.Atoi(c.Query("end"))

	path, err := h.Router.FindPath(start, end)
	if err != nil || path == nil {
		c.JSON(404, gin.H{"error": "No path found"})
		return
	}
	c.JSON(200, path)
}
