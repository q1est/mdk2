const Storage = {
  setToken(token) {
    if (!token) {
      console.warn('Token is empty');
      return false;
    }
    localStorage.setItem(CONFIG.STORAGE_KEYS.TOKEN, token);
    return true;
  },

  getToken() {
    return localStorage.getItem(CONFIG.STORAGE_KEYS.TOKEN) || null;
  },

  hasToken() {
    return !!this.getToken();
  },

  removeToken() {
    localStorage.removeItem(CONFIG.STORAGE_KEYS.TOKEN);
  },

  setUser(user) {
    if (!user) return false;
    localStorage.setItem(CONFIG.STORAGE_KEYS.USER, JSON.stringify(user));
    return true;
  },

  getUser() {
    const user = localStorage.getItem(CONFIG.STORAGE_KEYS.USER);
    return user ? JSON.parse(user) : null;
  },

  removeUser() {
    localStorage.removeItem(CONFIG.STORAGE_KEYS.USER);
  },

  clear() {
    this.removeToken();
    this.removeUser();
  },
};
