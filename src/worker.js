const API_ORIGIN = 'https://pulsekeep.fly.dev';

export default {
	async fetch(request, env, ctx) {
		const url = new URL(request.url);

		if (url.pathname.startsWith('/api/') || url.pathname.startsWith('/auth/') || url.pathname === '/health' || url.pathname === '/stats') {
			const apiUrl = API_ORIGIN + url.pathname + url.search;
			const headers = new Headers(request.headers);
			headers.set('Accept', 'application/json');

			const apiRequest = new Request(apiUrl, {
				method: request.method,
				headers: headers,
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
							headers: { Location: newLocation },
						});
					}
				}
				return response;
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
					headers: {
						'Content-Type': 'application/json',
						'Access-Control-Allow-Origin': '*',
						'Access-Control-Allow-Methods': 'GET, OPTIONS',
						'Access-Control-Allow-Headers': '*',
					},
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
