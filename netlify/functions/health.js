const apiBase = process.env.PULSEKEEP_API_BASE;

exports.handler = async () => {
    if (!apiBase) {
        return json(500, {
            status: 'unconfigured',
            database: 'unknown',
            uptime: '--',
            error: 'PULSEKEEP_API_BASE is not configured.',
        });
    }

    try {
        const response = await fetch(`${apiBase.replace(/\/$/, '')}/health`, {
            headers: { Accept: 'application/json' },
        });
        const data = await response.json();
        return json(response.status, data);
    } catch (error) {
        return json(502, {
            status: 'unavailable',
            database: 'unknown',
            uptime: '--',
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
