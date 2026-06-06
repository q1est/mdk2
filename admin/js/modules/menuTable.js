const MenuTable = {
  state: {
    allItems: [],
    filteredItems: [],
    searchQuery: '',
    currentCategory: null,
  },

  async init(tableId, tableBodyId) {
    this.tableElement = document.getElementById(tableId);
    this.tableBody = document.getElementById(tableBodyId);
    
    if (!this.tableBody || !this.tableElement) {
      console.error('Table elements not found');
      return;
    }

    this.setupEventListeners();
    await this.loadMenu();
  },

  async loadMenu() {
    try {
      UI.showTableSkeleton(this.tableBody);
      
      const data = await API.getMenu();
      this.state.allItems = data.items || data || [];
      this.renderTable();
    } catch (error) {
      console.error('Error loading menu:', error);
      Utils.showError('Ошибка загрузки меню: ' + error.message);
      UI.showEmptyState(this.tableBody, 'Ошибка загрузки меню');
    }
  },

  renderTable() {
    if (this.state.filteredItems.length === 0) {
      UI.showEmptyState(this.tableBody, 'Блюд не найдено');
      return;
    }

    this.tableBody.innerHTML = '';
    this.state.filteredItems.forEach(item => {
      const row = UI.createMenuRow(item);
      
      row.querySelector('.btn-edit')?.addEventListener('click', (e) => {
        e.preventDefault();
        this.onEditClick(item.id);
      });

      row.querySelector('.btn-delete')?.addEventListener('click', (e) => {
        e.preventDefault();
        this.onDeleteClick(item.id, item.name);
      });

      this.tableBody.appendChild(row);
    });
  },

  filter(searchQuery = '', category = null) {
    this.state.searchQuery = searchQuery.toLowerCase();
    this.state.currentCategory = category;

    this.state.filteredItems = this.state.allItems.filter(item => {
      const matchesSearch = !this.state.searchQuery || 
        item.name.toLowerCase().includes(this.state.searchQuery) ||
        (item.description && item.description.toLowerCase().includes(this.state.searchQuery));

      const matchesCategory = !this.state.currentCategory || 
        item.category === this.state.currentCategory;

      return matchesSearch && matchesCategory;
    });

    this.renderTable();
  },

  onSearch(e) {
    const query = e.target.value;
    this.filter(query, this.state.currentCategory);
  },

  onCategoryFilter(category) {
    this.filter(this.state.searchQuery, category);
  },

  onEditClick(itemId) {
    const item = this.state.allItems.find(i => i.id === itemId);
    if (item && window.MenuForm) {
      window.MenuForm.openEditModal(item);
    }
  },

  async onDeleteClick(itemId, itemName) {
    if (!Utils.confirmDelete(itemName)) return;

    try {
      Utils.showNotification('Удаление...', 'info', 0);
      await API.deleteMenuItem(itemId);
      
      Utils.showSuccess(`Блюдо "${itemName}" удалено`);
      this.state.allItems = this.state.allItems.filter(i => i.id !== itemId);
      this.filter(this.state.searchQuery, this.state.currentCategory);
    } catch (error) {
      Utils.showError('Ошибка удаления: ' + error.message);
    }
  },

  async toggleAvailability(itemId, currentStatus) {
    try {
      await API.toggleMenuItemAvailability(itemId, !currentStatus);
      
      const item = this.state.allItems.find(i => i.id === itemId);
      if (item) {
        item.available = !currentStatus;
        this.renderTable();
        Utils.showSuccess('Статус обновлен');
      }
    } catch (error) {
      Utils.showError('Ошибка обновления статуса: ' + error.message);
    }
  },

  setupEventListeners() {
    const searchInput = document.getElementById('menu-search');
    if (searchInput) {
      searchInput.addEventListener('input', this.onSearch.bind(this));
    }

    const categorySelect = document.getElementById('menu-category-filter');
    if (categorySelect) {
      categorySelect.addEventListener('change', (e) => {
        const category = e.target.value || null;
        this.onCategoryFilter(category);
      });
    }

    const addButton = document.getElementById('menu-add-button');
    if (addButton) {
      addButton.addEventListener('click', () => {
        if (window.MenuForm) {
          window.MenuForm.openCreateModal();
        }
      });
    }

    const refreshButton = document.getElementById('menu-refresh-button');
    if (refreshButton) {
      refreshButton.addEventListener('click', () => this.loadMenu());
    }
  },

  addItemToTable(item) {
    this.state.allItems.unshift(item);
    this.filter(this.state.searchQuery, this.state.currentCategory);
  },

  updateItemInTable(itemId, updatedItem) {
    const index = this.state.allItems.findIndex(i => i.id === itemId);
    if (index !== -1) {
      this.state.allItems[index] = { ...this.state.allItems[index], ...updatedItem };
      this.filter(this.state.searchQuery, this.state.currentCategory);
    }
  },
};
