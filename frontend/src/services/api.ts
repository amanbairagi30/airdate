const API_BASE_URL = 'http://localhost:8080/api';

export async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const token = localStorage.getItem('token');
  const headers = {
    ...options.headers,
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };

  const response = await fetch(`${API_BASE_URL}${url}`, { ...options, headers });
  if (!response.ok) {
    throw new Error('API request failed');
  }
  return response.json();
}

export const api = {
  health: () => fetch(`${API_BASE_URL}/health`).then(res => res.json()),
  register: (username: string, password: string) => 
    fetch(`${API_BASE_URL}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    }).then(res => res.json()),
  login: (username: string, password: string) => 
    fetch(`${API_BASE_URL}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    }).then(res => res.json()),
  getProfile: () => fetchWithAuth('/profile'),
  connectTwitch: (twitchUsername: string) => 
    fetchWithAuth('/connect/twitch', {
      method: 'POST',
      body: JSON.stringify({ twitchUsername }),
    }),
  connectDiscord: (discordUsername: string) => 
    fetchWithAuth('/connect/discord', {
      method: 'POST',
      body: JSON.stringify({ discordUsername }),
    }),
  searchGames: (query: string) => 
    fetchWithAuth(`/games/search?q=${encodeURIComponent(query)}`),
  // Add more API calls here as needed
};
