const MenuForm = {
  state: {
    mode: 'create',
    currentItemId: null,
    isLoading: false,
  },

  init(modalId, formId) {
    this.modal = document.getElementById(modalId);
    this.form = document.getElementById(formId);
    this.modalElement = new bootstrap.Modal(this.modal);

    if (!this.form) {
      console.error('Form element not found');
      return;
    }

    this.setupEventListeners();
    this.setupCategorySelect();
  },

  setupCategorySelect() {
    const categorySelect = this.form.querySelector('[name="category"]');
    if (!categorySelect) return;

    CONFIG.MENU_CATEGORIES.forEach(cat => {
      const option = document.createElement('option');
      option.value = cat.value;
      option.textContent = cat.label;
      categorySelect.appendChild(option);
    });
  },

  openCreateModal() {
    this.state.mode = 'create';
    this.state.currentItemId = null;
    
    const title = this.modal.querySelector('.modal-title');
    if (title) title.textContent = 'Добавить блюдо';

    this.form.reset();
    this.form.querySelector('[name="available"]').checked = true;

    this.modalElement.show();
  },

  openEditModal(item) {
    this.state.mode = 'edit';
    this.state.currentItemId = item.id;

    const title = this.modal.querySelector('.modal-title');
    if (title) title.textContent = 'Редактировать блюдо';

    this.form.querySelector('[name="name"]').value = item.name || '';
    this.form.querySelector('[name="description"]').value = item.description || '';
    this.form.querySelector('[name="price"]').value = item.price || '';
    this.form.querySelector('[name="category"]').value = item.category || '';
    this.form.querySelector('[name="image_url"]').value = item.image_url || '';
    this.form.querySelector('[name="available"]').checked = item.available !== false;

    this.modalElement.show();
  },

  setupEventListeners() {
    const submitButton = this.form.querySelector('button[type="submit"]');
    if (submitButton) {
      this.form.addEventListener('submit', (e) => this.onSubmit(e));
    }
  },

  async onSubmit(e) {
    e.preventDefault();

    if (this.state.isLoading) return;

    if (!this.validate()) {
      return;
    }

    try {
      this.state.isLoading = true;
      const submitButton = this.form.querySelector('button[type="submit"]');
      Utils.setButtonLoading(submitButton, true);

      const formData = Utils.getFormData(this.form);
      
      const data = {
        name: formData.name,
        description: formData.description,
        price: parseFloat(formData.price),
        category: formData.category,
        image_url: formData.image_url,
        available: formData.available === 'on' || formData.available === true,
      };

      if (this.state.mode === 'create') {
        await this.handleCreate(data);
      } else {
        await this.handleEdit(data);
      }

      this.modalElement.hide();
    } catch (error) {
      Utils.showError('Ошибка: ' + error.message);
      console.error('Form error:', error);
    } finally {
      this.state.isLoading = false;
      const submitButton = this.form.querySelector('button[type="submit"]');
      Utils.setButtonLoading(submitButton, false);
    }
  },

  async handleCreate(data) {
    const result = await API.createMenuItem(data);
    Utils.showSuccess('Блюдо создано успешно');
    
    if (window.MenuTable) {
      MenuTable.addItemToTable(result.item || result);
    }
  },

  async handleEdit(data) {
    await API.updateMenuItem(this.state.currentItemId, data);
    Utils.showSuccess('Блюдо обновлено успешно');

    if (window.MenuTable) {
      MenuTable.updateItemInTable(this.state.currentItemId, data);
    }
  },

  validate() {
    const name = this.form.querySelector('[name="name"]').value.trim();
    const price = this.form.querySelector('[name="price"]').value;
    const category = this.form.querySelector('[name="category"]').value;

    if (!name) {
      Utils.showError('Укажите название блюда');
      return false;
    }

    if (name.length < 3) {
      Utils.showError('Название должно быть не менее 3 символов');
      return false;
    }

    if (!Utils.isValidPrice(price)) {
      Utils.showError('Укажите корректную цену');
      return false;
    }

    if (!category) {
      Utils.showError('Выберите категорию');
      return false;
    }

    return true;
  },
};

window.MenuForm = MenuForm;
window.MenuTable = MenuTable;
