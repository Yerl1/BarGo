// Координаты центра Астаны доступны везде
const ASTANA_COORDS = [51.1605, 71.4704];

document.addEventListener("DOMContentLoaded", () => {
  const mapContainer = document.getElementById("map");
  if (!mapContainer) {
    console.warn("Map container #map not found");
    return;
  }

  // Инициализация карты
  const map = L.map("map").setView(ASTANA_COORDS, 13);
  window._leafletMap = map; // чтобы можно было дергать invalidateSize() из других файлов

  // Тайлы (локальный TileServer GL)
  L.tileLayer("http://localhost:8081/styles/osm-bright/{z}/{x}/{y}.png", {
    maxZoom: 18,
    attribution: "&copy; OpenStreetMap contributors",
  }).addTo(map);

  // Когда переключаемся на вкладку "Карта" — чинить размер
  const mapViewBtn = document.querySelector('.view-tab[data-view="map"]');
  if (mapViewBtn) {
    mapViewBtn.addEventListener("click", () => {
      setTimeout(() => {
        map.invalidateSize();
      }, 200);
    });
  }

  // --- КОНТРОЛ: выбор режима (авто / велик / пешком) ---
  const modeSelectControl = L.control({ position: "topright" });
  modeSelectControl.onAdd = () => {
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
  modeSelectControl.addTo(map);

  // --- КОНТРОЛ: радиус поиска (1–10 км) ---
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

  let start = null;                // точка старта (пользователь)
  let control = null;              // маршрут
  const storeMarkers = L.layerGroup().addTo(map);
  let radiusCircle = null;         // круг радиуса

  // --- ПОСТРОЕНИЕ МАРШРУТА ---
  function createRoute(dest, mode) {
    if (!start || !dest) return;

    const port = mode === "car" ? 5000 : mode === "bike" ? 5001 : 5002;

    if (control) {
      map.removeControl(control);
    }

    control = L.Routing.control({
      waypoints: [L.latLng(start), L.latLng(dest)],
      router: L.Routing.osrmv1({
        serviceUrl: `http://localhost:${port}/route/v1`,
      }),
      lineOptions: { styles: [{ color: "#007bff", weight: 5 }] },
      createMarker: (i, wp) =>
        L.marker(wp.latLng, { draggable: i === 0 }),
    }).addTo(map);
  }

  // --- ЗАГРУЗКА МАГАЗИНОВ В РАДИУСЕ ---
  function fetchNearbyStores(lat, lon, radius = 1000) {
    fetch(`http://localhost:3002/stores?lat=${lat}&lon=${lon}&radius=${radius}`)
      .then((res) => res.json())
      .then((stores) => {
        storeMarkers.clearLayers();

        if (radiusCircle) {
          map.removeLayer(radiusCircle);
        }
        radiusCircle = L.circle([lat, lon], {
          radius,
          color: "#007bff",
          fill: false,
        }).addTo(map);

        stores.forEach((store) => {
          if (!store.latitude || !store.longitude) return;

          const marker = L.marker([store.latitude, store.longitude]).addTo(storeMarkers);
          marker.bindPopup(`
            <b>${store.name}</b><br>
            ${store.address || ""}<br>
            <button class="route-btn"
                    data-lat="${store.latitude}"
                    data-lon="${store.longitude}">
              📍 Route here
            </button>
          `);
        });
      })
      .catch((err) => console.error("Error fetching stores:", err));
  }

  // --- ОБРАБОТКА ПОЛОЖЕНИЯ ПОЛЬЗОВАТЕЛЯ ---
  function handleLocation(lat, lon) {
    start = [lat, lon];
    map.setView(start, 14);

    L.marker(start).addTo(map).bindPopup("📍 You are here").openPopup();

    const radiusSelect = document.getElementById("radiusSelect");
    const radius = radiusSelect ? parseInt(radiusSelect.value, 10) : 3000;
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

  // --- СМЕНА РАДИУСА ---
  document.body.addEventListener("change", (e) => {
    if (e.target.id === "radiusSelect" && start) {
      const r = parseInt(e.target.value, 10);
      fetchNearbyStores(start[0], start[1], r);
    }
  });

  // --- КНОПКА "Route here" В ПОПАПЕ ---
  map.on("popupopen", (e) => {
    const btn = e.popup._contentNode.querySelector(".route-btn");
    if (!btn) return;

    btn.addEventListener("click", () => {
      const lat = parseFloat(btn.dataset.lat);
      const lon = parseFloat(btn.dataset.lon);
      const modeSelect = document.getElementById("modeSelect");
      const mode = modeSelect ? modeSelect.value : "car";
      createRoute([lat, lon], mode);
    });
  });
});
