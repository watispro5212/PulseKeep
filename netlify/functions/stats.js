const apiBase = process.env.PULSEKEEP_API_BASE;

exports.handler = async () => {
    if (!apiBase) {
        return json(500, {
            bot: 'PulseKeep',
            status: 'unconfigured',
            servers: 0,
            users: 0,
            commands_run: 0,
            uptime: '--',
            latency_ms: 0,
            error: 'PULSEKEEP_API_BASE is not configured.',
        });
    }

    try {
        const response = await fetch(`${apiBase.replace(/\/$/, '')}/stats`, {
            headers: { Accept: 'application/json' },
        });
        const data = await response.json();
        return json(response.status, data);
    } catch (error) {
        return json(502, {
            bot: 'PulseKeep',
            status: 'unavailable',
            servers: 0,
            users: 0,
            commands_run: 0,
            uptime: '--',
            latency_ms: 0,
            error: 'PulseKeep backend is unavailable.',
        });
    }
};

function json(statusCode, body) {
    return {
        statusCode,
        headers: {
            'Cache-Control': 'no-store',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    };
}
