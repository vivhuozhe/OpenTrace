# OpenTrace
A community-driven map app that makes micro-navigation possible. Our goal is to provide offline navigation inside public infrastructures like government hospitals, courts etc.

# Directory Structure
```
/OpenTrace
├── /apps
│   ├── /web-dashboard      # React.js (Vite + Tailwind)
│       ├── /src
│       │   ├── /components # Map components (MapLibre), UI kit
│       │   ├── /hooks      # Custom hooks for GeoJSON fetching
│       │   └── /services   # API client for the Go backend
│       └── tailwind.config.js
├── /backend
│   ├── /cmd
│   │   └── /api            # Entry point: main.go
│   ├── /internal           # Private business logic
│   │   ├── /api            # Gin handlers (routes.go, controllers/)
│   │   ├── /graph          # Pathfinding logic (A*, Dijkstra)
│   │   ├── /spatial        # PostGIS interaction layer
│   │   └── /models         # Go structs (Point, Edge, FloorPlan)
│   ├── /pkg                # Public libraries (e.g., GeoJSON utils)
│   ├── go.mod              # Module dependencies
│   └── go.sum
├── /infrastructure         # Devops & Database setup
│   ├── /db                 # SQL migrations and seed data
│   │   └── /migrations
│   │       ├── 001_init_postgis.sql
│   │       └── 002_create_map_tables.sql
│   └── docker-compose.yml  # Spins up Postgres/PostGIS
├── /docs                   # API documentation (OpenAPI/Swagger)
└── README.md
```
# Features
-  Shows the interior maps of buildings like university campuses, govt. hospitals, courts etc.
-  An easy-to-use map editor for users to update the interior maps.
-  Offline navigation

# Vision and principles
OpenTrace envisions a future where easy navigation inside large buildings is hassle free for everyone while being privacy-focused, accessible, easy-to-use and transparent. 

## Our principles
- Open Source: The complete codebase is open source and released under a Free and Open source license.
- Privacy focused: No data collection from users, no trackers, no ads etc.
- Community first: Fueled primarily by user-input and feedback.

# Tech Stack
- Backend -     Go + Gin
- Frontend -    HTML + Tailwind CSS + Reactjs
- Spatial DB -  Postgres + PostGIS
- Map engine -  MapLibre
