import { useState, useEffect, useMemo, useCallback } from 'react'
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from 'react-leaflet'
import L from 'leaflet'

const API_BASE = ''

const categoryColors = {
  Room: '#4caf50',
  Corridor: '#ff9800',
  Stairs: '#9c27b0',
  Elevator: '#2196f3'
}

const floorColors = ['#e91e63', '#00bcd4', '#ff5722', '#607d8b', '#795548']
const routeColors = ['#f44336', '#2196f3', '#4caf50', '#ff9800', '#9c27b0']

const createNodeIcon = (color, label, size) => L.divIcon({
  className: 'custom-marker-node',
  html: `
    <div style="
      background:${color};
      width:${size}px;
      height:${size}px;
      border-radius:50%;
      border:3px solid white;
      box-shadow:0 3px 8px rgba(0,0,0,0.4);
      display:flex;
      align-items:center;
      justify-content:center;
      font-weight:bold;
      font-size:${Math.max(8, size/3.5)}px;
      color:white;
    ">${label}</div>
  `,
  iconSize: [size, size],
  iconAnchor: [size/2, size/2]
})

const createRouteNodeIcon = (size) => L.divIcon({
  className: 'route-node-marker',
  html: `
    <div style="
      background:#f44336;
      width:${size}px;
      height:${size}px;
      border-radius:50%;
      border:3px solid white;
      box-shadow:0 2px 6px rgba(0,0,0,0.4);
    "></div>
  `,
  iconSize: [size, size],
  iconAnchor: [size/2, size/2]
})

function MapController({ center, zoom }) {
  const map = useMap()
  useEffect(() => {
    if (center) map.setView(center, zoom)
  }, [center, zoom, map])
  return null
}

function calculateDirection(prev, curr, next) {
  if (!prev || !next) return 'start'
  
  const angle1 = Math.atan2(curr.lat - prev.lat, curr.lon - prev.lon) * 180 / Math.PI
  const angle2 = Math.atan2(next.lat - curr.lat, next.lon - curr.lon) * 180 / Math.PI
  let angleDiff = angle2 - angle1
  if (angleDiff > 180) angleDiff -= 360
  if (angleDiff < -180) angleDiff += 360
  
  if (Math.abs(angleDiff) < 30) return 'straight'
  if (angleDiff > 30 && angleDiff < 150) return 'right'
  if (angleDiff < -30 && angleDiff > -150) return 'left'
  if (Math.abs(angleDiff) > 150) return 'turnaround'
  return 'straight'
}

function App() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [levels, setLevels] = useState([])
  const [nodes, setNodes] = useState([])
  const [edges, setEdges] = useState([])
  const [selectedLevel, setSelectedLevel] = useState(null)
  const [selectedNode, setSelectedNode] = useState(null)
  const [startNode, setStartNode] = useState(null)
  const [endNode, setEndNode] = useState(null)
  const [routes, setRoutes] = useState([])
  const [selectedRouteIndex, setSelectedRouteIndex] = useState(0)
  const [showAllLevels, setShowAllLevels] = useState(true)
  const [mapCenter, setMapCenter] = useState([28.4747, 77.4765])
  const [mapZoom, setMapZoom] = useState(18)
  
  const [nodeSize, setNodeSize] = useState(32)
  const [showNodes, setShowNodes] = useState(true)
  const [showEdges, setShowEdges] = useState(true)
  const [showRouteNodes, setShowRouteNodes] = useState(true)
  const [alternateCount, setAlternateCount] = useState(1)
  const [selectedBuilding, setSelectedBuilding] = useState('')

  useEffect(() => {
    fetchDashboardData()
  }, [])

  const fetchDashboardData = async () => {
    try {
      setLoading(true)
      const res = await fetch(`${API_BASE}/dashboard`)
      if (!res.ok) throw new Error('Failed to fetch data')
      const data = await res.json()
      setLevels(data.levels)
      setNodes(data.nodes)
      setEdges(data.edges)
      if (data.levels.length > 0) {
        setSelectedBuilding(data.levels[0].building_name)
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const findRoutes = async () => {
    if (!startNode || !endNode) return
    try {
      const res = await fetch(`${API_BASE}/route/alternate?start=${startNode.id}&end=${endNode.id}&count=${alternateCount}`)
      if (!res.ok) throw new Error('Route not found')
      const data = await res.json()
      setRoutes(data.routes)
      setSelectedRouteIndex(0)
      if (data.routes.length > 0 && data.routes[0].path.length > 0) {
        setMapCenter([data.routes[0].path[0].lat, data.routes[0].path[0].lon])
        setMapZoom(18)
      }
    } catch (err) {
      alert('No route found between these nodes')
    }
  }

  const currentRoute = routes[selectedRouteIndex]

  const filteredNodes = useMemo(() => {
    if (showAllLevels) return nodes
    return selectedLevel ? nodes.filter(n => n.level_id === selectedLevel) : nodes
  }, [nodes, showAllLevels, selectedLevel])

  const filteredEdges = useMemo(() => {
    if (showAllLevels) return edges
    return selectedLevel ? edges.filter(e => e.level_id === selectedLevel) : edges
  }, [edges, showAllLevels, selectedLevel])

  const getEdgeCoords = (edge) => {
    const source = nodes.find(n => n.id === edge.source_node_id)
    const target = nodes.find(n => n.id === edge.target_node_id)
    if (source && target) {
      return [[source.lat, source.lon], [target.lat, target.lon]]
    }
    return null
  }

  const routePathByFloor = useMemo(() => {
    if (!currentRoute?.path) return {}
    const floors = {}
    currentRoute.path.forEach((node, idx) => {
      const floor = node.floor
      if (!floors[floor]) floors[floor] = []
      floors[floor].push({ ...node, index: idx })
    })
    return floors
  }, [currentRoute])

  const navigationSteps = useMemo(() => {
    if (!currentRoute?.path) return []
    const steps = []
    currentRoute.path.forEach((node, idx) => {
      const prev = currentRoute.path[idx - 1]
      const next = currentRoute.path[idx + 1]
      const distance = idx < currentRoute.path.length - 1 
        ? Math.sqrt(
            Math.pow(currentRoute.path[idx+1].lat - node.lat, 2) + 
            Math.pow(currentRoute.path[idx+1].lon - node.lon, 2)
          ) * 111000
        : 0
      steps.push({
        step: idx + 1,
        node,
        direction: calculateDirection(prev, node, next),
        distance
      })
    })
    return steps
  }, [currentRoute])

  const buildings = [...new Set(levels.map(l => l.building_name))]
  
  const buildingLevels = useMemo(() => {
    if (!selectedBuilding) return []
    return levels.filter(l => l.building_name === selectedBuilding)
  }, [levels, selectedBuilding])

  if (loading) return <div className="loading">Loading dashboard...</div>
  if (error) return <div className="error">Error: {error}</div>

  return (
    <div className="dashboard">
      <div className="sidebar">
        <h1>OpenTrace</h1>
        <p style={{fontSize: '12px', color: '#888', marginBottom: '15px'}}>Micro Navigation</p>
        
        <div className="stats">
          <div className="stat-card">
            <div className="value">{nodes.length}</div>
            <div className="label">Nodes</div>
          </div>
          <div className="stat-card">
            <div className="value">{levels.length}</div>
            <div className="label">Floors</div>
          </div>
        </div>

        <div className="settings-panel">
          <h3>Map Settings</h3>
          
          <div className="setting-item">
            <label>Node Size: {nodeSize}px</label>
            <input 
              type="range" 
              min="16" 
              max="64" 
              value={nodeSize}
              onChange={(e) => setNodeSize(Number(e.target.value))}
            />
          </div>

          <div className="layer-toggles">
            <label className="toggle-label">
              <input type="checkbox" checked={showNodes} onChange={(e) => setShowNodes(e.target.checked)} />
              Show Nodes
            </label>
            <label className="toggle-label">
              <input type="checkbox" checked={showEdges} onChange={(e) => setShowEdges(e.target.checked)} />
              Show Paths
            </label>
            <label className="toggle-label">
              <input type="checkbox" checked={showRouteNodes} onChange={(e) => setShowRouteNodes(e.target.checked)} />
              Show Route Stops
            </label>
            <label className="toggle-label">
              <input type="checkbox" checked={showAllLevels} onChange={(e) => setShowAllLevels(e.target.checked)} />
              All Floors
            </label>
          </div>
        </div>

        {!showAllLevels && (
          <>
            <div className="filter-group">
              <label>Building</label>
              <select 
                value={selectedBuilding}
                onChange={(e) => {
                  setSelectedBuilding(e.target.value)
                  const firstLevel = levels.find(l => l.building_name === e.target.value)
                  if (firstLevel) setSelectedLevel(firstLevel.id)
                }}
              >
                {buildings.map(b => (
                  <option key={b} value={b}>{b}</option>
                ))}
              </select>
            </div>

            <div className="filter-group">
              <label>Floor</label>
              <select 
                value={selectedLevel || ''}
                onChange={(e) => setSelectedLevel(Number(e.target.value))}
              >
                {buildingLevels.map(l => (
                  <option key={l.id} value={l.id}>{l.label}</option>
                ))}
              </select>
            </div>
          </>
        )}

        <h2>Route Planner</h2>
        <div className="filter-group">
          <label>Start Node</label>
          <select 
            value={startNode?.id || ''}
            onChange={(e) => {
              const node = nodes.find(n => n.id === Number(e.target.value))
              setStartNode(node)
            }}
          >
            <option value="">Select start</option>
            {nodes.map(n => (
              <option key={n.id} value={n.id}>{n.name}</option>
            ))}
          </select>
        </div>
        <div className="filter-group">
          <label>End Node</label>
          <select 
            value={endNode?.id || ''}
            onChange={(e) => {
              const node = nodes.find(n => n.id === Number(e.target.value))
              setEndNode(node)
            }}
          >
            <option value="">Select destination</option>
            {nodes.map(n => (
              <option key={n.id} value={n.id}>{n.name}</option>
            ))}
          </select>
        </div>
        
        <div className="filter-group">
          <label>Alternate Routes: {alternateCount}</label>
          <input 
            type="range" 
            min="1" 
            max="3" 
            value={alternateCount}
            onChange={(e) => setAlternateCount(Number(e.target.value))}
            style={{width: '100%'}}
          />
        </div>

        <button className="btn" onClick={findRoutes} disabled={!startNode || !endNode}>
          Find Routes
        </button>

        {routes.length > 0 && (
          <>
            <div className="route-selector">
              {routes.map((route, idx) => (
                <button 
                  key={idx}
                  className={`route-btn ${selectedRouteIndex === idx ? 'active' : ''}`}
                  onClick={() => setSelectedRouteIndex(idx)}
                  style={{borderColor: routeColors[idx % routeColors.length]}}
                >
                  Route {idx + 1} ({route.summary.total_distance_meters.toFixed(0)}m)
                </button>
              ))}
            </div>

            <div className="route-panel">
              <h3>Route {selectedRouteIndex + 1}</h3>
              <div className="route-info">
                <p><strong>From:</strong> {currentRoute.summary.start_point}</p>
                <p><strong>To:</strong> {currentRoute.summary.end_point}</p>
                <p><strong>Distance:</strong> {currentRoute.summary.total_distance_meters.toFixed(1)}m</p>
                <p><strong>Stops:</strong> {currentRoute.summary.node_count}</p>
              </div>

              <h4 style={{marginTop: '15px', fontSize: '14px'}}>Path by Floor:</h4>
              <div className="floor-layers">
                {Object.entries(routePathByFloor).map(([floor, pathNodes], idx) => (
                  <div 
                    key={floor} 
                    className="floor-layer-item"
                    onClick={() => {
                      if (pathNodes[0]) {
                        setMapCenter([pathNodes[0].lat, pathNodes[0].lon])
                        setMapZoom(18)
                      }
                    }}
                  >
                    <span className="floor-color" style={{background: floorColors[idx % floorColors.length]}}></span>
                    <span>{floor}</span>
                    <span style={{marginLeft: 'auto', fontSize: '12px'}}>{pathNodes.length} stops</span>
                  </div>
                ))}
              </div>

              <h4 style={{marginTop: '15px', fontSize: '14px'}}>Turn-by-Turn:</h4>
              <div className="nav-steps">
                {navigationSteps.map((step, idx) => (
                  <div key={idx} className="nav-step">
                    <div className="step-number">{step.step}</div>
                    <div className="step-content">
                      <div className="step-direction">
                        {step.direction === 'start' && '🚶 Start'}
                        {step.direction === 'straight' && '⬆️ Continue'}
                        {step.direction === 'left' && '⬅️ Turn left'}
                        {step.direction === 'right' && '➡️ Turn right'}
                        {step.direction === 'turnaround' && '🔄 Turn around'}
                      </div>
                      <div className="step-location">{step.node.name}</div>
                      <div className="step-floor">{step.node.floor}</div>
                    </div>
                    {step.distance > 0 && (
                      <div className="step-distance">{step.distance.toFixed(0)}m</div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </>
        )}

        <h2>Nodes ({filteredNodes.length})</h2>
        <div className="node-list">
          {filteredNodes.slice(0, 50).map(node => (
            <div 
              key={node.id} 
              className={`node-item ${selectedNode?.id === node.id ? 'selected' : ''}`}
              onClick={() => {
                setSelectedNode(node)
                setMapCenter([node.lat, node.lon])
                setMapZoom(19)
              }}
            >
              <div className="name">{node.name}</div>
              <div className="node-meta">
                <span className="category-tag" style={{background: categoryColors[node.category]}}>{node.category}</span>
                <span className="floor">{node.floor_label}</span>
              </div>
            </div>
          ))}
          {filteredNodes.length > 50 && (
            <p style={{fontSize: '12px', color: '#888', textAlign: 'center', padding: '10px'}}>
              +{filteredNodes.length - 50} more nodes
            </p>
          )}
        </div>
      </div>

      <div className="map-container">
        <MapContainer 
          center={mapCenter} 
          zoom={mapZoom} 
          maxZoom={19}
          minZoom={15}
          scrollWheelZoom={true}
          zoomControl={true}
        >
          <TileLayer
            attribution='&copy; OpenStreetMap'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            maxZoom={19}
          />
          <MapController center={mapCenter} zoom={mapZoom} />
          
          {showEdges && filteredEdges.map(edge => {
            const coords = getEdgeCoords(edge)
            if (!coords) return null
            const isRoute = currentRoute?.path?.some((_, i) => 
              currentRoute.path[i]?.id === edge.source_node_id && 
              currentRoute.path[i + 1]?.id === edge.target_node_id
            )
            return (
              <Polyline 
                key={edge.id} 
                positions={coords} 
                color={isRoute ? routeColors[selectedRouteIndex % routeColors.length] : '#ccc'}
                weight={isRoute ? 4 : 1}
                opacity={isRoute ? 1 : 0.3}
              />
            )
          })}

          {routes.map((route, routeIdx) => (
            route.path && (
              <Polyline 
                key={`route-${routeIdx}`}
                positions={route.path.map(n => [n.lat, n.lon])}
                color={routeColors[routeIdx % routeColors.length]}
                weight={4}
                opacity={routeIdx === selectedRouteIndex ? 1 : 0.4}
              />
            )
          ))}

          {showNodes && filteredNodes.map(node => {
            const iconLabel = node.category === 'Room' || node.category === 'Corridor' 
              ? node.name.substring(0, 2).toUpperCase() 
              : node.category === 'Stairs' ? 'S' : 'E'
            return (
              <Marker 
                key={node.id} 
                position={[node.lat, node.lon]}
                icon={createNodeIcon(categoryColors[node.category] || '#666', iconLabel, nodeSize)}
                eventHandlers={{
                  click: () => setSelectedNode(node)
                }}
              >
                <Popup>
                  <div className="node-popup">
                    <h4>{node.name}</h4>
                    <p><strong>Category:</strong> {node.category}</p>
                    <p><strong>Floor:</strong> {node.floor}</p>
                    <button 
                      className="btn" 
                      style={{marginTop: '10px', padding: '5px 10px', fontSize: '12px'}}
                      onClick={() => {
                        if (!startNode) setStartNode(node)
                        else setEndNode(node)
                      }}
                    >
                      {!startNode ? 'Set as Start' : 'Set as End'}
                    </button>
                  </div>
                </Popup>
              </Marker>
            )
          })}

          {showRouteNodes && currentRoute?.path && currentRoute.path.map((node, idx) => (
            <Marker
              key={`route-node-${idx}`}
              position={[node.lat, node.lon]}
              icon={createRouteNodeIcon(12)}
            >
              <Popup>
                <div style={{textAlign: 'center'}}>
                  <strong>Step {idx + 1}</strong><br/>
                  {node.name}
                </div>
              </Popup>
            </Marker>
          ))}
        </MapContainer>

        <div className="legend">
          <h4>Categories</h4>
          {Object.entries(categoryColors).map(([cat, color]) => (
            <div key={cat} className="legend-item">
              <div className="legend-color" style={{background: color}}></div>
              {cat}
            </div>
          ))}
          {routes.length > 0 && (
            <>
              <h4 style={{marginTop: '10px'}}>Routes</h4>
              {routes.map((_, idx) => (
                <div key={idx} className="legend-item">
                  <div className="legend-color" style={{background: routeColors[idx % routeColors.length]}}></div>
                  Route {idx + 1}
                </div>
              ))}
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default App
