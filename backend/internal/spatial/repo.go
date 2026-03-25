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
	var nodes []models.Node
	var edges []models.Edge

	err := r.DB.Select(&nodes, "SELECT id, level_id, name, category, ST_Y(geom) as lat, ST_X(geom) as lon FROM nodes")
	if err != nil {
		return nil, nil, err
	}

	err = r.DB.Select(&edges, "SELECT id, source_node_id, target_node_id, distance_meters FROM edges")
	return nodes, edges, err
}
