const isBrowser = typeof window !== 'undefined';

export const auth = {
  isAuthenticated: () => {
    if (!isBrowser) return false;
    return !!localStorage.getItem('token');
  },

  getToken: () => {
    if (!isBrowser) return null;
    return localStorage.getItem('token');
  },

  getUsername: () => {
    if (!isBrowser) return null;
    return localStorage.getItem('username');
  },

  login: (token: string, username: string) => {
    if (!isBrowser) return;
    localStorage.setItem('token', token);
    localStorage.setItem('username', username);
  },

  logout: () => {
    if (!isBrowser) return;
    localStorage.removeItem('token');
    localStorage.removeItem('username');
  }
}; 