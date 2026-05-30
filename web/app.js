'use strict';

(function () {
  const API_FLY = 'https://pulsekeep.fly.dev';

  document.addEventListener('DOMContentLoaded', function () {
    var yearEl = document.getElementById('footer-year');
    if (yearEl) yearEl.textContent = new Date().getFullYear().toString();

    /* --- Mobile menu toggle --- */
    var toggle = document.getElementById('mobile-menu-toggle');
    var mobileNav = document.getElementById('mobile-nav');
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
    var statServers = document.getElementById('stat-servers');
    var statUsers = document.getElementById('stat-users');
    var statCommands = document.getElementById('stat-commands');
    var statUptime = document.getElementById('stat-uptime');
    var statusBadge = document.getElementById('bot-status-badge');

    function setStatus(state, label) {
      if (!statusBadge) return;
      var dot = statusBadge.querySelector('.pulse-dot');
      var txt = statusBadge.querySelector('.badge-text');
      statusBadge.className = 'hero-eyebrow';
      if (dot) dot.className = 'pulse-dot' + (state === 'online' ? '' : ' offline');
      if (txt) txt.textContent = label;
    }

    function fetchStats() {
      fetch(API_FLY + '/health')
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

    if (statusBadge || statServers) fetchStats();
  });
})();
