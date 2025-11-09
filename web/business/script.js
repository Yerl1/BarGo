const API_BASE = "http://localhost:3003";
let currentStoreId = null;

/* ---------- STORES ---------- */

// Load all stores
async function loadStores() {
  const token = localStorage.getItem("jwt");

  try {
    const res = await fetch(`${API_BASE}/stores`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    });

    if (!res.ok) throw new Error("Failed to load stores");
    const stores = await res.json();
    renderStores(stores);
  } catch (err) {
    console.error(err);
    alert("Ошибка при загрузке магазинов");
  }
}

// Render stores grid
function renderStores(stores) {
  const grid = document.getElementById("storesGrid");
  grid.innerHTML = "";

  if (!stores.length) {
    grid.innerHTML = `<p class="text-center text-muted">Нет магазинов</p>`;
    return;
  }

  stores.forEach((store) => {
    const card = document.createElement("div");
    card.className = "col";
    card.innerHTML = `
      <div class="card h-100 shadow-sm">
        <img src="${
          store.photo || "https://via.placeholder.com/400"
        }" class="card-img-top" alt="Фото магазина">
        <div class="card-body">
          <h5 class="card-title">${store.name}</h5>
          <p class="card-text text-muted">${store.description || ""}</p>
          <p class="small text-secondary">${store.address || ""}</p>
        </div>
        <div class="card-footer d-flex justify-content-between">
          <button class="btn btn-sm btn-outline-primary" onclick="openProducts('${
            store.store_id
          }', '${store.name}')">
            <i class="bi bi-box"></i> Товары
          </button>
          <div>
            <button class="btn btn-sm btn-outline-warning me-1" onclick="openUpdateStoreModal(${JSON.stringify(
              store
            ).replace(/"/g, "&quot;")})">
              <i class="bi bi-pencil"></i>
            </button>
            <button class="btn btn-sm btn-outline-danger" onclick="deleteStore('${
              store.store_id
            }')">
              <i class="bi bi-trash"></i>
            </button>
          </div>
        </div>
      </div>`;
    grid.appendChild(card);
  });
}

// Add new store
async function addStore() {
  const token = localStorage.getItem("jwt");
  if (!token) {
    alert("You must be logged in to add a store.");
    return;
  }

  const data = {
    name: document.getElementById("storeName").value,
    address: document.getElementById("storeAddress").value,
    latitude: parseFloat(document.getElementById("storeLat").value),
    longitude: parseFloat(document.getElementById("storeLon").value),
    description: document.getElementById("storeDescription").value,
    photo: document.getElementById("storePhoto").value,
  };

  try {
    const res = await fetch(`${API_BASE}/stores`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to add store");
    await res.json();
    bootstrap.Modal.getInstance(
      document.getElementById("addStoreModal")
    ).hide();
    await loadStores();
  } catch (err) {
    console.error(err);
    alert("Ошибка при добавлении магазина");
  }
}

// Update store
async function updateStore() {
  const token = localStorage.getItem("jwt");
  if (!token) {
    alert("You must be logged in to add a store.");
    return;
  }

  const data = {
    store_id: document.getElementById("updateStoreId").value,
    name: document.getElementById("updateStoreName").value,
    address: document.getElementById("updateStoreAddress").value,
    latitude: parseFloat(document.getElementById("updateStoreLat").value),
    longitude: parseFloat(document.getElementById("updateStoreLon").value),
    description: document.getElementById("updateStoreDescription").value,
    photo: document.getElementById("updateStorePhoto").value,
  };

  try {
    const res = await fetch(`${API_BASE}/stores`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to update store");
    bootstrap.Modal.getInstance(
      document.getElementById("updateStoreModal")
    ).hide();
    await loadStores();
  } catch (err) {
    console.error(err);
    alert("Ошибка при обновлении магазина");
  }
}

// Delete store
async function deleteStore(id) {
  const token = localStorage.getItem("jwt");
  if (!token) {
    alert("You must be logged in to add a store.");
    return;
  }

  if (!confirm("Удалить этот магазин?")) return;
  try {
    const res = await fetch(`${API_BASE}/stores?store_id=${id}`, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      method: "DELETE",
    });
    if (!res.ok) throw new Error("Failed to delete store");
    await loadStores();
  } catch (err) {
    console.error(err);
    alert("Ошибка при удалении магазина");
  }
}

function openUpdateStoreModal(store) {
  document.getElementById("updateStoreId").value = store.store_id;
  document.getElementById("updateStoreName").value = store.name;
  document.getElementById("updateStoreAddress").value = store.address;
  document.getElementById("updateStoreLat").value = store.latitude;
  document.getElementById("updateStoreLon").value = store.longitude;
  document.getElementById("updateStoreDescription").value = store.description;
  document.getElementById("updateStorePhoto").value = store.photo;

  const modal = new bootstrap.Modal(
    document.getElementById("updateStoreModal")
  );
  modal.show();
}

/* ---------- PRODUCTS ---------- */

function openUpdateProductModal(product) {
  document.getElementById("updateProductId").value = product.product_id;
  document.getElementById("updateProductName").value = product.name;
  document.getElementById("updateProductCategory").value =
    product.category || "";
  document.getElementById("updateProductPrice").value = product.price;
  document.getElementById("updateProductDescription").value =
    product.description;
  document.getElementById("updateProductPhoto").value = product.photo;

  const modal = new bootstrap.Modal(
    document.getElementById("updateProductModal")
  );
  modal.show();
}

async function openProducts(storeId, storeName) {
  currentStoreId = storeId;
  document.getElementById("currentStoreName").textContent = storeName;
  document.getElementById("storesView").classList.add("hidden");
  document.getElementById("productsView").classList.remove("hidden");
  await loadProducts();
}

function backToStores() {
  document.getElementById("storesView").classList.remove("hidden");
  document.getElementById("productsView").classList.add("hidden");
  currentStoreId = null;
}

// Load products for current store
async function loadProducts() {
  const token = localStorage.getItem("jwt");
  try {
    const res = await fetch(`${API_BASE}/products?store_id=${currentStoreId}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    });

    if (!res.ok) throw new Error("Failed to load products");
    const products = await res.json();
    renderProducts(products);
  } catch (err) {
    console.error(err);
    alert("Ошибка при загрузке продуктов");
  }
}

function renderProducts(products) {
  const grid = document.getElementById("productsGrid");
  grid.innerHTML = "";

  if (!products.length) {
    grid.innerHTML = `<p class="text-center text-muted">Нет продуктов</p>`;
    return;
  }

  products.forEach((prod) => {
    const card = document.createElement("div");
    card.className = "col";
    card.innerHTML = `
      <div class="card h-100 shadow-sm">
        <img src="${
          prod.photo || "https://via.placeholder.com/400"
        }" class="card-img-top">
        <div class="card-body">
          <h5 class="card-title">${prod.name}</h5>
          <p class="text-muted">${prod.description || ""}</p>
          <p><strong>${prod.price} ₸</strong></p>
        </div>
        <div class="card-footer d-flex justify-content-end">
          <button class="btn btn-sm btn-outline-warning me-1" onclick="openUpdateProductModal(${JSON.stringify(
            prod
          ).replace(/"/g, "&quot;")})">
            <i class="bi bi-pencil"></i>
          </button>
          <button class="btn btn-sm btn-outline-danger" onclick="deleteProduct('${
            prod.product_id
          }')">
            <i class="bi bi-trash"></i>
          </button>
        </div>
      </div>`;
    grid.appendChild(card);
  });
}

// Add product
async function addProduct() {
  const token = localStorage.getItem("jwt");
  if (!token) {
    alert("You must be logged in to add a store.");
    return;
  }

  const data = {
    name: document.getElementById("productName").value,
    description: document.getElementById("productDescription").value,
    photo: document.getElementById("productPhoto").value,
    price: parseFloat(document.getElementById("productPrice").value),
  };

  try {
    const res = await fetch(`${API_BASE}/products?store_id=${currentStoreId}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to add product");
    bootstrap.Modal.getInstance(
      document.getElementById("addProductModal")
    ).hide();
    await loadProducts();
  } catch (err) {
    console.error(err);
    alert("Ошибка при добавлении продукта");
  }
}

// Update product
async function updateProduct() {
  const token = localStorage.getItem("jwt");
  if (!token) {
    alert("You must be logged in to add a store.");
    return;
  }

  const data = {
    product_id: document.getElementById("updateProductId").value,
    name: document.getElementById("updateProductName").value,
    description: document.getElementById("updateProductDescription").value,
    photo: document.getElementById("updateProductPhoto").value,
    price: parseFloat(document.getElementById("updateProductPrice").value),
  };

  try {
    const res = await fetch(`${API_BASE}/products`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to update product");
    bootstrap.Modal.getInstance(
      document.getElementById("updateProductModal")
    ).hide();
    await loadProducts();
  } catch (err) {
    console.error(err);
    alert("Ошибка при обновлении продукта");
  }
}

// Delete product
async function deleteProduct(id) {
  const token = localStorage.getItem("jwt");
  if (!token) {
    alert("You must be logged in to add a store.");
    return;
  }

  if (!confirm("Удалить этот продукт?")) return;
  try {
    const res = await fetch(`${API_BASE}/products?product_id=${id}`, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      method: "DELETE",
    });
    if (!res.ok) throw new Error("Failed to delete product");
    await loadProducts();
  } catch (err) {
    console.error(err);
    alert("Ошибка при удалении продукта");
  }
}

/* ---------- INIT ---------- */

document.addEventListener("DOMContentLoaded", loadStores);
