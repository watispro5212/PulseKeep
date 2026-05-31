export default {
	async fetch(request, env, ctx) {
		const API_ORIGIN = env.API_ORIGIN;
		const url = new URL(request.url);

		if (url.pathname.startsWith('/api/') || url.pathname.startsWith('/auth/') || url.pathname === '/health' || url.pathname === '/stats') {
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

			const apiUrl = API_ORIGIN + url.pathname + url.search;
			const headers = new Headers(request.headers);
			headers.set('Accept', 'application/json');
			headers.set('Host', new URL(API_ORIGIN).host);

			const apiRequest = new Request(apiUrl, {
				method: request.method,
				headers: headers,
				body: request.method === 'GET' || request.method === 'HEAD' ? undefined : request.body,
				redirect: 'manual',
			});

			try {
				const response = await fetch(apiRequest);
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
				for (const [key, value] of Object.entries(corsHeaders())) {
					outHeaders.set(key, value);
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
		}

		const response = await env.ASSETS.fetch(request);
		if (response.status === 404) {
			const notFound = await env.ASSETS.fetch(new Request(url.origin + '/404.html'));
			return new Response(notFound.body, { status: 404, headers: notFound.headers });
		}
		return response;
	},
};

function corsHeaders() {
	return {
		'Access-Control-Allow-Origin': '*',
		'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
		'Access-Control-Allow-Headers': 'Content-Type, Authorization, Accept',
		'Cache-Control': 'no-store',
	};
}
