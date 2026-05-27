const API_BASE_URL = window.location.origin;

const statEls = {
    servers: document.getElementById('stat-servers'),
    users: document.getElementById('stat-users'),
    commands: document.getElementById('stat-commands'),
    uptime: document.getElementById('stat-uptime'),
    apiSpeed: document.getElementById('api-speed'),
    database: document.getElementById('stat-database'),
    version: document.getElementById('stat-version'),
    runtime: document.getElementById('stat-runtime'),
};

function setStatus(state, message) {
    const badge = document.getElementById('bot-status-badge');
    if (!badge) return;
    const dot = badge.querySelector('.pulse-dot');
    const text = badge.querySelector('.badge-text') || badge.querySelector('[data-status-text]');
    if (text) text.textContent = message;
    if (dot) {
        dot.className = 'pulse-dot';
        if (state) dot.classList.add(state);
    }
    // Preserve status-badge class (used on status page), add/remove state class
    badge.classList.remove('online', 'offline', 'degraded');
    if (state) badge.classList.add(state);
    if (state === 'online') {
        const icon = badge.querySelector('i');
        if (icon) icon.className = 'fa-solid fa-circle-check';
    }
}

function animateCounter(element, target) {
    if (!element) return;
    const value = Number(target);
    if (!Number.isFinite(value)) {
        element.textContent = '--';
        return;
    }
    const duration = 1000;
    const start = performance.now();
    function tick(now) {
        const progress = Math.min((now - start) / duration, 1);
        const eased = 1 - Math.pow(1 - progress, 3);
        element.textContent = Math.round(value * eased).toLocaleString('en-US');
        if (progress < 1) requestAnimationFrame(tick);
    }
    requestAnimationFrame(tick);
}

function setFallback() {
    if (statEls.servers) statEls.servers.textContent = '--';
    if (statEls.users) statEls.users.textContent = '--';
    if (statEls.commands) statEls.commands.textContent = '--';
    if (statEls.uptime) statEls.uptime.textContent = '--';
    if (statEls.apiSpeed) statEls.apiSpeed.textContent = '-- ms';
    if (statEls.database) statEls.database.textContent = 'Unknown';
    if (statEls.version) statEls.version.textContent = '--';
    if (statEls.runtime) statEls.runtime.textContent = 'Offline';
}

async function timedFetch(path, options = {}) {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 8000);
    const startedAt = performance.now();
    try {
        const response = await fetch(`${API_BASE_URL.replace(/\/$/, '')}${path}`, {
            ...options,
            signal: controller.signal,
            headers: { Accept: 'application/json', ...(options.headers || {}) },
        });
        return { response, duration: Math.round(performance.now() - startedAt) };
    } finally {
        clearTimeout(timeout);
    }
}

async function fetchHealth() {
    try {
        const { response, duration } = await timedFetch('/health');
        if (!response.ok) throw new Error(`Health request failed with HTTP ${response.status}`);
        const data = await response.json();
        if (statEls.apiSpeed) statEls.apiSpeed.textContent = `${duration} ms`;
        if (statEls.database) statEls.database.textContent = data.database === 'ok' ? 'Healthy' : data.database || 'Unknown';
        if (statEls.runtime) statEls.runtime.textContent = data.uptime || 'Online';
        return data;
    } catch (error) {
        if (statEls.apiSpeed) statEls.apiSpeed.textContent = '-- ms';
        if (statEls.database) statEls.database.textContent = 'Unavailable';
        throw error;
    }
}

async function fetchStats() {
    if (!API_BASE_URL) {
        setFallback();
        setStatus('offline', 'API not configured');
        return;
    }

    try {
        const { response } = await timedFetch('/stats');
        if (!response.ok) throw new Error(`Stats request failed with HTTP ${response.status}`);
        const data = await response.json();

        animateCounter(statEls.servers, data.servers);
        animateCounter(statEls.users, data.users);
        animateCounter(statEls.commands, data.commands_run);
        if (statEls.uptime) statEls.uptime.textContent = data.uptime || '--';
        if (statEls.version) statEls.version.textContent = (data.bot || '').replace(/^PulseKeep\s*/i, '') || 'v5.8';

        const healthOk = await fetchHealth().catch(() => null);

        setStatus('online', 'All systems online');
    } catch (error) {
        console.warn('Could not reach PulseKeep API:', error.message);
        setFallback();
        setStatus('offline', 'Service offline');
    }
}

function initMobileMenu() {
    const toggle = document.getElementById('mobile-menu-toggle');
    const mobileNav = document.getElementById('mobile-nav');
    if (!toggle || !mobileNav) return;
    toggle.addEventListener('click', () => {
        const isOpen = mobileNav.classList.toggle('active');
        toggle.setAttribute('aria-expanded', String(isOpen));
        const icon = toggle.querySelector('i');
        if (icon) icon.className = isOpen ? 'fa-solid fa-xmark' : 'fa-solid fa-bars';
    });
}

function initSmoothScroll() {
    document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
        anchor.addEventListener('click', (event) => {
            const targetId = anchor.getAttribute('href');
            const target = targetId ? document.querySelector(targetId) : null;
            if (!target) return;
            event.preventDefault();
            target.scrollIntoView({ behavior: 'smooth', block: 'start' });
        });
    });
}

function initCopyCommands() {
    document.querySelectorAll('.command-item code').forEach(el => {
        el.addEventListener('click', () => {
            const text = el.textContent.replace(/[<>\[\]]/g, '').split(' ')[0];
            navigator.clipboard.writeText(text).catch(() => {});
            const original = el.textContent;
            el.textContent = 'Copied!';
            setTimeout(() => { el.textContent = original; }, 1200);
        });
        el.style.cursor = 'pointer';
    });
}

function initScrollAnimations() {
    if (!('IntersectionObserver' in window)) return;
    const observer = new IntersectionObserver((entries) => {
        entries.forEach((entry) => {
            if (!entry.isIntersecting) return;
            entry.target.classList.add('visible');
            observer.unobserve(entry.target);
        });
    }, { threshold: 0.1 });
    document.querySelectorAll('.scroll-reveal').forEach((element) => {
        observer.observe(element);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    initMobileMenu();
    initSmoothScroll();
    initScrollAnimations();
    initCopyCommands();
    fetchStats();
    window.setInterval(fetchStats, 60_000);
});
