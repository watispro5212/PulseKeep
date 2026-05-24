const API_BASE_URL = '/.netlify/functions';

const statEls = {
    servers: document.getElementById('stat-servers'),
    users: document.getElementById('stat-users'),
    commands: document.getElementById('stat-commands'),
    uptime: document.getElementById('stat-uptime'),
    latency: document.getElementById('stat-latency'),
    apiSpeed: document.getElementById('api-speed'),
    database: document.getElementById('stat-database'),
};

function setStatus(state, message) {
    document.querySelectorAll('#bot-status-badge .badge-text, [data-status-text]').forEach((element) => {
        element.textContent = message;
    });

    document.querySelectorAll('#bot-status-badge .pulse-dot, [data-status-dot]').forEach((dot) => {
        dot.classList.remove('online', 'offline');
        if (state) dot.classList.add(state);
    });
}

function animateCounter(element, target) {
    if (!element) return;
    const value = Number(target);

    if (!Number.isFinite(value)) {
        element.textContent = '--';
        return;
    }

    const duration = 800;
    const start = performance.now();

    function tick(now) {
        const progress = Math.min((now - start) / duration, 1);
        const eased = 1 - Math.pow(1 - progress, 3);
        element.textContent = Math.round(value * eased).toLocaleString('en-US');

        if (progress < 1) {
            requestAnimationFrame(tick);
        }
    }

    requestAnimationFrame(tick);
}

function setFallback() {
    ['servers', 'users', 'commands', 'uptime', 'latency', 'apiSpeed'].forEach((key) => {
        if (statEls[key]) statEls[key].textContent = key === 'latency' || key === 'apiSpeed' ? '-- ms' : '--';
    });
    if (statEls.database) statEls.database.textContent = 'Unknown';
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
        setStatus('offline', 'PulseKeep API not configured');
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
        if (statEls.latency) statEls.latency.textContent = data.latency_ms ? `${data.latency_ms} ms` : '-- ms';

        await fetchHealth().catch(() => null);
        setStatus('online', 'PulseKeep service online');
    } catch (error) {
        console.warn('Could not reach PulseKeep API:', error.message);
        setFallback();
        setStatus('offline', 'PulseKeep service offline');
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

function initScrollAnimations() {
    if (!('IntersectionObserver' in window)) return;

    const observer = new IntersectionObserver((entries) => {
        entries.forEach((entry) => {
            if (!entry.isIntersecting) return;
            entry.target.classList.add('visible');
            observer.unobserve(entry.target);
        });
    }, { threshold: 0.12 });

    document.querySelectorAll('.stat-card, .feature-card, .command-category, .command-item, .setup-steps li, .metric-card, .team-card, .listing-card, .support-panel').forEach((element) => {
        element.classList.add('scroll-reveal');
        observer.observe(element);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    initMobileMenu();
    initSmoothScroll();
    initScrollAnimations();
    fetchStats();
    window.setInterval(fetchStats, 60_000);
});
