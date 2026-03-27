INSERT INTO levels (id, building_name, level_number, label) values
(1, 'A Block', 0, 'Ground Floor'),
(2, 'A Block', 1, 'First Floor'),
(3, 'A Block', 2, 'Second Floor'),
(4, 'A Block', 3, 'Third Floor'),
(5, 'Admin Block', 0, 'Ground Floor'),
(6, 'Admin Block', 1, 'First Floor'),
(7, 'B Block', 0, 'Ground Floor'),
(8, 'B Block', 1, 'First Floor'),
(9, 'B Block', 2, 'Second Floor'),
(10, 'B Block', 3, 'Third Floor');

INSERT INTO nodes (id, level_id, name, category, geom) values
(1, 3, 'Incubation Hub', 'Room', ST_GeomFromText('POINT(77.4765 28.4749)', 4326)),
(2, 5, 'Reception', 'Room', ST_GeomFromText('POINT(77.4763 28.4747)', 4326)),
(3, 4, 'Library', 'Room', ST_GeomFromText('POINT(77.4765 28.4749)', 4326)),
(4, 4, 'Library', 'Room', ST_GeomFromText('POINT(77.4765 28.4749)', 4326)),
(5, 7, 'Seminar Hall', 'Room', ST_GeomFromText('POINT(77.4760 28.4747)', 4326)),
(6, 7, 'Corridor Seminar Hall', 'Corridor', ST_GeomFromText('POINT(77.4760 28.4746)', 4326)),
(7, 1, 'A Block Ground Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4765 28.4748)', 4326)),
(8, 2, 'A Block First Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4765 28.4748)', 4326)),
(9, 3, 'A Block Second Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4765 28.4748)', 4326)),
(10, 4, 'A Block Third Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4765 28.4748)', 4326)),
(11, 1, 'A Block Ground Floor Elevator', 'Elevator', ST_GeomFromText('POINT(77.4765 28.4749)', 4326)),
(12, 2, 'A Block First Floor Elevator', 'Elevator', ST_GeomFromText('POINT(77.4765 28.4749)', 4326)),
(13, 3, 'A Block Second Floor Elevator', 'Elevator', ST_GeomFromText('POINT(77.4765 28.4749)', 4326)),
(14, 4, 'A Block Third Floor Elevator', 'Elevator', ST_GeomFromText('POINT(77.4765 28.4749)', 4326)),
(15, 7, 'B Block Ground Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4762 28.4747)', 4326)),
(16, 8, 'B Block First Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4762 28.4747)', 4326)),
(17, 9, 'B Block Second Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4762 28.4747)', 4326)),
(18, 10, 'B Block Third Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4762 28.4747)', 4326));

-- Connect Node seminal hall corridor to B block ground floor stairs
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    7,
    6,
    15,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 6)::geography, 
        (SELECT geom FROM nodes WHERE id = 15)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 6), 
        (SELECT geom FROM nodes WHERE id = 15)
    );

-- connecting node b block ground floor stairs to reception
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    5,
    15,
    2,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 15)::geography, 
        (SELECT geom FROM nodes WHERE id = 2)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 15), 
        (SELECT geom FROM nodes WHERE id = 2)
    );

-- connecting reception ground floor to A Block ground floor elevator
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    5,
    2,
    11,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 2)::geography, 
        (SELECT geom FROM nodes WHERE id = 11)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 2), 
        (SELECT geom FROM nodes WHERE id = 11)
    );

-- connecting A block ground floor elevator to A Block first floor elevator
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    1,
    11,
    12,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 11)::geography, 
        (SELECT geom FROM nodes WHERE id = 12)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 11), 
        (SELECT geom FROM nodes WHERE id = 12)
    );

-- connecting A block first floor elevator to A Block second floor elevator
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    2,
    12,
    13,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 12)::geography, 
        (SELECT geom FROM nodes WHERE id = 13)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 12), 
        (SELECT geom FROM nodes WHERE id = 13)
    );

-- connecting A block second floor elevator to A Block third floor elevator
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    3,
    13,
    14,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 13)::geography, 
        (SELECT geom FROM nodes WHERE id = 14)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 13), 
        (SELECT geom FROM nodes WHERE id = 14)
    );

-- connecting A block ground second elevator to Incubation Hub
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    3,
    13,
    1,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 13)::geography, 
        (SELECT geom FROM nodes WHERE id = 1)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 13), 
        (SELECT geom FROM nodes WHERE id = 1)
    );