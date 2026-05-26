const API_ORIGIN = 'https://pulsekeep.fly.dev';

export default {
	async fetch(request, env, ctx) {
		const url = new URL(request.url);

		if (url.pathname === '/health' || url.pathname === '/stats') {
			const apiUrl = API_ORIGIN + url.pathname + url.search;
			const apiRequest = new Request(apiUrl, {
				method: request.method,
				headers: {
					'Accept': 'application/json',
				},
			});
			return fetch(apiRequest);
		}

		return env.ASSETS.fetch(request);
	},
};
