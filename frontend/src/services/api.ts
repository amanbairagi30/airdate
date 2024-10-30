const API_BASE_URL = 'http://localhost:8080/api';

export async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('No authentication token found');
  }

  const headers = {
    ...options.headers,
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };

  try {
    const response = await fetch(`${API_BASE_URL}${url}`, { 
      ...options, 
      headers 
    });

    if (!response.ok) {
      const data = await response.json().catch(() => ({}));
      throw new Error(data.error || `Request failed with status ${response.status}`);
    }

    return response.json();
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
}

export const api = {
  health: () => fetch(`${API_BASE_URL}/health`).then(res => res.json()),
  
  register: async (username: string, password: string) => {
    const response = await fetch(`${API_BASE_URL}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'Registration failed');
    }
    
    return data;
  },
  
  login: async (username: string, password: string) => {
    const response = await fetch(`${API_BASE_URL}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.trim(),
        password: password.trim(),
      }),
    });

    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || 'Login failed');
    }
    
    return data;
  },
  
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
  connectInstagram: (instagramHandle: string) => 
    fetchWithAuth('/connect/instagram', {
      method: 'POST',
      body: JSON.stringify({ instagramHandle }),
    }),
  connectYoutube: (youtubeChannel: string) => 
    fetchWithAuth('/connect/youtube', {
      method: 'POST',
      body: JSON.stringify({ youtubeChannel }),
    }),
  searchGames: async (query: string) => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No authentication token found');
      }

      const response = await fetch(`${API_BASE_URL}/games/search?q=${encodeURIComponent(query)}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error(`Search failed with status ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Search games error:', error);
      throw error;
    }
  },
  disconnectTwitch: () => 
    fetchWithAuth('/disconnect/twitch', {
      method: 'POST'
    }),
  disconnectDiscord: () => 
    fetchWithAuth('/disconnect/discord', {
      method: 'POST'
    }),
  connectGame: (gameName: string, gameId: string) => 
    fetchWithAuth('/connect/game', {
      method: 'POST',
      body: JSON.stringify({ gameName, gameId }),
    }),
  // Add more API calls here as needed
};
