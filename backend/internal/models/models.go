package models

type Node struct {
	ID       int     `db:"id" json:"id"`
	LevelID  int     `db:"level_id" json:"level_id"`
	Name     string  `db:"name" json:"name"`
	Category string  `db:"category" json:"category"`
	Lat      float64 `db:"lat" json:"lat"`
	Lon      float64 `db:"lon" json:"lon"`
}

type Edge struct {
	ID             int     `db:"id" json:"id"`
	LevelID        int     `db:"level_id" json:"level_id"`
	SourceNodeID   int     `db:"source_node_id" json:"source_node_id"`
	TargetNodeID   int     `db:"target_node_id" json:"target_node_id"`
	DistanceMeters float64 `db:"distance_meters" json:"distance_meters"`
}
