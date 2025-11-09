// productInfo.js

document.addEventListener("DOMContentLoaded", async () => {
  // 1️⃣ Get product_id from URL
  const urlParams = new URLSearchParams(window.location.search);
  const productId = urlParams.get("id");

  if (!productId) {
    console.error("Product ID not found in URL");
    return;
  }

  try {
    // 2️⃣ Fetch product info from your backend
    const response = await fetch(
      `http://localhost:3002/product-info?product_id=${productId}`
    );
    if (!response.ok) throw new Error("Failed to fetch product info");
    const data = await response.json();

    const product = data.product;
    const stores = data.stores;

    // 3️⃣ Update product section
    document.querySelector("h1").textContent = product.name;
    document.querySelector(".breadcrumb-item.active").textContent =
      product.name;
    document.querySelector(".ratio img").src = product.photo;
    document.querySelector(".ratio img").alt = product.name;
    document.querySelector("p.mb-3").textContent = product.description;
    document.querySelector(
      ".fs-5"
    ).textContent = `от ${product.price.toLocaleString("ru-RU")} ₸`;

    // 4️⃣ Fill stores list
    const shopsList = document.getElementById("shopsList");
    const shopsEmpty = document.getElementById("shopsEmpty");
    const storesCountTop = document.getElementById("storesCountTop");
    const storesFound = document.getElementById("storesFound");

    shopsList.innerHTML = "";
    storesCountTop.textContent = stores.length;
    storesFound.textContent = stores.length;

    if (stores.length === 0) {
      shopsEmpty.classList.remove("d-none");
      return;
    } else {
      shopsEmpty.classList.add("d-none");
    }

    stores.forEach((store) => {
      const priceText = store.price
        ? `${store.price.toLocaleString("ru-RU")} ₸`
        : "—";
      const availabilityText = store.availability ?? "Неизвестно";

      const col = document.createElement("div");
      col.className = "col shop-item";
      col.innerHTML = `
    <div class="card h-100 product-card">
      <div class="card-body d-flex flex-column">
        <div class="d-flex justify-content-between align-items-start">
          <div>
            <h6 class="mb-1 js-name">${store.name}</h6>
            <div class="text-muted small">
              <i class="bi bi-geo-alt"></i> ${
                store.address ?? "Адрес неизвестен"
              }
            </div>
          </div>
          <a href="/store.html?id=${
            store.store_id
          }" class="btn btn-outline-primary btn-sm">
            В магазин
          </a>
        </div>
        <div class="mt-3 d-flex gap-4">
          <div>
            <span class="fw-semibold">${priceText}</span>
            <span class="small text-muted d-block">цена</span>
          </div>
          <div>
            <span class="fw-semibold">${availabilityText}</span>
            <span class="small text-muted d-block">обновлено сегодня</span>
          </div>
        </div>
      </div>
    </div>
  `;
      shopsList.appendChild(col);
    });
  } catch (error) {
    console.error(error);
    document.getElementById("shopsList").innerHTML = `
      <div class="col-12 text-center text-danger">Ошибка загрузки продукта</div>
    `;
  }
});
