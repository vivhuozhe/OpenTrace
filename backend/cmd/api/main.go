package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vivhuozhe/OpenTrace/backend/internal/models"
	"github.com/vivhuozhe/OpenTrace/backend/internal/api"
	"github.com/vivhuozhe/OpenTrace/backend/internal/graph"
	"github.com/vivhuozhe/OpenTrace/backend/internal/spatial"
)

func main() {
	dsn := "host=localhost port=5432 user=admin password=lemmaballs dbname=opentrace_map sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	repo := &spatial.MapRepo{DB: db}
	nodes, edges, _ := repo.GetGraphData()

	router := &graph.Router{
		Nodes: make(map[int]models.Node),
		Edges: make(map[int][]models.Edge),
	}
	for _, n := range nodes { router.Nodes[n.ID] = n }
	for _, e := range edges { router.Edges[e.SourceNodeID] = append(router.Edges[e.SourceNodeID], e) }

	h := &api.Handler{Repo: repo, Router: router}
	r := gin.Default()

	r.GET("/map/:id", h.GetMap)
	r.GET("/route", h.GetRoute)

	log.Println("🚀 OpenTrace Backend running on :8080")
	r.Run(":8080")
}
