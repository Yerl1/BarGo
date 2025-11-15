// storeScript.js
document.addEventListener("DOMContentLoaded", () => {
  const $q = document.getElementById("q");
  const $sort = document.getElementById("sortSelect");
  const $grid = document.getElementById("storeGrid");
  const $found = document.getElementById("foundCount");
  const $empty = document.getElementById("emptyState");
  const $emptyQuery = document.getElementById("emptyQuery");
  const $pagination = document.getElementById("pagination");

  if (!$grid) return;

  const PER_PAGE = 9;
  let currentPage = 1;
  let lastQuery = $q ? $q.value : "";

  let allCards = [];

  function recomputeCards() {
    allCards = Array.from($grid.querySelectorAll(".store-prod"));
    allCards.forEach((card, idx) => {
      card.dataset.index = String(idx);
    });
  }

  recomputeCards();

  const norm = (v) =>
    (v || "").toString().trim().toLowerCase();

  function relevanceScore(card, qn) {
    if (!qn) return Number(card.dataset.index);
    const title = norm(card.dataset.title);
    const kw = norm(card.dataset.keywords);
    const tPos = title.indexOf(qn);
    const kPos = kw.indexOf(qn);
    let pos = 1e9;
    if (tPos >= 0) pos = Math.min(pos, tPos);
    if (kPos >= 0) pos = Math.min(pos, kPos + 100);
    return pos === 1e9
      ? 1e9 + Number(card.dataset.index)
      : pos + Number(card.dataset.index) / 1000;
  }
  async function loadStoreFromApi() {
    const params = new URLSearchParams(window.location.search);
    const storeId = params.get("id");
    if (!storeId) return;

    try {
      const resp = await fetch(
        `http://localhost:3002/store-info?store_id=${encodeURIComponent(storeId)}`
      );
      if (!resp.ok) {
        throw new Error(`Failed to load store info: ${resp.status}`);
      }

      const data = await resp.json();
      const store = data.store;
      const products = data.products || [];

      // --- шапка магазина ---
      const titleEl = document.getElementById("storeTitle");
      if (titleEl && store.name) {
        titleEl.textContent = store.name;
      }

      const metaEl = document.getElementById("storeMeta");
      if (metaEl && (store.address || store.distance_km != null)) {
        const distancePart =
          store.distance_km != null ? ` • <b>${store.distance_km}</b> км` : "";
        metaEl.innerHTML = `
          <i class="bi bi-geo-alt"></i>
          ${store.address || ""}
          ${distancePart}
          <span class="ms-3 badge bg-chip">Открыто сейчас</span>
        `;
      }

      // --- товары магазина ---
      if (products.length) {
        $grid.innerHTML = "";
        products.forEach((p) => {
          const col = document.createElement("div");
          col.className = "col store-prod";
          col.dataset.title = p.name || "";
          col.dataset.keywords = (p.keywords || p.name || "").toLowerCase();
          col.dataset.price = p.price ?? 0;

          col.innerHTML = `
            <div class="card h-100 product-card">
              <div class="ratio ratio-16x9">
                <img src="${p.photo || "https://via.placeholder.com/640x360"}"
                     class="card-img-top img-cover"
                     alt="${p.name || ""}">
              </div>
              <div class="card-body d-flex flex-column">
                <h6 class="card-title js-title mb-1">${p.name || "Товар"}</h6>
                <div class="text-muted small mb-2">В наличии</div>
                <div class="mt-auto d-flex align-items-center justify-content-between">
                  <div class="fw-semibold text-primary-emphasis">${p.price ?? 0} ₸</div>
                  ${p.id
              ? `<a href="product.html?id=${encodeURIComponent(
                p.id
              )}" class="btn btn-outline-primary btn-sm">К описанию</a>`
              : `<button class="btn btn-outline-primary btn-sm" disabled>К описанию</button>`
            }
                </div>
              </div>
            </div>
          `;
          $grid.appendChild(col);
        });

        // После того, как перерисовали сетку, пересчитываем карточки и обновляем фильтры/пагинацию
        recomputeCards();
        update(false);
      }
    } catch (err) {
      console.error("Failed to load store info:", err);
    }
  }

  function filterCards(query) {
    const qn = norm(query);
    const visible = [];

    allCards.forEach((card) => {
      const title = norm(card.dataset.title);
      const kw = norm(card.dataset.keywords);
      const match =
        !qn || title.includes(qn) || kw.includes(qn);

      card.classList.toggle("d-none", !match);
      card.classList.remove("d-none-page");

      if (match) visible.push(card);
    });

    if ($found) $found.textContent = String(visible.length);

    const nothing = visible.length === 0;
    if ($empty) $empty.classList.toggle("d-none", !nothing);
    if ($emptyQuery) $emptyQuery.textContent = query || "";

    return { qn, visible };
  }

  function sortCards(visible, qn) {
    const sort = $sort ? $sort.value : "relevance";
    const cmpMap = {
      relevance: (a, b) =>
        relevanceScore(a, qn) - relevanceScore(b, qn),
      price_asc: (a, b) =>
        Number(a.dataset.price || 0) -
        Number(b.dataset.price || 0) ||
        Number(a.dataset.index) - Number(b.dataset.index),
      price_desc: (a, b) =>
        Number(b.dataset.price || 0) -
        Number(a.dataset.price || 0) ||
        Number(a.dataset.index) - Number(b.dataset.index),
    };
    const cmp = cmpMap[sort] || cmpMap.relevance;
    return [...visible].sort(cmp);
  }

  function buildPagination(total) {
    if (!$pagination) return;

    const pages = Math.max(1, Math.ceil(total / PER_PAGE));
    if (pages <= 1) {
      $pagination.classList.add("d-none");
      $pagination.innerHTML = "";
      return;
    }
    $pagination.classList.remove("d-none");
    if (currentPage > pages) currentPage = pages;

    $pagination.innerHTML = `
      <button class="page-btn" data-page="prev" ${currentPage === 1 ? "disabled" : ""
      }>Предыдущая</button>
      <span class="page-current">${currentPage}</span>
      <button class="page-btn" data-page="next" ${currentPage === pages ? "disabled" : ""
      }>Следующая</button>
    `;

    $pagination
      .querySelectorAll(".page-btn")
      .forEach((btn) => {
        btn.addEventListener("click", () => {
          const type = btn.dataset.page;
          if (type === "prev" && currentPage > 1) {
            currentPage--;
          } else if (
            type === "next" &&
            currentPage < pages
          ) {
            currentPage++;
          }
          update(true);
        });
      });
  }

  function applyPage(sortedVisible) {
    const start = (currentPage - 1) * PER_PAGE;
    const end = start + PER_PAGE;
    const pageItems = new Set(
      sortedVisible.slice(start, end)
    );

    const frag = document.createDocumentFragment();
    allCards.forEach((card) => {
      if (card.classList.contains("d-none")) {
        // уже отфильтрован — всё равно скрыт
        frag.appendChild(card);
        return;
      }
      if (pageItems.has(card)) {
        card.classList.remove("d-none-page");
      } else {
        card.classList.add("d-none-page");
      }
      frag.appendChild(card);
    });
    $grid.appendChild(frag);
  }

  function update(keepPage = false) {
    if (!keepPage) currentPage = 1;
    const { qn, visible } = filterCards(lastQuery);
    const sortedVisible = sortCards(visible, qn);
    buildPagination(sortedVisible.length);
    applyPage(sortedVisible);
  }

  // initial
  update();

  if ($q) {
    let tId;
    $q.addEventListener("input", () => {
      clearTimeout(tId);
      tId = setTimeout(() => {
        lastQuery = $q.value;
        update();
      }, 120);
    });
  }

  if ($sort) {
    $sort.addEventListener("change", () => {
      update(true); // страницу сохраняем
    });
  }

    // initial
    update();

    // ... слушатели $q, $sort и т.д.
  
    window.storeUpdate = () => update(false);
  
    // Загрузить магазин и товары по id из URL
    loadStoreFromApi();
});

document.addEventListener("DOMContentLoaded", () => {
  const params = new URLSearchParams(window.location.search);
  const storeId = params.get("id");
  if (!storeId) return;

  fetch(`http://localhost:3002/stores/${encodeURIComponent(storeId)}`)
    .then((resp) => {
      if (!resp.ok) {
        throw new Error(`Failed to load store: ${resp.status}`);
      }
      return resp.json();
    })
    .then((store) => {
      // Название
      const titleEl = document.getElementById("storeTitle");
      if (titleEl && store.name) {
        titleEl.textContent = store.name;
      }

      // Адрес + бейдж "Открыто"
      const metaEl = document.getElementById("storeMeta");
      if (metaEl && store.address) {
        metaEl.innerHTML = `
          <i class="bi bi-geo-alt"></i> ${store.address}
          <span class="ms-3 badge bg-chip">Открыто сейчас</span>
        `;
      }

      // Описание
      const descEl = document.getElementById("storeDescription");
      if (descEl && store.description) {
        descEl.textContent = store.description;
      }

      // Рейтинг сверху справа (используется твоим скриптом отзывов)
      const ratingValEl = document.getElementById("infoRatingValue");
      const ratingCntEl = document.getElementById("infoRatingCount");

      if (ratingValEl && typeof store.average_rating === "number") {
        ratingValEl.textContent = `${store.average_rating.toFixed(1)}/5`;
      }
      if (ratingCntEl && typeof store.review_count === "number") {
        ratingCntEl.textContent = store.review_count;
      }

      // Если захочешь — сюда можно добавить установку фоновой фотки магазина
      // например, поменять картинку в блоке карты или сделать баннер.
    })
    .catch((err) => {
      console.error("Error loading store:", err);
    });
});
