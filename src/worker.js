export default {
	async fetch(request, env, ctx) {
		const API_ORIGIN = normalizeOrigin(env.API_ORIGIN);
		const url = new URL(request.url);

		if (!API_ORIGIN) {
			return new Response(JSON.stringify({
				status: 'offline',
				error: 'Backend origin is not configured',
			}), {
				status: 503,
				headers: { 'Content-Type': 'application/json', ...corsHeaders() },
			});
		}

		if (request.method === 'OPTIONS') {
			return new Response(null, {
				status: 204,
				headers: corsHeaders(),
			});
		}

		const originUrl = API_ORIGIN + url.pathname + url.search;
		const headers = new Headers(request.headers);
		headers.set('Host', new URL(API_ORIGIN).host);

		if (expectsJson(url.pathname)) {
			headers.set('Accept', 'application/json');
		}

		const originRequest = new Request(originUrl, {
			method: request.method,
			headers: headers,
			body: request.method === 'GET' || request.method === 'HEAD' ? undefined : request.body,
			redirect: 'manual',
		});

		try {
			const response = await fetch(originRequest);
			if (response.status >= 300 && response.status < 400) {
				const location = response.headers.get('Location');
				if (location && !location.startsWith('http')) {
					const newLocation = url.origin + location;
					return new Response(response.body, {
						status: response.status,
						statusText: response.statusText,
						headers: { ...corsHeaders(), Location: newLocation },
					});
				}
			}

			const outHeaders = new Headers(response.headers);
			if (expectsJson(url.pathname)) {
				for (const [key, value] of Object.entries(corsHeaders())) {
					outHeaders.set(key, value);
				}
			}
			return new Response(response.body, {
				status: response.status,
				statusText: response.statusText,
				headers: outHeaders,
			});
		} catch (err) {
			return new Response(JSON.stringify({
				status: 'offline',
				database: 'unavailable',
				servers: 0,
				users: 0,
				commands: 0,
				commands_run: 0,
				uptime: 'unknown',
				bot_uptime: 'unknown',
				go_version: 'unknown',
				error: 'Backend unreachable',
			}), {
				status: 503,
				headers: { 'Content-Type': 'application/json', ...corsHeaders() },
			});
		}
	},
};

function normalizeOrigin(origin) {
	if (!origin) return '';
	return origin.replace(/\/+$/, '');
}

function expectsJson(pathname) {
	return pathname.startsWith('/api/') || pathname.startsWith('/auth/') || pathname === '/health' || pathname === '/stats';
}

function corsHeaders() {
	return {
		'Access-Control-Allow-Origin': '*',
		'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
		'Access-Control-Allow-Headers': 'Content-Type, Authorization, Accept',
		'Cache-Control': 'no-store',
	};
}
