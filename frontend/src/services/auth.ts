export const auth = {
  isAuthenticated: () => {
    return !!localStorage.getItem('token');
  },

  getToken: () => {
    return localStorage.getItem('token');
  },

  login: (token: string, username: string) => {
    localStorage.setItem('token', token);
    localStorage.setItem('username', username);
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
  }
}; 