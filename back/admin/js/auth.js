const Auth = {
  async login(email, password) {
    try {
      const response = await fetch(CONFIG.API.login(), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || 'Ошибка при входе');
      }

      const data = await response.json();
      
      if (!data.token) {
        throw new Error('Token не получен от сервера');
      }

      Storage.setToken(data.token);
      if (data.user) {
        Storage.setUser(data.user);
      }

      return data;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  },

  logout() {
    Storage.clear();
    window.location.href = '/admin/login.html';
  },

  isAuthenticated() {
    return Storage.hasToken();
  },

  getCurrentUser() {
    return Storage.getUser();
  },

  getAuthHeader() {
    const token = Storage.getToken();
    if (!token) return null;
    return {
      'Authorization': `Bearer ${token}`,
    };
  },

  protectRoute() {
    if (!this.isAuthenticated()) {
      window.location.href = '/admin/login.html';
    }
  },
};
