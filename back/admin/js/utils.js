const Utils = {
  showNotification(message, type = 'info', duration = 5000) {
    let container = document.getElementById('notifications');
    if (!container) {
      container = document.createElement('div');
      container.id = 'notifications';
      container.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 9999;
      `;
      document.body.appendChild(container);
    }

    const notification = document.createElement('div');
    notification.className = `alert alert-${type} alert-dismissible fade show`;
    notification.setAttribute('role', 'alert');
    notification.innerHTML = `
      ${message}
      <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    notification.style.cssText = 'min-width: 300px; margin-bottom: 10px;';

    container.appendChild(notification);

    if (duration > 0) {
      setTimeout(() => notification.remove(), duration);
    }

    return notification;
  },

  showError(message) {
    this.showNotification(message, 'danger');
  },

  showSuccess(message) {
    this.showNotification(message, 'success');
  },

  setButtonLoading(button, isLoading = true) {
    if (isLoading) {
      button.disabled = true;
      button.dataset.originalText = button.textContent;
      button.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Загрузка...';
    } else {
      button.disabled = false;
      button.textContent = button.dataset.originalText || 'Отправить';
    }
  },

  isValidEmail(email) {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  },

  isValidPrice(price) {
    return !isNaN(parseFloat(price)) && parseFloat(price) > 0;
  },

  formatPrice(price) {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB',
    }).format(price);
  },

  formatDate(date) {
    return new Intl.DateTimeFormat('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(new Date(date));
  },

  clearForm(form) {
    if (form && form.reset) {
      form.reset();
    }
  },

  getFormData(form) {
    const formData = new FormData(form);
    return Object.fromEntries(formData);
  },

  debounce(func, delay = CONFIG.SEARCH_DEBOUNCE_MS) {
    let timeoutId;
    return function (...args) {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => func.apply(this, args), delay);
    };
  },

  confirmDelete(itemName = 'элемент') {
    return confirm(`Вы уверены, что хотите удалить ${itemName}? Это действие необратимо.`);
  },
};

const UI = {
  createMenuRow(item) {
    const row = document.createElement('tr');
    row.dataset.itemId = item.id;

    const statusBadge = item.available 
      ? '<span class="badge bg-success">Доступно</span>'
      : '<span class="badge bg-secondary">Недоступно</span>';

    const categoryLabel = CONFIG.MENU_CATEGORIES.find(c => c.value === item.category)?.label || item.category;

    row.innerHTML = `
      <td>${Utils.escapeHtml(item.id)}</td>
      <td>
        <div class="d-flex align-items-center">
          ${item.image_url ? `<img src="${Utils.escapeHtml(item.image_url)}" alt="Image" class="img-thumbnail me-2" style="width: 40px; height: 40px; object-fit: cover;">` : '<div class="img-thumbnail me-2" style="width: 40px; height: 40px;"></div>'}
          <strong>${Utils.escapeHtml(item.name)}</strong>
        </div>
      </td>
      <td>${Utils.formatPrice(item.price)}</td>
      <td>${Utils.escapeHtml(categoryLabel)}</td>
      <td>${statusBadge}</td>
      <td>
        <div class="btn-group btn-group-sm" role="group">
          <button type="button" class="btn btn-outline-primary btn-edit" data-id="${item.id}" title="Редактировать">
            <i class="icon-edit"></i>
          </button>
          <button type="button" class="btn btn-outline-danger btn-delete" data-id="${item.id}" title="Удалить">
            <i class="icon-trash"></i>
          </button>
        </div>
      </td>
    `;

    return row;
  },

  showTableSkeleton(tableBody, rows = 5) {
    tableBody.innerHTML = '';
    for (let i = 0; i < rows; i++) {
      const row = document.createElement('tr');
      row.innerHTML = `
        <td><div class="placeholder col-6"></div></td>
        <td><div class="placeholder col-12"></div></td>
        <td><div class="placeholder col-4"></div></td>
        <td><div class="placeholder col-6"></div></td>
        <td><div class="placeholder col-5"></div></td>
        <td><div class="placeholder col-8"></div></td>
      `;
      tableBody.appendChild(row);
    }
  },

  showEmptyState(tableBody, message = 'Блюд не найдено') {
    tableBody.innerHTML = `
      <tr>
        <td colspan="6" class="text-center text-muted py-5">
          <p>${message}</p>
        </td>
      </tr>
    `;
  },
};

Utils.escapeHtml = function(text) {
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  return String(text).replace(/[&<>"']/g, m => map[m]);
};
