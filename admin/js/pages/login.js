document.addEventListener('DOMContentLoaded', async () => {
  if (Auth.isAuthenticated()) {
    window.location.href = '/admin/menu.html';
    return;
  }

  const form = document.getElementById('login-form');
  const submitButton = form?.querySelector('button[type="submit"]');

  if (!form) {
    console.error('Login form not found');
    return;
  }

  form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const email = form.querySelector('[name="email"]').value.trim();
    const password = form.querySelector('[name="password"]').value;

    if (!email) {
      Utils.showError('Укажите email');
      return;
    }

    if (!Utils.isValidEmail(email)) {
      Utils.showError('Укажите корректный email');
      return;
    }

    if (!password) {
      Utils.showError('Укажите пароль');
      return;
    }

    try {
      Utils.setButtonLoading(submitButton, true);
      
      const result = await Auth.login(email, password);
      
      Utils.showSuccess('Вход успешен!');
      
      setTimeout(() => {
        window.location.href = '/admin/menu.html';
      }, 500);
    } catch (error) {
      console.error('Login error:', error);
      Utils.showError('Ошибка входа: ' + error.message);
    } finally {
      Utils.setButtonLoading(submitButton, false);
    }
  });

  const demoButton = document.getElementById('demo-login-button');
  if (demoButton) {
    demoButton.addEventListener('click', () => {
      form.querySelector('[name="email"]').value = 'admin@mdk2.com';
      form.querySelector('[name="password"]').value = 'password123';
      form.dispatchEvent(new Event('submit'));
    });
  }
});
