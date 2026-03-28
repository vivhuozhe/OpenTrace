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
(1, 3, 'Incubation Hub', 'Room', ST_GeomFromText('POINT(77.4766105 28.4749348)', 4326)),
(2, 5, 'Reception', 'Room', ST_GeomFromText('POINT(77.4763316 28.4748098)', 4326)),
(3, 4, 'Library', 'Room', ST_GeomFromText('POINT(77.4766105 28.4749348)', 4326)),
(4, 7, 'B-04', 'Room', ST_GeomFromText('POINT(77.4760996 28.4747285)', 4326)),
(5, 7, 'Seminar Hall', 'Room', ST_GeomFromText('POINT(77.4759400 28.4746954)', 4326)),
(6, 7, 'Corridor Seminar Hall', 'Corridor', ST_GeomFromText('POINT(77.4759614 28.4746613)', 4326)),
(7, 1, 'A Block Ground Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4765287 28.4748334)', 4326)),
(8, 2, 'A Block First Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4765287 28.4748334)', 4326)),
(9, 3, 'A Block Second Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4765287 28.4748334)', 4326)),
(10, 4, 'A Block Third Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4765287 28.4748334)', 4326)),
(11, 1, 'A Block Ground Floor Elevator', 'Elevator', ST_GeomFromText('POINT(77.4764670 28.4748841)', 4326)),
(12, 2, 'A Block First Floor Elevator', 'Elevator', ST_GeomFromText('POINT(77.4764670 28.4748841)', 4326)),
(13, 3, 'A Block Second Floor Elevator', 'Elevator', ST_GeomFromText('POINT(77.4764670 28.4748841)', 4326)),
(14, 4, 'A Block Third Floor Elevator', 'Elevator', ST_GeomFromText('POINT(77.4764670 28.4748841)', 4326)),
(15, 7, 'B Block Ground Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4761827 28.4747320)', 4326)),
(16, 8, 'B Block First Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4761827 28.4747320)', 4326)),
(17, 9, 'B Block Second Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4761827 28.4747320)', 4326)),
(18, 10, 'B Block Third Floor Stairs', 'Stairs', ST_GeomFromText('POINT(77.4761827 28.4747320)', 4326));

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

-- connecting A block second floor elevator to Incubation Hub
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

-- connecting A block third floor elevator to Library
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    4,
    14,
    3,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 14)::geography, 
        (SELECT geom FROM nodes WHERE id = 3)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 14), 
        (SELECT geom FROM nodes WHERE id = 3)
    );

-- connecting seminar hall to seminar hall corridor
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    7,
    5,
    6,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 5)::geography, 
        (SELECT geom FROM nodes WHERE id = 6)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 5), 
        (SELECT geom FROM nodes WHERE id = 6)
    );

-- connecting A Block stairs 0-1 
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    1,
    7,
    8,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 7)::geography, 
        (SELECT geom FROM nodes WHERE id = 8)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 7), 
        (SELECT geom FROM nodes WHERE id = 8)
    );

-- connecting A Block stairs 1-2
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    2,
    8,
    9,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 8)::geography, 
        (SELECT geom FROM nodes WHERE id = 9)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 8), 
        (SELECT geom FROM nodes WHERE id = 9)
    );

-- connecting A Block stairs 2-3
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    3,
    9,
    10,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 9)::geography, 
        (SELECT geom FROM nodes WHERE id = 10)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 9), 
        (SELECT geom FROM nodes WHERE id = 10)
    );

-- connecting A Block 2nd Floor Stairs to Incubation Hub 
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    3,
    9,
    1,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 9)::geography, 
        (SELECT geom FROM nodes WHERE id = 1)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 9), 
        (SELECT geom FROM nodes WHERE id = 1)
    );

-- connecting A Block 3nd Floor Stairs to Library
INSERT INTO edges (level_id, source_node_id, target_node_id, distance_meters, geom)
SELECT 
    4,
    10,
    3,
    ST_Distance(
        (SELECT geom FROM nodes WHERE id = 10)::geography, 
        (SELECT geom FROM nodes WHERE id = 3)::geography
    ),
    ST_MakeLine(
        (SELECT geom FROM nodes WHERE id = 10), 
        (SELECT geom FROM nodes WHERE id = 3)
    );