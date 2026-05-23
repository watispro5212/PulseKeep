<<<<<<< HEAD
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
=======
/* ==========================================================
   PulseKeep Web App — Live Stats + Animations + UI Interactivity
   ========================================================== */

// Configuration — dynamically use current host for API calls
const API_BASE_URL = window.location.origin;

// ── Mobile Menu Toggle ─────────────────────────────────────
function initMobileMenu() {
    const toggle = document.getElementById('mobile-menu-toggle');
    const mobileNav = document.getElementById('mobile-nav');
    
    if (!toggle || !mobileNav) return;
    
    toggle.addEventListener('click', () => {
        mobileNav.classList.toggle('active');
        const icon = toggle.querySelector('i');
        if (mobileNav.classList.contains('active')) {
            icon.className = 'fa-solid fa-xmark';
        } else {
            icon.className = 'fa-solid fa-bars';
        }
    });
    
    // Close menu when clicking a link
    mobileNav.querySelectorAll('a').forEach(link => {
        link.addEventListener('click', () => {
            mobileNav.classList.remove('active');
            toggle.querySelector('i').className = 'fa-solid fa-bars';
        });
    });
}

// ── Header Scroll Effect ──────────────────────────────────
function initHeaderScroll() {
    const header = document.getElementById('main-header');
    if (!header) return;
    
    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            header.classList.add('scrolled');
        } else {
            header.classList.remove('scrolled');
        }
    });
}
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715

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
<<<<<<< HEAD
        animateCounter(statEls.servers, data.servers);
        animateCounter(statEls.users, data.users);
        animateCounter(statEls.commands, data.commands_run);
        if (statEls.uptime) statEls.uptime.textContent = data.uptime || '--';
        setStatus('online', 'PulseKeep service online');
    } catch (error) {
        console.warn('Could not reach PulseKeep API:', error.message);
        setStatsFallback();
        setStatus('offline', 'PulseKeep service offline');
=======

        // Populate stat cards with animated counters
        animateCounter(document.getElementById('stat-servers'), data.servers || 0);
        animateCounter(document.getElementById('stat-users'), data.users || 0);
        animateCounter(document.getElementById('stat-commands'), data.commands_run || 0);
        
        const uptimeEl = document.getElementById('stat-uptime');
        if (uptimeEl) uptimeEl.textContent = data.uptime || '—';

        // Update status page specific elements
        const latencyEl = document.getElementById('stat-latency');
        if (latencyEl) latencyEl.textContent = `${data.latency || 0} ms`;

        const apiSpeedEl = document.getElementById('api-speed');
        if (apiSpeedEl) {
            const startTime = performance.now();
            try {
                await fetch(`${API_BASE_URL}/health`);
                const duration = Math.round(performance.now() - startTime);
                apiSpeedEl.textContent = `${duration} ms`;
            } catch (e) {
                apiSpeedEl.textContent = '— ms';
            }
        }

        // Show online status
        if (pulseDot) pulseDot.classList.remove('offline');
        if (statusBadge) statusBadge.classList.remove('offline');
        if (badgeText) badgeText.textContent = 'PulseKeep Service Online';

    } catch (err) {
        console.warn('Could not reach PulseKeep API:', err.message);

        // Fallback to offline state
        document.getElementById('stat-servers').textContent = '—';
        document.getElementById('stat-users').textContent = '—';
        document.getElementById('stat-commands').textContent = '—';
        document.getElementById('stat-uptime').textContent = '—';

        if (document.getElementById('stat-latency')) document.getElementById('stat-latency').textContent = '— ms';
        if (document.getElementById('api-speed')) document.getElementById('api-speed').textContent = '— ms';

        if (pulseDot) {
            pulseDot.classList.add('offline');
            pulseDot.style.backgroundColor = 'var(--red-error)';
            pulseDot.style.boxShadow = '0 0 15px var(--red-glow)';
        }
        if (statusBadge) statusBadge.classList.add('offline');
        if (badgeText) badgeText.textContent = 'PulseKeep Service Offline';
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
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
<<<<<<< HEAD
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
=======
    }, { threshold: 0.1 });

    document.querySelectorAll('.scroll-reveal').forEach((el) => {
        observer.observe(el);
    });
}

// ── Smooth Scroll for Nav Links ───────────────────────────
function initSmoothScroll() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }
        });
    });
}

// ── Initialize on DOM Ready ───────────────────────────────
document.addEventListener('DOMContentLoaded', () => {
    initMobileMenu();
    initHeaderScroll();
    initScrollAnimations();
    initSmoothScroll();
    fetchStats();

    // Refresh stats every 60 seconds
    setInterval(fetchStats, 60_000);
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
});
