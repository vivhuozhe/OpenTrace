package models

type Node struct {
    ID         int     `json:"id"`
    LevelID    int     `json:"level_id"`
    Name       string  `json:"name"`
    Category   string  `json:"category"`
    Geom       string  `json:"-"`
    Lat        float64 `json:"lat"`
    Lon        float64 `json:"lon"`
    FloorLabel string  `json:"floor_label"`
}

type Edge struct {
	ID             int     `db:"id" json:"id"`
	LevelID        int     `db:"level_id" json:"level_id"`
	SourceNodeID   int     `db:"source_node_id" json:"source_node_id"`
	TargetNodeID   int     `db:"target_node_id" json:"target_node_id"`
	DistanceMeters float64 `db:"distance_meters" json:"distance_meters"`
}
