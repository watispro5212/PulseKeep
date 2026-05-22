/* ==========================================================
   PulseKeep Web App — Live Stats + Animations
   ========================================================== */

// Configuration — change this to your Fly.io app URL when deployed
const API_BASE_URL = window.location.hostname === 'localhost'
    ? 'http://localhost:8080'
    : 'https://pulsekeep.fly.dev';

// ── Animated Counter ──────────────────────────────────────
function animateCounter(element, target, suffix = '') {
    const isNumber = typeof target === 'number';
    if (!isNumber) {
        element.textContent = target;
        return;
    }

    const duration = 2000;
    const frameDuration = 1000 / 60;
    const totalFrames = Math.round(duration / frameDuration);
    let frame = 0;

    const counter = setInterval(() => {
        frame++;
        const progress = easeOutCubic(frame / totalFrames);
        const currentValue = Math.round(target * progress);

        element.textContent = formatNumber(currentValue) + suffix;

        if (frame === totalFrames) {
            clearInterval(counter);
            element.textContent = formatNumber(target) + suffix;
        }
    }, frameDuration);
}

function easeOutCubic(t) {
    return 1 - Math.pow(1 - t, 3);
}

function formatNumber(num) {
    return num.toLocaleString('en-US');
}

// ── Fetch Live Stats from API ─────────────────────────────
async function fetchStats() {
    const statusBadge = document.getElementById('bot-status-badge');
    const badgeText = statusBadge?.querySelector('.badge-text');
    const pulseDot = statusBadge?.querySelector('.pulse-dot');

    try {
        const response = await fetch(`${API_BASE_URL}/stats`, {
            signal: AbortSignal.timeout(8000),
        });

        if (!response.ok) throw new Error(`HTTP ${response.status}`);

        const data = await response.json();

        // Populate stat cards with animated counters
        animateCounter(document.getElementById('stat-servers'), data.servers || 0);
        animateCounter(document.getElementById('stat-users'), data.users || 0);
        animateCounter(document.getElementById('stat-commands'), data.commands_run || 0);
        
        const uptimeEl = document.getElementById('stat-uptime');
        if (uptimeEl) uptimeEl.textContent = data.uptime || '—';

        // Show online status
        if (pulseDot) pulseDot.classList.add('green');
        if (badgeText) badgeText.textContent = 'PulseKeep Service Online';

    } catch (err) {
        console.warn('Could not reach PulseKeep API:', err.message);

        // Fallback to offline state
        document.getElementById('stat-servers').textContent = '—';
        document.getElementById('stat-users').textContent = '—';
        document.getElementById('stat-commands').textContent = '—';
        document.getElementById('stat-uptime').textContent = '—';

        if (pulseDot) {
            pulseDot.classList.remove('green');
            pulseDot.style.backgroundColor = '#ef4444';
            pulseDot.style.boxShadow = '0 0 10px #ef4444';
        }
        if (badgeText) badgeText.textContent = 'PulseKeep Service Offline';
    }
}

// ── Intersection Observer for Scroll Animations ───────────
function initScrollAnimations() {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach((entry) => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
                observer.unobserve(entry.target);
            }
        });
    }, { threshold: 0.15 });

    document.querySelectorAll('.stat-card, .feature-card, .cta-container').forEach((el) => {
        el.classList.add('scroll-reveal');
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

// ── Header Glass Effect on Scroll ─────────────────────────
function initHeaderScroll() {
    const header = document.querySelector('.header');
    if (!header) return;

    window.addEventListener('scroll', () => {
        if (window.scrollY > 60) {
            header.style.backgroundColor = 'rgba(7, 8, 13, 0.92)';
            header.style.boxShadow = '0 4px 30px rgba(0, 0, 0, 0.4)';
        } else {
            header.style.backgroundColor = 'rgba(7, 8, 13, 0.7)';
            header.style.boxShadow = 'none';
        }
    });
}

// ── Initialize on DOM Ready ───────────────────────────────
document.addEventListener('DOMContentLoaded', () => {
    fetchStats();
    initScrollAnimations();
    initSmoothScroll();
    initHeaderScroll();

    // Refresh stats every 60 seconds
    setInterval(fetchStats, 60_000);
});
