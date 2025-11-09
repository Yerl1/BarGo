document.addEventListener("DOMContentLoaded", () => {
  // Center on Astana
  const ASTANA_COORDS = [51.1605, 71.4704];
  const map = L.map("map").setView(ASTANA_COORDS, 13);

  // Base tiles from your local TileServer GL
  L.tileLayer("http://localhost:8081/styles/osm-bright/{z}/{x}/{y}.png", {
    maxZoom: 18,
    attribution: "&copy; OpenStreetMap contributors",
  }).addTo(map);

  // Fix Leaflet sizing when switching to map tab
  const mapTabBtn = document.getElementById("map-tab");
  if (mapTabBtn) {
    mapTabBtn.addEventListener("shown.bs.tab", () => map.invalidateSize());
  }

  // Mode selector (Car / Bike / Foot)
  const modeSelect = L.control({ position: "topright" });
  modeSelect.onAdd = () => {
    const div = L.DomUtil.create("div", "leaflet-bar leaflet-control");
    div.innerHTML = `
      <select id="modeSelect" style="padding:4px;border:none;outline:none;">
        <option value="car" selected>🚗 Car</option>
        <option value="bike">🚴 Bicycle</option>
        <option value="foot">🚶 Foot</option>
      </select>
    `;
    return div;
  };
  modeSelect.addTo(map);

  // Radius selector (1 km – 10 km)
  const radiusControl = L.control({ position: "topright" });
  radiusControl.onAdd = () => {
    const div = L.DomUtil.create("div", "leaflet-bar leaflet-control");
    div.innerHTML = `
      <select id="radiusSelect" style="padding:4px;border:none;outline:none;">
        <option value="1000">1 km</option>
        <option value="2000">2 km</option>
        <option value="3000" selected>3 km</option>
        <option value="5000">5 km</option>
        <option value="10000">10 km</option>
      </select>
    `;
    return div;
  };
  radiusControl.addTo(map);

  let start = null;
  let control = null;
  const storeMarkers = L.layerGroup().addTo(map);
  let radiusCircle = null;

  // --- ROUTING FUNCTION ---
  function createRoute(dest, mode) {
    if (!start || !dest) return;

    const port = mode === "car" ? 5000 : mode === "bike" ? 5001 : 5002;
    if (control) map.removeControl(control);

    control = L.Routing.control({
      waypoints: [L.latLng(start), L.latLng(dest)],
      router: L.Routing.osrmv1({
        serviceUrl: `http://localhost:${port}/route/v1`,
      }),
      lineOptions: { styles: [{ color: "#007bff", weight: 5 }] },
      createMarker: (i, wp) => L.marker(wp.latLng, { draggable: i === 0 }),
    }).addTo(map);
  }

  // --- FETCH AND DISPLAY NEARBY STORES ---
  function fetchNearbyStores(lat, lon, radius = 1000) {
    fetch(`http://localhost:3002/stores?lat=${lat}&lon=${lon}&radius=${radius}`)
      .then((res) => res.json())
      .then((stores) => {
        storeMarkers.clearLayers(); // clear old markers

        if (radiusCircle) map.removeLayer(radiusCircle); // remove old circle
        radiusCircle = L.circle([lat, lon], {
          radius,
          color: "#007bff",
          fill: false,
        }).addTo(map);

        stores.forEach((store) => {
          if (store.latitude && store.longitude) {
            const marker = L.marker([store.latitude, store.longitude]).addTo(
              storeMarkers
            ).bindPopup(`
                <b>${store.name}</b><br>
                ${store.address || ""}<br>
                <button class="route-btn" data-lat="${
                  store.latitude
                }" data-lon="${store.longitude}">
                  📍 Route here
                </button>
              `);
          }
        });
      })
      .catch((err) => console.error("Error fetching stores:", err));
  }

  // --- HANDLE USER LOCATION ---
  function handleLocation(lat, lon) {
    start = [lat, lon];
    map.setView(start, 14);
    L.marker(start).addTo(map).bindPopup("📍 You are here").openPopup();

    const radius = parseInt(document.getElementById("radiusSelect").value, 10);
    fetchNearbyStores(lat, lon, radius);
  }

  if ("geolocation" in navigator) {
    navigator.geolocation.getCurrentPosition(
      (pos) => handleLocation(pos.coords.latitude, pos.coords.longitude),
      (err) => {
        console.warn("Geolocation failed:", err.message);
        alert("Could not get location. Defaulting to Astana center.");
        handleLocation(ASTANA_COORDS[0], ASTANA_COORDS[1]);
      }
    );
  } else {
    alert("Geolocation not supported, defaulting to Astana.");
    handleLocation(ASTANA_COORDS[0], ASTANA_COORDS[1]);
  }

  // --- HANDLE RADIUS CHANGE ---
  document.body.addEventListener("change", (e) => {
    if (e.target.id === "radiusSelect" && start) {
      const r = parseInt(e.target.value, 10);
      fetchNearbyStores(start[0], start[1], r);
    }
  });

  // --- ROUTE BUTTON CLICK HANDLER ---
  map.on("popupopen", (e) => {
    const btn = e.popup._contentNode.querySelector(".route-btn");
    if (btn) {
      btn.addEventListener("click", () => {
        const lat = parseFloat(btn.dataset.lat);
        const lon = parseFloat(btn.dataset.lon);
        const mode = document.getElementById("modeSelect").value;
        createRoute([lat, lon], mode);
      });
    }
  });
});
