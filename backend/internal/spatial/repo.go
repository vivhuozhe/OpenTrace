package spatial

import (
	"github.com/jmoiron/sqlx"
	"github.com/vivhuozhe/OpenTrace/backend/internal/models"
)

type MapRepo struct {
	DB *sqlx.DB
}

func (r *MapRepo) GetLevelGeoJSON(levelID int) (string, error) {
	var geojson string
	query := `
		SELECT jsonb_build_object(
			'type',     'FeatureCollection',
			'features', jsonb_agg(features.feature)
		)
		FROM (
		  SELECT jsonb_build_object(
			'type',       'Feature',
			'geometry',   ST_AsGeoJSON(geom)::jsonb,
			'properties', jsonb_build_object('id', id, 'name', name, 'category', category)
		  ) AS feature
		  FROM nodes WHERE level_id = $1
		) AS features;`

	err := r.DB.Get(&geojson, query, levelID)
	return geojson, err
}

func (r *MapRepo) GetGraphData() ([]models.Node, []models.Edge, error) {
	nodeQuery := `
		SELECT
			n.id,
			n.level_id,
			n.name,
			n.category,
			ST_X(n.geom) as lon,
			ST_Y(n.geom) as lat,
			l.label as floor_label
		FROM nodes n
		JOIN levels l ON n.level_id = l.id
	`

	rows, err := r.DB.Query(nodeQuery)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var n models.Node
		err := rows.Scan(&n.ID, &n.LevelID, &n.Name, &n.Category, &n.Lon, &n.Lat, &n.FloorLabel)
		if err != nil {
			return nil, nil, err
		}
		nodes = append(nodes, n)
	}

	edgeQuery := `
		SELECT id, level_id, source_node_id, target_node_id, distance_meters
		FROM edges
	`
	edgeRows, err := r.DB.Query(edgeQuery)
	if err != nil {
		return nil, nil, err
	}
	defer edgeRows.Close()

	var edges []models.Edge
	for edgeRows.Next() {
		var e models.Edge
		err := edgeRows.Scan(&e.ID, &e.LevelID, &e.SourceNodeID, &e.TargetNodeID, &e.DistanceMeters)
		if err != nil {
			return nil, nil, err
		}
		edges = append(edges, e)
	}

	return nodes, edges, nil
}
