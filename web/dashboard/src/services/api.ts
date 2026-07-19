const API_BASE = 'http://localhost:8083/api/v1';

let cachedToken: string | null = null;

async function getToken(): Promise<string> {
    if (cachedToken) return cachedToken;
    
    // Auto-fetch dev token
    const res = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST'
    });
    const data = await res.json();
    cachedToken = data.token;
    return data.token;
}

async function fetchWithAuth(endpoint: string, options: RequestInit = {}) {
    const token = await getToken();
    const headers = {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        ...options.headers,
    };
    
    const res = await fetch(`${API_BASE}${endpoint}`, { ...options, headers });
    if (!res.ok) {
        throw new Error(`API error: ${res.status}`);
    }
    return res.json();
}

export const DashboardService = {
    getMetrics: () => fetchWithAuth('/dashboard/metrics'),
    getArtifacts: (limit = 50, offset = 0) => fetchWithAuth(`/artifacts?limit=${limit}&offset=${offset}`)
};
