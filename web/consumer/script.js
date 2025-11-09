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

const PER_PAGE = 9; // сколько карточек на страницу
let currentPage = 1;

// Карточки и исходный порядок
const allCards = Array.from(document.querySelectorAll(".prod-card")).map(
  (el, idx) => {
    el.dataset.index = idx;
    return el;
  }
);

const norm = (s) => (s || "").toString().trim().toLowerCase();

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

// поиск -> список видимых
function applySearch(query) {
  const qn = norm(query);
  const visible = [];

  allCards.forEach((card) => {
    const title = norm(card.dataset.title);
    const keywords = norm(card.dataset.keywords);
    const match = !qn || title.includes(qn) || keywords.includes(qn);

    card.classList.toggle("d-none", !match); // фильтрация
    card.classList.remove("d-none-page"); // сброс скрытия пагинации
    highlightTitle(card, qn ? query : "");
    if (match) visible.push(card);
  });

  // UI
  $found.textContent = visible.length.toString();
  const showChip = !!qn;
  $chip.classList.toggle("d-none", !showChip);
  if (showChip) $chip.textContent = `Запрос: «${query}»`;

  const showClearAll = showChip || $sortSelect.value !== "relevance";
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

  // Pages (простая версия: показываем до 7 страниц; при большем — можно добавить «...»)
  const totalPagesToShow = Math.min(7, Math.ceil(total / perPage));
  for (let p = 1; p <= totalPages; p++) {
    // если страниц много — показываем «окно» вокруг текущей
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
      $pagination.appendChild(makeItem(String(p), p, false, p === currentPage));
    }
  }

  // Next
  $pagination.appendChild(
    makeItem("Следующая", currentPage + 1, currentPage === totalPages)
  );
}

function applyPage(visible) {
  // Скрываем элементы, которые не попали на текущую страницу
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

// listeners
let tId;
$q.addEventListener("input", () => {
  clearTimeout(tId);
  tId = setTimeout(() => applyAll($q.value), 140);
});

$form.addEventListener("submit", (e) => {
  e.preventDefault();
  applyAll($q.value);
});

$clear.addEventListener("click", () => {
  $q.value = "";
  $sortSelect.value = "relevance";
  applyAll("");
  $q.focus();
});

$clearAll.addEventListener("click", () => {
  $q.value = "";
  $sortSelect.value = "relevance";
  applyAll("");
  $q.focus();
});

$sortSelect.addEventListener("change", () => {
  applyAll($q.value);
});

// init
applyAll("");
