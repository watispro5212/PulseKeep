'use strict';
const API = window.API_BASE || window.location.origin;

function $(id) { return document.getElementById(id) }
function qs(s, p) { return (p||document).querySelector(s) }
function qsa(s, p) { return (p||document).querySelectorAll(s) }

function toast(msg, type) {
  const el = $('toast');
  if (!el) return;
  el.className = 'toast show ' + type;
  el.innerHTML = '<i class="fa-solid fa-' + (type==='success'?'check-circle':'exclamation-circle') + '"></i>' + msg;
  clearTimeout(el._t);
  el._t = setTimeout(() => el.classList.remove('show'), 3000);
}

document.addEventListener('DOMContentLoaded', function () {
  const fy = $('footer-year');
  if (fy) fy.textContent = String(new Date().getFullYear());

  // Active nav
  const path = window.location.pathname.split('/').pop() || 'index.html';
  document.querySelectorAll('header nav a, .mobile-nav a').forEach(a => {
    if (a.getAttribute('href') === path) a.classList.add('active');
  });

  // Mobile menu
  const toggle = $('mobile-toggle');
  const mnav = $('mobile-nav');
  if (toggle && mnav) {
    toggle.addEventListener('click', () => {
      mnav.classList.toggle('open');
      toggle.innerHTML = mnav.classList.contains('open')
        ? '<i class="fa-solid fa-xmark"></i>'
        : '<i class="fa-solid fa-bars"></i>';
    });
    document.addEventListener('click', (e) => {
      if (!mnav.contains(e.target) && !toggle.contains(e.target)) {
        mnav.classList.remove('open');
        toggle.innerHTML = '<i class="fa-solid fa-bars"></i>';
      }
    });
  }

  // Command search
  const cmdSearch = $('cmd-search');
  if (cmdSearch) {
    cmdSearch.addEventListener('input', function () {
      const q = this.value.toLowerCase().trim();
      qsa('.cmd-group').forEach(g => {
        let visible = 0;
        qsa('span[data-name]', g).forEach(s => {
          const match = !q || s.dataset.name.includes(q);
          s.style.display = match ? 'inline-block' : 'none';
          if (match) visible++;
        });
        g.style.display = visible === 0 ? 'none' : '';
        const cnt = qs('.cmd-count', g);
        if (cnt) cnt.textContent = '(' + visible + '/' + qsa('span[data-name]', g).length + ')';
      });
    });
  }

  // Stats fetching
  async function fetchStats() {
    try {
      const res = await fetch(API + '/api/stats');
      if (!res.ok) throw new Error('API down');
      const d = await res.json();

      ['servers','users','commands','uptime'].forEach(key => {
        const el = $('stat-' + key);
        if (!el) return;
        const val = key === 'servers' ? d.guilds : key === 'users' ? d.users : key === 'commands' ? d.commandsRun : formatUptime(d.uptime);
        if (key !== 'uptime' && typeof val === 'number') {
          animateCounter(el, val);
        } else {
          el.textContent = val ?? '--';
        }
      });

      const dot = $('status-dot');
      const txt = $('status-text');
      if (dot) dot.className = 'dot dot-green';
      if (txt) txt.textContent = 'Online \u00b7 ' + d.guilds + ' servers';

      ['st-bot','st-api'].forEach(id => {
        const el = $(id);
        if (el) el.className = 'dot dot-green';
        const tel = $(id + '-text');
        if (tel) tel.textContent = id === 'st-bot' ? 'Online \u00b7 ' + d.guilds + ' servers' : 'Online';
      });
      const dbEl = $('st-db');
      const dbTxt = $('st-db-text');
      if (dbEl) dbEl.className = d.dbConnected ? 'dot dot-green' : 'dot dot-rose';
      if (dbTxt) dbTxt.textContent = d.dbConnected ? 'Connected' : 'Disconnected';
      const latEl = $('st-latency');
      const latTxt = $('st-latency-text');
      if (latEl) latEl.className = d.avgLatency < 200 ? 'dot dot-green' : 'dot dot-amber';
      if (latTxt) latTxt.textContent = (d.avgLatency || '--') + 'ms';
      const uptEl = $('st-uptime-text');
      if (uptEl) uptEl.textContent = formatUptime(d.uptime);
      const shEl = $('st-shards');
      const shTxt = $('st-shards-text');
      if (shEl) shEl.className = (d.shards || 0) > 0 ? 'dot dot-green' : 'dot dot-amber';
      if (shTxt) shTxt.textContent = (d.shards || 0) > 0 ? d.shards + ' online' : 'Starting...';
      const upEl = $('st-updated');
      if (upEl) upEl.textContent = new Date().toLocaleTimeString();

    } catch {
      const dot = $('status-dot');
      const txt = $('status-text');
      if (dot) dot.className = 'dot dot-rose';
      if (txt) txt.textContent = 'Offline';

      ['st-bot','st-api','st-db','st-shards'].forEach(id => {
        const el = $(id);
        if (el) el.className = 'dot dot-rose';
        const tel = $(id + '-text');
        if (tel) tel.textContent = id === 'st-db' ? 'Unknown' : 'Offline';
      });
      const upEl = $('st-updated');
      if (upEl) upEl.textContent = '--';
    }
  }

  function animateCounter(el, target) {
    const current = parseInt(el.textContent.replace(/,/g,''), 10);
    if (isNaN(current) || current === target) { el.textContent = target.toLocaleString(); return; }
    const duration = 1200;
    const start = performance.now();
    const from = current;
    function tick(now) {
      const p = Math.min((now - start) / duration, 1);
      const ease = 1 - Math.pow(1 - p, 3);
      el.textContent = Math.round(from + (target - from) * ease).toLocaleString();
      if (p < 1) requestAnimationFrame(tick);
    }
    requestAnimationFrame(tick);
  }

  function formatUptime(s) {
    if (!s && s !== 0) return '--';
    if (s < 0) return '--';
    const d = Math.floor(s / 86400);
    const h = Math.floor((s % 86400) / 3600);
    const m = Math.floor((s % 3600) / 60);
    let r = '';
    if (d > 0) r += d + 'd ';
    if (h > 0 || d > 0) r += h + 'h ';
    r += m + 'm';
    return r || '0m';
  }

  if ($('stat-servers') || $('status-dot') || $('st-bot')) {
    fetchStats();
    setInterval(fetchStats, 30000);
  }

  // Scroll reveal
  if ('IntersectionObserver' in window) {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(e => { if (e.isIntersecting) { e.target.classList.add('visible'); observer.unobserve(e.target) } });
    }, { threshold: .1 });
    document.querySelectorAll('.reveal').forEach(el => observer.observe(el));
  } else {
    document.querySelectorAll('.reveal').forEach(el => el.classList.add('visible'));
  }
});
