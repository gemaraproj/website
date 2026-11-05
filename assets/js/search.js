// SPDX-License-Identifier: Apache-2.0

(function () {
  'use strict';

  const form = document.querySelector('.sidebar-search');
  if (!form) return;
  const input = form.querySelector('.sidebar-search-input');
  const results = form.querySelector('.sidebar-search-results');
  if (!input || !results) return;

  const indexUrl = input.dataset.index || '/search.json';
  let index = null;
  let pending = null;
  let selectedIndex = -1;

  function loadIndex() {
    if (index) return Promise.resolve();
    if (pending) return pending;
    pending = fetch(indexUrl)
      .then(function (res) {
        if (!res.ok) throw new Error('HTTP ' + res.status);
        return res.json();
      })
      .then(function (data) {
        index = data;
        pending = null;
      })
      .catch(function (err) {
        pending = null;
        console.error('Failed to load search index:', err);
      });
    return pending;
  }

  function escapeHTML(s) {
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }

  function search(query) {
    if (!index || !query) return [];
    const q = query.toLowerCase();
    const matches = [];
    for (let i = 0; i < index.length; i++) {
      const entry = index[i];
      const title = (entry.title || '').toLowerCase();
      const content = (entry.content || '').toLowerCase();
      const titleMatch = title.indexOf(q) !== -1;
      const contentMatch = content.indexOf(q) !== -1;
      if (!titleMatch && !contentMatch) continue;
      matches.push({
        entry: entry,
        score: (titleMatch ? 10 : 0) + (contentMatch ? 1 : 0),
      });
    }
    matches.sort(function (a, b) { return b.score - a.score; });
    return matches.slice(0, 10).map(function (m) { return m.entry; });
  }

  function updateSelected() {
    const items = results.querySelectorAll('li:not(.search-empty)');
    for (let i = 0; i < items.length; i++) {
      if (i === selectedIndex) {
        items[i].classList.add('selected');
        items[i].scrollIntoView({ block: 'nearest' });
      } else {
        items[i].classList.remove('selected');
      }
    }
  }

  function render(matches, query) {
    selectedIndex = -1;
    if (!query) {
      results.innerHTML = '';
      results.hidden = true;
      return;
    }
    if (matches.length === 0) {
      results.innerHTML = '<li class="search-empty">No matches for &ldquo;' + escapeHTML(query) + '&rdquo;.</li>';
      results.hidden = false;
      return;
    }
    results.innerHTML = matches
      .map(function (m) {
        return '<li><a href="' + escapeHTML(m.url) + '">' + escapeHTML(m.title) + '</a></li>';
      })
      .join('');
    results.hidden = false;
  }

  input.addEventListener('focus', function () {
    loadIndex().then(function () {
      const q = input.value.trim();
      if (q) render(search(q), q);
    });
  });

  input.addEventListener('input', function () {
    loadIndex().then(function () {
      const q = input.value.trim();
      render(search(q), q);
    });
  });

  // Hide results when focus leaves the search area, but allow link clicks first.
  input.addEventListener('blur', function () {
    setTimeout(function () {
      if (!results.contains(document.activeElement)) {
        results.hidden = true;
        selectedIndex = -1;
      }
    }, 150);
  });

  input.addEventListener('keydown', function (e) {
    if (results.hidden) {
      if (e.key === 'Escape') {
        input.value = '';
        input.blur();
      }
      return;
    }

    const items = results.querySelectorAll('li:not(.search-empty)');
    const count = items.length;

    if (e.key === 'ArrowDown') {
      if (count === 0) return;
      e.preventDefault();
      selectedIndex = (selectedIndex + 1) % count;
      updateSelected();
    } else if (e.key === 'ArrowUp') {
      if (count === 0) return;
      e.preventDefault();
      if (selectedIndex === -1) {
        selectedIndex = count - 1;
      } else {
        selectedIndex = (selectedIndex - 1 + count) % count;
      }
      updateSelected();
    } else if (e.key === 'Enter') {
      if (selectedIndex >= 0 && selectedIndex < count) {
        e.preventDefault();
        const link = items[selectedIndex].querySelector('a');
        if (link) {
          link.click();
        }
      }
    } else if (e.key === 'Escape') {
      input.value = '';
      results.hidden = true;
      selectedIndex = -1;
      input.blur();
    }
  });
})();
