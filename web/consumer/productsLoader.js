// productsLoader.js

document.addEventListener("DOMContentLoaded", () => {
  loadProducts();

  // Optional: handle search and sort after loading products
  const searchForm = document.getElementById("searchForm");
  searchForm.addEventListener("submit", (e) => {
    e.preventDefault();
    loadProducts(document.getElementById("q").value.trim());
  });

  const sortSelect = document.getElementById("sortSelect");
  sortSelect.addEventListener("change", () => {
    loadProducts(document.getElementById("q").value.trim(), sortSelect.value);
  });
});

async function loadProducts(query = "", sort = "relevance") {
  const grid = document.getElementById("cardsGrid");
  grid.innerHTML = ""; // Clear previous cards
  const emptyState = document.getElementById("emptyState");
  const foundCount = document.getElementById("foundCount");

  try {
    let url = "http://localhost:3002/products";
    if (query) url += `?q=${encodeURIComponent(query)}`;
    const res = await fetch(url);
    if (!res.ok) throw new Error("Failed to fetch products");

    let products = await res.json();

    // Simple frontend filtering based on query keywords if backend doesn't support it
    if (query) {
      const qLower = query.toLowerCase();
      products = products.filter(
        (p) =>
          p.name.toLowerCase().includes(qLower) ||
          (p.description && p.description.toLowerCase().includes(qLower))
      );
    }

    // Sorting
    products.sort((a, b) => {
      switch (sort) {
        case "price_asc":
          return a.price - b.price;
        case "price_desc":
          return b.price - a.price;
        default:
          return 0; // relevance or unsupported: leave as-is
      }
    });

    // Show empty state if no products
    if (products.length === 0) {
      emptyState.classList.remove("d-none");
      document.getElementById("emptyQuery").textContent = query || "...";
      foundCount.textContent = 0;
      return;
    } else {
      emptyState.classList.add("d-none");
      foundCount.textContent = products.length;
    }

    products.forEach((p) => {
      const col = document.createElement("div");
      col.className = "col prod-card";
      col.dataset.title = p.name;
      col.dataset.price = p.price;
      col.innerHTML = `
        <div class="card h-100 product-card">
          <div class="ratio ratio-16x9">
            <img src="${
              p.photo ?? "placeholder.jpg"
            }" class="card-img-top img-cover" alt="${p.name}">
          </div>
          <div class="card-body d-flex flex-column">
            <h6 class="card-title js-title mb-1">${p.name}</h6>
            <div class="text-muted small mb-2">${p.description ?? ""}</div>
            <div class="mt-auto d-flex align-items-center justify-content-between">
              <div>
                <div class="fw-semibold text-primary-emphasis">${p.price.toLocaleString()} ₸</div>
              </div>
              <a href="product.html?id=${
                p.product_id
              }" class="btn btn-outline-primary btn-sm">Смотреть</a>
            </div>
          </div>
        </div>
      `;
      grid.appendChild(col);
    });
  } catch (err) {
    console.error(err);
    grid.innerHTML = `<div class="col-12 text-center text-danger">Ошибка загрузки продуктов</div>`;
  }
}
