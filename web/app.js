'use strict';

(function () {
  const API = window.location.origin;

  function getEl(id) { return document.getElementById(id); }
  function setText(id, val) { const e = getEl(id); if (e) e.textContent = val; }

  document.addEventListener('DOMContentLoaded', function () {
    const yearEl = getEl('footer-year');
    if (yearEl) yearEl.textContent = new Date().getFullYear().toString();

    /* --- Mobile menu toggle --- */
    var toggle = getEl('mobile-menu-toggle');
    var mobileNav = getEl('mobile-nav');
    if (toggle && mobileNav) {
      toggle.addEventListener('click', function () {
        var expanded = toggle.getAttribute('aria-expanded') === 'true' ? false : true;
        toggle.setAttribute('aria-expanded', expanded);
        mobileNav.classList.toggle('open');
      });
    }

    /* --- Scroll reveal animation --- */
    var revealEls = document.querySelectorAll('.scroll-reveal');
    if (revealEls.length) {
      var observer = new IntersectionObserver(function (entries) {
        entries.forEach(function (entry) {
          if (entry.isIntersecting) {
            entry.target.classList.add('visible');
            observer.unobserve(entry.target);
          }
        });
      }, { threshold: 0.1 });
      revealEls.forEach(function (el) { observer.observe(el); });
    }

    /* --- Home page: live stats + status badge --- */
    var statServers = getEl('stat-servers');
    var statUsers = getEl('stat-users');
    var statCommands = getEl('stat-commands');
    var statUptime = getEl('stat-uptime');
    var statusBadge = getEl('bot-status-badge');

    function fetchStats() {
      fetch(API + '/health')
        .then(function (r) { return r.ok ? r.json() : null; })
        .then(function (d) {
          if (!d) { setStatus('offline', 'Service offline'); return; }
          if (statServers) statServers.textContent = (d.servers || 0).toLocaleString();
          if (statUsers) statUsers.textContent = (d.users || 0).toLocaleString();
          if (statCommands) statCommands.textContent = (d.commands || 0).toLocaleString();
          if (statUptime) statUptime.textContent = d.uptime || '--';
          setStatus(d.database === 'ok' ? 'online' : 'degraded', d.database === 'ok' ? 'Service online' : 'Degraded');
        })
        .catch(function () { setStatus('offline', 'Offline'); });
    }

    function setStatus(state, label) {
      if (!statusBadge) return;
      statusBadge.className = 'hero-eyebrow';
      var dot = statusBadge.querySelector('.pulse-dot');
      var txt = statusBadge.querySelector('.badge-text');
      if (dot) dot.className = 'pulse-dot' + (state === 'online' ? '' : ' offline');
      if (txt) txt.textContent = label;
    }

    /* Only run on pages that have these elements */
    if (statusBadge || statServers) fetchStats();
  });
})();
