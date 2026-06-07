const API = {
  async request(url, options = {}) {
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    const authHeader = Auth.getAuthHeader();
    if (authHeader) {
      Object.assign(headers, authHeader);
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      if (response.status === 401) {
        Auth.logout();
        throw new Error('Сессия истекла. Пожалуйста, войдите заново.');
      }

      if (!response.ok) {
        const error = await this.parseError(response);
        throw error;
      }

      if (response.status === 204) {
        return { success: true };
      }

      return await response.json();
    } catch (error) {
      console.error('API request error:', error);
      throw error;
    }
  },

  async parseError(response) {
    try {
      const data = await response.json();
      return new Error(data.message || `Ошибка ${response.status}`);
    } catch {
      return new Error(`Ошибка ${response.status}: ${response.statusText}`);
    }
  },

  async getMenu() {
    return this.request(CONFIG.API.menuPublic());
  },

  async getMenuAdmin(filters = {}) {
    const params = new URLSearchParams(filters);
    const url = `${CONFIG.API.menuAdmin()}${params.toString() ? '?' + params.toString() : ''}`;
    return this.request(url);
  },

  async createMenuItem(data) {
    return this.request(CONFIG.API.menuAdmin(), {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async updateMenuItem(id, data) {
    return this.request(CONFIG.API.menuItem(id), {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  async deleteMenuItem(id) {
    return this.request(CONFIG.API.menuItem(id), {
      method: 'DELETE',
    });
  },

  async toggleMenuItemAvailability(id, available) {
    return this.request(CONFIG.API.menuItem(id), {
      method: 'PATCH',
      body: JSON.stringify({ available }),
    });
  },

  async getOrders(filters = {}) {
    const params = new URLSearchParams(filters);
    const url = `${CONFIG.API.orders()}${params.toString() ? '?' + params.toString() : ''}`;
    return this.request(url);
  },

  async getOrder(id) {
    return this.request(CONFIG.API.ordersItem(id));
  },

  async updateOrderStatus(id, status) {
    return this.request(CONFIG.API.ordersItem(id), {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    });
  },

  async uploadImage(file) {
    const formData = new FormData();
    formData.append('file', file);

    const headers = {};
    const authHeader = Auth.getAuthHeader();
    if (authHeader) {
      Object.assign(headers, authHeader);
    }

    try {
      const response = await fetch(CONFIG.API.upload(), {
        method: 'POST',
        headers,
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Ошибка загрузки файла');
      }

      return await response.json();
    } catch (error) {
      console.error('Upload error:', error);
      throw error;
    }
  },
};
