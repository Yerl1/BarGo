const $q = document.getElementById("q");
const $form = document.getElementById("searchForm");
const $clear = document.getElementById("clearBtn");
const $chip = document.getElementById("activeQueryChip");
const $sortChip = document.getElementById("activeSortChip");
const $clearAll = document.getElementById("clearAll");
const $found = document.getElementById("foundCount");
const $grid = document.getElementById("cardsGrid");
const $empty = document.getElementById("emptyState");
const $emptyQuery = document.getElementById("emptyQuery");
const $sortSelect = document.getElementById("sortSelect");
const $pagination = document.getElementById("pagination");

// элементы фильтров
const $category = document.getElementById("categoryFilter");
const $brand = document.getElementById("brandFilter");
const $priceMin = document.getElementById("priceMin");
const $priceMax = document.getElementById("priceMax");
const $applyFilters = document.getElementById("applyFilters");
const $activeFilters = document.getElementById("activeFilters");

const PER_PAGE = 9; // сколько карточек на страницу
let currentPage = 1;

// сохранённое состояние фильтров (то, что применено кнопкой "Применить фильтры")
let currentFilters = {
  category: $category ? $category.value : "",
  brand: $brand ? $brand.value : "",
  min: 0,
  max: Infinity,
};

// Карточки и исходный порядок
const allCards = Array.from(document.querySelectorAll(".prod-card")).map(
  (el, idx) => {
    el.dataset.index = idx;
    return el;
  }
);

const norm = (s) => (s || "").toString().trim().toLowerCase();

// --- Работа с фильтрами ---

function readFiltersFromUI() {
  const minVal =
    $priceMin && $priceMin.value !== "" ? Number($priceMin.value) : 0;
  const maxVal =
    $priceMax && $priceMax.value !== ""
      ? Number($priceMax.value)
      : Infinity;

  return {
    category: $category ? $category.value : "",
    brand: $brand ? $brand.value : "",
    min: Number.isFinite(minVal) ? minVal : 0,
    max: Number.isFinite(maxVal) ? maxVal : Infinity,
  };
}

function hasAnyFilters(f) {
  return (
    f.category ||
    f.brand ||
    (f.min && f.min > 0) ||
    (f.max && f.max !== Infinity)
  );
}

function renderFilterChips() {
  if (!$activeFilters) return;

  const f = currentFilters;
  $activeFilters.innerHTML = "";

  const chips = [];

  if (f.category && $category) {
    const label =
      $category.options[$category.selectedIndex].textContent || f.category;
    chips.push(`Категория: ${label}`);
  }

  if (f.brand && $brand) {
    const label =
      $brand.options[$brand.selectedIndex].textContent || f.brand;
    chips.push(`Бренд: ${label}`);
  }

  if (f.min && f.min > 0) {
    chips.push(`Цена от ${f.min}₸`);
  }

  if (f.max && f.max !== Infinity) {
    chips.push(`Цена до ${f.max}₸`);
  }

  chips.forEach((text) => {
    const span = document.createElement("span");
    span.className = "filter-chip";
    span.textContent = text;
    $activeFilters.appendChild(span);
  });
}

// подсветка совпадений
function highlightTitle(el, query) {
  const titleEl = el.querySelector(".js-title");
  const original = el.dataset.title;
  if (!query) {
    titleEl.innerHTML = original;
    return;
  }
  const q = norm(query).replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const re = new RegExp(`(${q})`, "ig");
  titleEl.innerHTML = original.replace(
    re,
    '<mark class="search-highlight">$1</mark>'
  );
}

// оценка релевантности
function relevanceScore(el, queryNorm) {
  if (!queryNorm) return Number(el.dataset.index);
  const title = norm(el.dataset.title);
  const kw = norm(el.dataset.keywords);
  const tPos = title.indexOf(queryNorm);
  const kPos = kw.indexOf(queryNorm);
  let pos = 1e9;
  if (tPos >= 0) pos = Math.min(pos, tPos);
  if (kPos >= 0) pos = Math.min(pos, kPos + 100); // хуже, чем в title
  return pos === 1e9
    ? 1e9 + Number(el.dataset.index)
    : pos + Number(el.dataset.index) / 1000;
}

// поиск + фильтры -> список видимых
function applySearch(query) {
  const qn = norm(query);
  const visible = [];
  const f = currentFilters;

  allCards.forEach((card) => {
    const title = norm(card.dataset.title);
    const keywords = norm(card.dataset.keywords);
    const category = card.dataset.category || "";
    const brand = card.dataset.brand || "";
    const price = Number(card.dataset.price || 0);

    let match =
      (!qn || title.includes(qn) || keywords.includes(qn)) &&
      (!f.category || category === f.category) &&
      (!f.brand || brand === f.brand);

    if (match) {
      if (f.min && price < f.min) match = false;
      if (f.max !== Infinity && price > f.max) match = false;
    }

    card.classList.toggle("d-none", !match); // фильтрация
    card.classList.remove("d-none-page"); // сброс скрытия пагинации
    highlightTitle(card, qn ? query : "");
    if (match) visible.push(card);
  });

  // чипы фильтров
  renderFilterChips();

  // UI по количеству
  $found.textContent = visible.length.toString();
  const showChip = !!qn;
  $chip.classList.toggle("d-none", !showChip);
  if (showChip) $chip.textContent = `Запрос: «${query}»`;

  const hasFiltersNow = hasAnyFilters(f);
  const showClearAll =
    showChip || $sortSelect.value !== "relevance" || hasFiltersNow;
  $clearAll.classList.toggle("d-none", !showClearAll);

  const nothing = visible.length === 0;
  $empty.classList.toggle("d-none", !nothing);
  if (nothing) $emptyQuery.textContent = query;

  return { qn, visible };
}

// сортировка видимых
function applySort(visible, qn) {
  const mode = $sortSelect.value;

  const sortLabel = {
    relevance: "По релевантности",
    price_asc: "Цена: сначала дешевле",
    price_desc: "Цена: сначала дороже",
    distance_asc: "Расстояние: ближайшие",
    distance_desc: "Расстояние: дальние",
  }[mode];

  const showSortChip = mode !== "relevance";
  $sortChip.classList.toggle("d-none", !showSortChip);
  if (showSortChip) $sortChip.textContent = `Сортировка: ${sortLabel}`;

  const cmp = {
    relevance: (a, b) => relevanceScore(a, qn) - relevanceScore(b, qn),
    price_asc: (a, b) =>
      Number(a.dataset.price) - Number(b.dataset.price) ||
      Number(a.dataset.index) - Number(b.dataset.index),
    price_desc: (a, b) =>
      Number(b.dataset.price) - Number(a.dataset.price) ||
      Number(a.dataset.index) - Number(b.dataset.index),
    distance_asc: (a, b) =>
      Number(a.dataset.distance) - Number(b.dataset.distance) ||
      Number(a.dataset.index) - Number(b.dataset.index),
    distance_desc: (a, b) =>
      Number(b.dataset.distance) - Number(a.dataset.distance) ||
      Number(a.dataset.index) - Number(b.dataset.index),
  }[mode];

  visible.sort(cmp);

  const frag = document.createDocumentFragment();
  visible.forEach((el) => frag.appendChild(el));
  allCards
    .filter((el) => el.classList.contains("d-none"))
    .forEach((el) => frag.appendChild(el));
  $grid.appendChild(frag);

  return visible;
}

// ---------- Пагинация ----------
function buildPagination(total, perPage) {
  const totalPages = Math.max(1, Math.ceil(total / perPage));
  currentPage = Math.min(currentPage, totalPages);

  $pagination.innerHTML = "";

  const makeItem = (label, page, disabled = false, active = false) => {
    const li = document.createElement("li");
    li.className = `page-item${disabled ? " disabled" : ""}${
      active ? " active" : ""
    }`;
    const a = document.createElement("a");
    a.className = "page-link";
    a.href = "#";
    a.textContent = label;
    a.addEventListener("click", (e) => {
      e.preventDefault();
      if (disabled || active) return;
      currentPage = page;
      applyAll($q.value, /*keepPage*/ true);
    });
    li.appendChild(a);
    return li;
  };

  // Prev
  $pagination.appendChild(
    makeItem("Предыдущая", currentPage - 1, currentPage === 1)
  );

  for (let p = 1; p <= totalPages; p++) {
    if (totalPages > 7) {
      if (
        p === 1 ||
        p === totalPages ||
        (p >= currentPage - 2 && p <= currentPage + 2)
      ) {
        $pagination.appendChild(
          makeItem(String(p), p, false, p === currentPage)
        );
      } else if (p === 2 && currentPage > 4) {
        $pagination.appendChild(makeItem("…", currentPage, true, false));
      } else if (p === totalPages - 1 && currentPage < totalPages - 3) {
        $pagination.appendChild(makeItem("…", currentPage, true, false));
      }
    } else {
      $pagination.appendChild(
        makeItem(String(p), p, false, p === currentPage)
      );
    }
  }

  // Next
  $pagination.appendChild(
    makeItem("Следующая", currentPage + 1, currentPage === totalPages)
  );
}

function applyPage(visible) {
  visible.forEach((el, idx) => {
    const onPage = Math.floor(idx / PER_PAGE) + 1 === currentPage;
    el.classList.toggle("d-none-page", !onPage);
  });
}

// Комбинированное применение
function applyAll(query, keepPage = false) {
  if (!keepPage) currentPage = 1;
  const { qn, visible } = applySearch(query);
  const sortedVisible = applySort(visible, qn);

  buildPagination(sortedVisible.length, PER_PAGE);
  applyPage(sortedVisible);
}
document.addEventListener("DOMContentLoaded", () => {
  const tabs = document.querySelectorAll(".view-tab");
  const grid = document.getElementById("cardsGrid");
  const mapView = document.getElementById("mapView");

  if (!tabs.length || !grid || !mapView) return;

  function setView(view) {
    if (view === "grid") {
      grid.style.display = "";
      mapView.style.display = "none";
    } else {
      grid.style.display = "none";
      mapView.style.display = "block";

      // Чиним размер карты, если она уже инициализирована
      if (window._leafletMap) {
        setTimeout(() => {
          window._leafletMap.invalidateSize();
        }, 200);
      }
    }
  }

  tabs.forEach((btn) => {
    btn.addEventListener("click", () => {
      tabs.forEach((b) => b.classList.remove("active"));
      btn.classList.add("active");

      const view = btn.dataset.view;
      setView(view);
    });
  });

  // По умолчанию – сетка
  setView("grid");
});


// ---------- LISTENERS ----------

let tId;
$q.addEventListener("input", () => {
  clearTimeout(tId);
  tId = setTimeout(() => applyAll($q.value), 140);
});

$form.addEventListener("submit", (e) => {
  e.preventDefault();
  applyAll($q.value);
});

// "Применить фильтры"
if ($applyFilters) {
  $applyFilters.addEventListener("click", () => {
    currentFilters = readFiltersFromUI();
    applyAll($q.value);
  });
}

// Сброс в блоке фильтров (кнопка "Сбросить")
$clear.addEventListener("click", () => {
  $q.value = "";
  $sortSelect.value = "relevance";

  if ($category) $category.value = "";
  if ($brand) $brand.value = "";
  if ($priceMin) $priceMin.value = "";
  if ($priceMax) $priceMax.value = "";

  currentFilters = {
    category: "",
    brand: "",
    min: 0,
    max: Infinity,
  };

  applyAll("");
  $q.focus();
});

// Чип "Очистить всё" рядом с поиском
$clearAll.addEventListener("click", () => {
  $q.value = "";
  $sortSelect.value = "relevance";

  if ($category) $category.value = "";
  if ($brand) $brand.value = "";
  if ($priceMin) $priceMin.value = "";
  if ($priceMax) $priceMax.value = "";

  currentFilters = {
    category: "",
    brand: "",
    min: 0,
    max: Infinity,
  };

  applyAll("");
  $q.focus();
});

$sortSelect.addEventListener("change", () => {
  applyAll($q.value);
});

// init
applyAll("");
