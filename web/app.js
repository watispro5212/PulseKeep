'use strict';
const API = window.API_BASE || window.location.origin;

document.addEventListener('DOMContentLoaded', function () {

  // Footer year
  const fy = document.getElementById('footer-year');
  if (fy) fy.textContent = new Date().getFullYear().toString();

  // Mobile menu
  const toggle = document.getElementById('mobile-toggle');
  const mnav = document.getElementById('mobile-nav');
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

  // Fetch stats from API
  async function fetchStats() {
    try {
      const res = await fetch(API + '/api/stats');
      if (!res.ok) throw new Error('API down');
      const d = await res.json();

      const ids = ['stat-servers','stat-users','stat-commands','stat-uptime'];
      const vals = [d.guilds, d.users, d.commandsRun, formatUptime(d.uptime)];
      ids.forEach((id, i) => {
        const el = document.getElementById(id);
        if (el) el.textContent = vals[i] ?? '--';
      });

      // Status badge on homepage
      const dot = document.getElementById('status-dot');
      const txt = document.getElementById('status-text');
      if (dot && txt) {
        dot.className = 'dot dot-green';
        txt.textContent = 'Online · ' + d.guilds + ' servers';
      }

      // Status page
      const els = ['st-bot','st-api','st-db'];
      const sts = ['dot-green','dot-green','dot-green'];
      const lst = [`Online · ${d.guilds} servers`, 'Online', 'Connected'];
      els.forEach((id, i) => {
        const el = document.getElementById(id);
        if (el) el.className = 'dot ' + sts[i];
        const tel = document.getElementById(id + '-text');
        if (tel) tel.textContent = lst[i];
      });

      const latEl = document.getElementById('st-latency');
      const latTxt = document.getElementById('st-latency-text');
      if (latEl) latEl.className = d.avgLatency < 200 ? 'dot dot-green' : 'dot dot-amber';
      if (latTxt) latTxt.textContent = d.avgLatency + 'ms';

    } catch {
      const dot = document.getElementById('status-dot');
      const txt = document.getElementById('status-text');
      if (dot) dot.className = 'dot dot-rose';
      if (txt) txt.textContent = 'Offline';

      ['st-bot','st-api','st-db'].forEach(id => {
        const el = document.getElementById(id);
        if (el) el.className = 'dot dot-rose';
        const tel = document.getElementById(id + '-text');
        if (tel) tel.textContent = 'Offline';
      });
    }
  }

  function formatUptime(s) {
    if (!s) return '--';
    const d = Math.floor(s / 86400);
    const h = Math.floor((s % 86400) / 3600);
    const m = Math.floor((s % 3600) / 60);
    let r = '';
    if (d > 0) r += d + 'd ';
    if (h > 0 || d > 0) r += h + 'h ';
    r += m + 'm';
    return r;
  }

  if (document.getElementById('stat-servers') || document.getElementById('status-dot')) {
    fetchStats();
    setInterval(fetchStats, 30000);
  }

  // --- Dashboard ---
  const guildList = document.getElementById('guild-list');
  const dashConfig = document.getElementById('dash-config');
  const dashPrompt = document.getElementById('dash-select-prompt');

  if (guildList) {
    // Check for saved token
    let accessToken = localStorage.getItem('pk_token');

    function renderGuildList(guilds) {
      guildList.innerHTML = '';
      if (guilds.length === 0) {
        guildList.innerHTML = '<div style="color:var(--muted);font-size:.82rem;padding:8px">No mutual servers found.</div>';
        return;
      }
      guilds.forEach(g => {
        const div = document.createElement('div');
        div.className = 'dash-guild-item';
        const icon = g.icon
          ? `https://cdn.discordapp.com/icons/${g.id}/${g.icon}.png`
          : 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 28 28"%3E%3Crect fill="%23333" width="28" height="28" rx="14"/%3E%3Ctext x="14" y="18" text-anchor="middle" fill="%23fff" font-size="14"%3E${g.name.charAt(0).toUpperCase()}%3C/text%3E%3C/svg%3E';
        div.innerHTML = `<img src="${icon}" alt="">${g.name}`;
        div.addEventListener('click', () => {
          document.querySelectorAll('.dash-guild-item').forEach(el => el.classList.remove('active'));
          div.classList.add('active');
          loadConfig(g.id, g.name);
        });
        guildList.appendChild(div);
      });
    }

    async function loadConfig(guildId, guildName) {
      document.getElementById('dash-guild-name').textContent = guildName;
      document.getElementById('dash-guild-id').textContent = 'ID: ' + guildId;
      dashPrompt.style.display = 'none';
      dashConfig.style.display = 'block';
      document.getElementById('dash-saved').style.display = 'none';

      try {
        const res = await fetch(API + '/api/guild/' + guildId + '/config');
        if (!res.ok) throw new Error('No config');
        const cfg = await res.json();

        ['economyEnabled','ticketsEnabled','modlogsEnabled','welcomeEnabled'].forEach(key => {
          const el = document.querySelector(`[data-key="${key}"]`);
          if (el) el.checked = cfg[key] !== false;
        });

        const tc = document.querySelector('[data-key="ticketCategoryId"]');
        if (tc) tc.value = cfg.ticketCategoryId || '';

        const lc = document.querySelector('[data-key="logChannelId"]');
        if (lc) lc.value = cfg.logChannelId || '';

      } catch {
        // Defaults
        document.querySelectorAll('[data-key]').forEach(el => {
          if (el.type === 'checkbox') el.checked = true;
          else el.value = '';
        });
      }
    }

    document.getElementById('dash-save')?.addEventListener('click', async () => {
      const guildId = document.getElementById('dash-guild-id').textContent.replace('ID: ','');
      if (!guildId) return;

      const payload = {};
      document.querySelectorAll('[data-key]').forEach(el => {
        const key = el.getAttribute('data-key');
        payload[key] = el.type === 'checkbox' ? el.checked : el.value;
      });

      try {
        const res = await fetch(API + '/api/guild/' + guildId + '/config', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
        });
        if (res.ok) {
          const saved = document.getElementById('dash-saved');
          saved.style.display = 'block';
          setTimeout(() => { saved.style.display = 'none'; }, 3000);
        }
      } catch (err) {
        alert('Failed to save config.');
      }
    });

    // Try OAuth login
    async function tryLogin() {
      guildList.innerHTML = '<div style="color:var(--muted);font-size:.82rem;padding:8px">Redirecting to Discord...</div>';
      // For demo purposes, show mock guilds
      renderGuildList([
        { id: '123456789', name: 'Demo Server', icon: null },
        { id: '987654321', name: 'Test Community', icon: null },
      ]);
    }

    tryLogin();
  }

});
