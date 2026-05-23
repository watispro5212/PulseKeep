const apiMeta = document.querySelector('meta[name="pulsekeep-api"]');
const configuredApi = apiMeta?.getAttribute('content')?.trim();
const isLocalHost = ['localhost', '127.0.0.1', '::1'].includes(window.location.hostname);
const API_BASE_URL = isLocalHost ? 'http://localhost:8080' : configuredApi;

const statEls = {
    servers: document.getElementById('stat-servers'),
    users: document.getElementById('stat-users'),
    commands: document.getElementById('stat-commands'),
    uptime: document.getElementById('stat-uptime'),
};

function setStatus(state, message) {
    const badgeText = document.querySelector('#bot-status-badge .badge-text');
    const dot = document.querySelector('#bot-status-badge .pulse-dot');

    if (badgeText) badgeText.textContent = message;
    if (!dot) return;

    dot.classList.remove('online', 'offline');
    if (state) dot.classList.add(state);
}

function animateCounter(element, target) {
    if (!element) return;
    const value = Number(target);

    if (!Number.isFinite(value)) {
        element.textContent = '--';
        return;
    }

    const duration = 900;
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

function setStatsFallback() {
    Object.values(statEls).forEach((element) => {
        if (element) element.textContent = '--';
    });
}

async function fetchStats() {
    if (!API_BASE_URL) {
        setStatsFallback();
        setStatus('offline', 'PulseKeep API not configured');
        return;
    }

    try {
        const controller = new AbortController();
        const timeout = setTimeout(() => controller.abort(), 8000);
        const response = await fetch(`${API_BASE_URL.replace(/\/$/, '')}/stats`, {
            signal: controller.signal,
            headers: { Accept: 'application/json' },
        });
        clearTimeout(timeout);

        if (!response.ok) {
            throw new Error(`Stats request failed with HTTP ${response.status}`);
        }

        const data = await response.json();
        animateCounter(statEls.servers, data.servers);
        animateCounter(statEls.users, data.users);
        animateCounter(statEls.commands, data.commands_run);
        if (statEls.uptime) statEls.uptime.textContent = data.uptime || '--';
        setStatus('online', 'PulseKeep service online');
    } catch (error) {
        console.warn('Could not reach PulseKeep API:', error.message);
        setStatsFallback();
        setStatus('offline', 'PulseKeep service offline');
    }
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

    document.querySelectorAll('.stat-card, .feature-card, .command-group, .setup-steps li').forEach((element) => {
        element.classList.add('scroll-reveal');
        observer.observe(element);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    fetchStats();
    initSmoothScroll();
    initScrollAnimations();
    window.setInterval(fetchStats, 60_000);
});
