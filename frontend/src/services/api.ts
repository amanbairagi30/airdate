import { auth } from './auth';

const API_BASE_URL = 'http://localhost:8080/api';

export const api = {
  register: async (username: string, password: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/register`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        const data = await response.json().catch(() => ({
          error: `Registration failed with status ${response.status}`
        }));
        throw new Error(data.error || 'Registration failed');
      }

      const data = await response.json();
      // After successful registration, automatically log in
      if (data.token) {
        auth.login(data.token, data.username);
      }
      return data;
    } catch (error) {
      console.error('Registration error:', error);
      throw error;
    }
  },

  login: async (username: string, password: string) => {
    try {
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

      if (!response.ok) {
        const data = await response.json().catch(() => ({
          error: `Login failed with status ${response.status}`
        }));
        throw new Error(data.error || 'Login failed');
      }

      const data = await response.json();
      // Store the token and username in localStorage
      auth.login(data.token, data.username);
      return data;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  },

  getAllUsers: async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/users`, {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      });

      if (!response.ok) {
        const data = await response.json().catch(() => ({
          error: `Failed to fetch users with status ${response.status}`
        }));
        throw new Error(data.error);
      }

      const data = await response.json();
      return data || [];
    } catch (error) {
      console.error('Error fetching users:', error);
      throw error;
    }
  },

  getProfile: async () => {
    return fetchWithAuth('/profile');
  },

  connectTwitch: async (twitchUsername: string) => {
    return fetchWithAuth('/connect/twitch', {
      method: 'POST',
      body: JSON.stringify({ twitchUsername }),
    });
  },

  connectDiscord: async (discordUsername: string) => {
    return fetchWithAuth('/connect/discord', {
      method: 'POST',
      body: JSON.stringify({ discordUsername }),
    });
  },

  connectInstagram: async (instagramHandle: string) => {
    return fetchWithAuth('/connect/instagram', {
      method: 'POST',
      body: JSON.stringify({ instagramHandle }),
    });
  },

  connectYoutube: async (youtubeChannel: string) => {
    return fetchWithAuth('/connect/youtube', {
      method: 'POST',
      body: JSON.stringify({ youtubeChannel }),
    });
  },

  searchGames: async (query: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/games/search?q=${encodeURIComponent(query)}`, {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      });

      if (!response.ok) {
        const data = await response.json().catch(() => ({
          error: `Failed to search games with status ${response.status}`
        }));
        throw new Error(data.error || 'Failed to search games');
      }

      return await response.json();
    } catch (error) {
      console.error('Game search error:', error);
      throw error;
    }
  },

  connectGame: async (gameName: string, gameId: string) => {
    return fetchWithAuth('/connect/game', {
      method: 'POST',
      body: JSON.stringify({ gameName, gameId }),
    });
  },

  getUserProfile: async (username: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/users/${username}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${auth.getToken()}`,
        },
      });

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('User not found');
        }
        throw new Error(`Failed to fetch profile with status ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Profile fetch error:', error);
      throw error;
    }
  },

  updatePrivacySettings: async (isPrivate: boolean) => {
    try {
      const response = await fetch(`${API_BASE_URL}/privacy`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${auth.getToken()}`,
        },
        body: JSON.stringify({ isPrivate }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        let errorMessage;
        try {
          const errorData = JSON.parse(errorText);
          errorMessage = errorData.error || `Failed to update privacy settings with status ${response.status}`;
        } catch {
          errorMessage = `Failed to update privacy settings with status ${response.status}`;
        }
        throw new Error(errorMessage);
      }

      const text = await response.text();
      return text ? JSON.parse(text) : { message: 'Privacy settings updated', isPrivate };
    } catch (error) {
      console.error('Privacy update error:', error);
      throw error;
    }
  },

  disconnectInstagram: async () => {
    return fetchWithAuth('/disconnect/instagram', {
      method: 'POST'
    });
  },

  disconnectYoutube: async () => {
    return fetchWithAuth('/disconnect/youtube', {
      method: 'POST'
    });
  },

  disconnectGame: async (gameName: string) => {
    return fetchWithAuth('/disconnect/game', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ gameName }),
    });
  },
};

async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const token = auth.getToken();
  if (!token) {
    throw new Error('No authentication token found');
  }

  const headers = {
    ...options.headers,
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  };

  try {
    const response = await fetch(`${API_BASE_URL}${url}`, { 
      ...options, 
      headers,
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
