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
  
    const allCards = Array.from(
      $grid.querySelectorAll(".store-prod")
    );
    allCards.forEach((card, idx) => {
      card.dataset.index = String(idx);
    });
  
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
        <button class="page-btn" data-page="prev" ${
          currentPage === 1 ? "disabled" : ""
        }>Предыдущая</button>
        <span class="page-current">${currentPage}</span>
        <button class="page-btn" data-page="next" ${
          currentPage === pages ? "disabled" : ""
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
  
    // на всякий случай
    window.storeUpdate = () => update(false);
  });
  