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
        el.textContent = val ?? '--';
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
      const upEl = $('st-updated');
      if (upEl) upEl.textContent = new Date().toLocaleTimeString();

    } catch {
      const dot = $('status-dot');
      const txt = $('status-text');
      if (dot) dot.className = 'dot dot-rose';
      if (txt) txt.textContent = 'Offline';

      ['st-bot','st-api','st-db'].forEach(id => {
        const el = $(id);
        if (el) el.className = 'dot dot-rose';
        const tel = $(id + '-text');
        if (tel) tel.textContent = id === 'st-db' ? 'Unknown' : 'Offline';
      });
      const upEl = $('st-updated');
      if (upEl) upEl.textContent = '--';
    }
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

  // --- Dashboard ---
  const guildList = $('guild-list');
  const dashConfig = $('dash-config');
  const dashPrompt = $('dash-select-prompt');
  const loginBtn = $('dash-login-btn');
  const logoutBtn = $('dash-logout-btn');
  const sidebar = $('dash-sidebar');
  const sidebarToggle = $('dash-sidebar-toggle');

  if (guildList) {
    let accessToken = localStorage.getItem('pk_token');

    const sidebarClose = $('dash-sidebar-close');

    function closeSidebar() {
      if (sidebar) sidebar.classList.remove('open');
    }

    if (sidebarToggle) {
      sidebarToggle.addEventListener('click', () => {
        sidebar.classList.toggle('open');
      });
      document.addEventListener('click', (e) => {
        if (sidebar && sidebar.classList.contains('open') && !sidebar.contains(e.target) && !sidebarToggle.contains(e.target)) {
          closeSidebar();
        }
      });
    }
    if (sidebarClose) {
      sidebarClose.addEventListener('click', closeSidebar);
    }

    function checkAuth() {
      const params = new URLSearchParams(window.location.search);
      const token = params.get('token');
      if (token) {
        accessToken = token;
        localStorage.setItem('pk_token', token);
        window.history.replaceState({}, '', window.location.pathname);
        updateAuthUI();
        loadGuilds();
        return true;
      }
      if (accessToken) {
        updateAuthUI();
        loadGuilds();
        return true;
      }
      updateAuthUI();
      return false;
    }

    function updateAuthUI() {
      if (loginBtn) loginBtn.style.display = accessToken ? 'none' : 'flex';
      if (logoutBtn) logoutBtn.style.display = accessToken ? 'flex' : 'none';
    }

    async function loadGuilds() {
      if (loginBtn) loginBtn.style.display = 'none';
      if (logoutBtn) logoutBtn.style.display = 'flex';
      guildList.innerHTML = '<div style="color:var(--muted);font-size:.82rem;padding:8px"><div class="loading-spinner" style="margin-right:8px;vertical-align:middle"></div> Loading servers...</div>';

      try {
        const res = await fetch(API + '/api/user/guilds', {
          headers: accessToken ? { 'Authorization': 'Bearer ' + accessToken } : {}
        });
        if (!res.ok) throw new Error('Failed to fetch');
        const guilds = await res.json();

        if (!guilds || guilds.length === 0) {
          guildList.innerHTML = '<div style="color:var(--muted);font-size:.82rem;padding:8px">No mutual servers found. Make sure PulseKeep is in your server.</div>';
          return;
        }

        guildList.innerHTML = '';
        guilds.forEach(g => {
          const div = document.createElement('div');
          div.className = 'guild-item';
          if (!g.hasAdmin) div.classList.add('no-admin');
          const icon = g.icon
            ? 'https://cdn.discordapp.com/icons/' + g.id + '/' + g.icon + '.png?size=64'
            : 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"%3E%3Crect fill="%23333" width="32" height="32" rx="8"/%3E%3Ctext x="16" y="20" text-anchor="middle" fill="%23fff" font-size="16"%3E' + (g.name ? g.name.charAt(0).toUpperCase() : '?') + '%3C/text%3E%3C/svg%3E';
          div.innerHTML = '<img src="' + icon + '" alt="" loading="lazy"><span class="guild-name">' + (g.name || 'Unknown') + '</span>';
          if (!g.hasAdmin) {
            div.innerHTML += '<span class="guild-lock"><i class="fa-solid fa-lock" title="Need Manage Server permission"></i></span>';
          }
          div.addEventListener('click', () => {
            if (!g.hasAdmin) {
              toast('You need Manage Server permission to configure this server.', 'error');
              return;
            }
            qsa('.guild-item').forEach(el => el.classList.remove('active'));
            div.classList.add('active');
            if (sidebar) sidebar.classList.remove('open');
            loadConfig(g.id, g.name);
          });
          guildList.appendChild(div);
        });
      } catch {
        guildList.innerHTML = '<div style="color:var(--muted);font-size:.82rem;padding:8px">Failed to load servers. Make sure you\'re logged in.</div>';
      }
    }

    async function loadConfig(guildId, guildName) {
      $('dash-guild-name').textContent = guildName;
      $('dash-guild-id').textContent = 'ID: ' + guildId;
      if (dashPrompt) dashPrompt.style.display = 'none';
      if (dashConfig) dashConfig.style.display = 'block';
      const saved = $('dash-saved');
      if (saved) saved.style.display = 'none';

      // Show loading state
      qsa('[data-key]').forEach(el => {
        if (el.type === 'checkbox') el.checked = true;
        else el.value = '';
      });

      try {
        const res = await fetch(API + '/api/guild/' + guildId + '/config', {
          headers: accessToken ? { 'Authorization': 'Bearer ' + accessToken } : {}
        });
        if (!res.ok) throw new Error('No config');
        const cfg = await res.json();

        const vc = qs('[data-key="voteChannelId"]');
        if (vc) vc.value = cfg.voteChannelId || '';

        ['economyEnabled','ticketsEnabled','modlogsEnabled','welcomeEnabled'].forEach(key => {
          const el = qs('[data-key="' + key + '"]');
          if (el) el.checked = cfg[key] !== false;
        });

        const tc = qs('[data-key="ticketCategoryId"]');
        if (tc) tc.value = cfg.ticketCategoryId || '';

        const lc = qs('[data-key="logChannelId"]');
        if (lc) lc.value = cfg.logChannelId || '';

        const wc = qs('[data-key="welcomeChannelId"]');
        if (wc) wc.value = cfg.welcomeChannelId || '';

      } catch {
        // Keep defaults on error
      }
    }

    const saveBtn = $('dash-save');
    if (saveBtn) {
      saveBtn.addEventListener('click', async () => {
        const guildIdEl = $('dash-guild-id');
        if (!guildIdEl) return;
        const guildId = guildIdEl.textContent.replace('ID: ','');
        if (!guildId) return;

        const payload = {};
        qsa('[data-key]').forEach(el => {
          const key = el.getAttribute('data-key');
          payload[key] = el.type === 'checkbox' ? el.checked : el.value;
        });

        saveBtn.disabled = true;
        saveBtn.innerHTML = '<div class="loading-spinner" style="width:16px;height:16px;border-width:2px"></div> Saving...';

        try {
          const res = await fetch(API + '/api/guild/' + guildId + '/config', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + accessToken },
            body: JSON.stringify(payload),
          });
          if (res.ok) {
            toast('Configuration saved', 'success');
          } else {
            const err = await res.json().catch(() => ({}));
            toast(err.error || 'Failed to save', 'error');
          }
        } catch {
          toast('Failed to save config', 'error');
        } finally {
          saveBtn.disabled = false;
          saveBtn.innerHTML = '<i class="fa-solid fa-floppy-disk"></i> Save Changes';
        }
      });
    }

    checkAuth();

    if (loginBtn) {
      loginBtn.addEventListener('click', (e) => {
        e.preventDefault();
        window.location.href = API + '/auth/discord/login';
      });
    }

    if (logoutBtn) {
      logoutBtn.addEventListener('click', (e) => {
        e.preventDefault();
        accessToken = null;
        localStorage.removeItem('pk_token');
        updateAuthUI();
        guildList.innerHTML = '<div style="color:var(--muted);font-size:.82rem;padding:8px">Logged out. Log in to see your servers.</div>';
        if (dashConfig) dashConfig.style.display = 'none';
        if (dashPrompt) dashPrompt.style.display = 'block';
        toast('Logged out', 'success');
      });
    }
  }
});
