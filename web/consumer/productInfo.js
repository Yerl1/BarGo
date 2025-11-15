// productInfo.js
document.addEventListener("DOMContentLoaded", () => {
  initProductDetails();
  initStores();
  initStoreFilters();
  initAddToCartButtons();
});

// ------ Product details from backend (optional) ------

async function initProductDetails() {
  const urlParams = new URLSearchParams(window.location.search);
  const productId = urlParams.get("id");
  if (!productId) return;

  try {
    const resp = await fetch(
      `http://localhost:3002/product-info?product_id=${encodeURIComponent(
        productId
      )}`
    );
    if (!resp.ok)
      throw new Error(`Failed to load product info: ${resp.status}`);

    const data = await resp.json();
    if (!data || !data.product) return;

    const product = data.product;

    const imgEl = document.getElementById("productImage");
    if (imgEl && product.photo) {
      imgEl.src = product.photo;
      imgEl.alt = product.name || "Изображение товара";
    }

    if (product.name)
      document.getElementById("productTitle").textContent = product.name;
    if (typeof product.price === "number")
      document.getElementById(
        "productPrice"
      ).textContent = `${product.price} ₸`;
    if (product.description)
      document.getElementById("productDescription").textContent =
        product.description;
  } catch (err) {
    console.error("Error loading product info:", err);
  }
}

// ------ Stores list / filters ------

let allShops = [];

function initStores() {
  const list = document.getElementById("shopsList");
  if (!list) return;

  const cards = Array.from(list.querySelectorAll(".store-card"));
  allShops = cards.map((card, index) => {
    const nameEl = card.querySelector(".store-name");
    const addrEl = card.querySelector(".store-address");
    const priceEl = card.querySelector(".store-price");
    const distEl = card.querySelector(".store-distance");

    const name = nameEl ? nameEl.textContent.trim() : "";
    const address = addrEl ? addrEl.textContent.trim() : "";

    const priceRaw = priceEl ? priceEl.textContent : "";
    const distanceRaw = distEl ? distEl.textContent : "";

    const price = parseFloat(priceRaw.replace(/[^\d.]/g, "")) || 0;
    const distance =
      parseFloat(distanceRaw.replace(",", ".").replace(/[^\d.]/g, "")) ||
      Infinity;

    card.dataset.initialIndex = index;
    return {
      element: card,
      name,
      address,
      price,
      distance,
    };
  });

  applyStoreFilters(); // первый рендер
}

function renderShopsFromApi(shops) {
  const list = document.getElementById("shopsList");
  if (!list) return;

  list.innerHTML = ""; // убираем заглушки

  allShops = shops.map((shop, index) => {
    const card = createShopCard(shop);
    card.dataset.initialIndex = index;
    list.appendChild(card);
    return {
      element: card,
      name: shop.name || "",
      address: shop.address || "",
      price: shop.price || 0,
      distance: shop.distance_km ?? Infinity,
    };
  });

  applyStoreFilters();
}

function createShopCard(shop) {
  const card = document.createElement("div");
  card.className = "store-card";

  const price = shop.price ?? 0;
  const distance = shop.distance_km ?? null;

  card.innerHTML = `
      <div class="row align-items-center">
        <div class="col-md-6">
          <div class="store-name">${shop.name || "Магазин"}</div>
          <div class="store-address">
            <i class="bi bi-geo-alt me-1"></i> ${shop.address || ""}
          </div>
        </div>
        <div class="col-md-3">
          <div class="store-price">${price} ₸</div>
          ${
            distance != null
              ? `<div class="store-distance">${distance} км</div>`
              : ""
          }
        </div>
        <div class="col-md-3">
          <div class="store-actions">
            <button class="btn btn-outline btn-sm me-2 js-route-btn">
              <i class="bi bi-signpost me-1"></i> Маршрут
            </button>
            <button class="btn btn-primary btn-sm js-add-to-cart">
              <i class="bi bi-cart me-1"></i> В корзину
            </button>
          </div>
        </div>
      </div>
    `;

  return card;
}

function initStoreFilters() {
  const shopQuery = document.getElementById("shopQuery");
  const minPrice = document.getElementById("minPrice");
  const maxPrice = document.getElementById("maxPrice");
  const maxDist = document.getElementById("maxDist");
  const sortShops = document.getElementById("sortShops");

  const inputs = [shopQuery, minPrice, maxPrice, maxDist];
  inputs.forEach((el) => {
    if (!el) return;
    el.addEventListener("input", () => applyStoreFilters());
  });
  if (sortShops) {
    sortShops.addEventListener("change", () => applyStoreFilters());
  }
}

function applyStoreFilters() {
  const list = document.getElementById("shopsList");
  if (!list || !allShops.length) return;

  const shopQuery = document.getElementById("shopQuery");
  const minPrice = document.getElementById("minPrice");
  const maxPrice = document.getElementById("maxPrice");
  const maxDist = document.getElementById("maxDist");
  const sortShops = document.getElementById("sortShops");

  const q = shopQuery ? shopQuery.value.trim().toLowerCase() : "";
  const min = minPrice && minPrice.value ? Number(minPrice.value) : 0;
  const max = maxPrice && maxPrice.value ? Number(maxPrice.value) : Infinity;
  const maxDistance =
    maxDist && maxDist.value ? Number(maxDist.value) : Infinity;
  const sort = sortShops ? sortShops.value : "distance_asc";

  let filtered = allShops.filter((shop) => {
    if (
      q &&
      !shop.name.toLowerCase().includes(q) &&
      !shop.address.toLowerCase().includes(q)
    ) {
      return false;
    }
    if (min && shop.price < min) return false;
    if (max !== Infinity && shop.price > max) return false;
    if (maxDistance !== Infinity && shop.distance > maxDistance) return false;
    return true;
  });

  const cmp = {
    price_asc: (a, b) => a.price - b.price,
    price_desc: (a, b) => b.price - a.price,
    distance_asc: (a, b) => a.distance - b.distance,
    distance_desc: (a, b) => b.distance - a.distance,
  }[sort];

  if (cmp) {
    filtered.sort(cmp);
  } else {
    filtered.sort(
      (a, b) =>
        (a.distance || Infinity) - (b.distance || Infinity) || a.price - b.price
    );
  }

  list.innerHTML = "";

  if (!filtered.length) {
    list.innerHTML =
      '<div class="text-muted text-center py-4">Магазины по заданным параметрам не найдены</div>';
    return;
  }

  filtered.forEach((shop) => list.appendChild(shop.element));
}

// ------ Add-to-cart feedback ------

function initAddToCartButtons() {
  document.addEventListener("click", (e) => {
    const btn = e.target.closest(".js-add-to-cart");
    if (!btn) return;

    const originalHTML = btn.innerHTML;
    btn.innerHTML = '<i class="bi bi-check me-2"></i> Добавлено';
    btn.classList.add("btn-success");

    setTimeout(() => {
      btn.innerHTML = originalHTML;
      btn.classList.remove("btn-success");
    }, 2000);
  });
}
