const CONFIG = {
API_BASE_URL: (window.__API_BASE_URL__ || 'https://qwefsdfsdsg-mdk.hf.space') + '/api',
  
  ENDPOINTS: {
    LOGIN: '/admin/login',
    MENU_PUBLIC: '/menu',
    MENU_ADMIN: '/admin/menu',
    ORDERS: '/admin/orders',
    UPLOAD: '/admin/upload',
  },
  
  STORAGE_KEYS: {
    TOKEN: 'admin_token',
    USER: 'admin_user',
  },
  
  TOKEN_EXPIRY_HOURS: 24,
  ITEMS_PER_PAGE: 10,
  SEARCH_DEBOUNCE_MS: 300,
};

CONFIG.API = {
  login: () => `${CONFIG.API_BASE_URL}${CONFIG.ENDPOINTS.LOGIN}`,
  menuPublic: () => `${CONFIG.API_BASE_URL}${CONFIG.ENDPOINTS.MENU_PUBLIC}`,
  menuAdmin: () => `${CONFIG.API_BASE_URL}${CONFIG.ENDPOINTS.MENU_ADMIN}`,
  menuItem: (id) => `${CONFIG.API_BASE_URL}${CONFIG.ENDPOINTS.MENU_ADMIN}/${id}`,
  orders: () => `${CONFIG.API_BASE_URL}${CONFIG.ENDPOINTS.ORDERS}`,
  ordersItem: (id) => `${CONFIG.API_BASE_URL}${CONFIG.ENDPOINTS.ORDERS}/${id}`,
  upload: () => `${CONFIG.API_BASE_URL}${CONFIG.ENDPOINTS.UPLOAD}`,
};

CONFIG.MENU_CATEGORIES = [
  { value: 'snacks', label: 'Закуски' },
  { value: 'main', label: 'Основные блюда' },
  { value: 'desserts', label: 'Десерты' },
  { value: 'drinks', label: 'Напитки' },
  { value: 'salads', label: 'Салаты' },
];

CONFIG.ORDER_STATUSES = [
  { value: 'pending', label: 'В ожидании', color: 'warning' },
  { value: 'confirmed', label: 'Подтвержден', color: 'info' },
  { value: 'preparing', label: 'Готовится', color: 'primary' },
  { value: 'ready', label: 'Готов', color: 'success' },
  { value: 'delivered', label: 'Доставлен', color: 'success' },
  { value: 'cancelled', label: 'Отменен', color: 'danger' },
];
