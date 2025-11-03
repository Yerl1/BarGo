document.addEventListener("DOMContentLoaded", () => {
  const map = L.map("map").setView([43.238949, 76.889709], 12); // Fallback center (e.g. Almaty)

  // Add OpenStreetMap base layer
  // Reference your map

  L.tileLayer("http://localhost:8081/styles/osm-bright/{z}/{x}/{y}.png", {
    maxZoom: 18,
    attribution: "&copy; OpenStreetMap contributors",
  }).addTo(map);

  // Invalidate size when the map tab is shown
  const mapTabBtn = document.getElementById("map-tab");
  mapTabBtn.addEventListener("shown.bs.tab", () => {
    map.invalidateSize(); // Recalculate Leaflet map dimensions
  });

  // Dropdown to select routing mode
  const modeSelect = L.control({ position: "topright" });
  modeSelect.onAdd = () => {
    const div = L.DomUtil.create("div", "leaflet-bar leaflet-control");
    div.innerHTML = `
      <select id="modeSelect" style="padding:4px;border:none;outline:none;">
        <option value="car">🚗 Car</option>
        <option value="bike">🚴 Bicycle</option>
        <option value="foot">🚶 Foot</option>
      </select>
    `;
    return div;
  };
  modeSelect.addTo(map);

  // Default routing variables
  let start = null;
  const end = [43.250447, 76.948394]; // Example destination (store)
  let control = null;

  // Initialize route control function
  function createRoute(mode) {
    const port = mode === "car" ? 5000 : mode === "bike" ? 5001 : 5002;
    if (control) map.removeControl(control);

    control = L.Routing.control({
      waypoints: [L.latLng(start[0], start[1]), L.latLng(end[0], end[1])],
      router: L.Routing.osrmv1({
        serviceUrl: `http://localhost:${port}/route/v1`,
      }),
      lineOptions: { styles: [{ color: "#007bff", weight: 5 }] },
      createMarker: (i, wp, nWps) =>
        L.marker(wp.latLng, { draggable: i === 0 }), // only start is draggable
    }).addTo(map);
  }

  // Handle mode switching
  document.body.addEventListener("change", (e) => {
    if (e.target.id === "modeSelect" && start) {
      createRoute(e.target.value);
    }
  });

  // Try to get user location
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const lat = pos.coords.latitude;
        const lon = pos.coords.longitude;
        start = [lat, lon];
        map.setView(start, 14);

        // Add a marker for user's position
        L.marker(start).addTo(map).bindPopup("📍 You are here").openPopup();

        // Create initial car route
        createRoute("car");
      },
      (err) => {
        console.error("Geolocation failed:", err.message);
        alert("Could not get your location. Using default center.");
        start = [43.238949, 76.889709]; // fallback
        createRoute("car");
      }
    );
  } else {
    alert("Geolocation not supported by this browser.");
    start = [43.238949, 76.889709];
    createRoute("car");
  }
});
