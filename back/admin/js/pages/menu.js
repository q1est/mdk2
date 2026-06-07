document.addEventListener('DOMContentLoaded', async () => {
  Auth.protectRoute();

  await MenuTable.init('menu-table', 'menu-table-body');
  MenuForm.init('menu-modal', 'menu-form');
  updateUserInfo();

  const logoutButton = document.getElementById('logout-button');
  if (logoutButton) {
    logoutButton.addEventListener('click', (e) => {
      e.preventDefault();
      if (confirm('Вы уверены, что хотите выйти?')) {
        Auth.logout();
      }
    });
  }
});

function updateUserInfo() {
  const user = Auth.getCurrentUser();
  const userNameElement = document.getElementById('user-name');
  const userEmailElement = document.getElementById('user-email');

  if (userNameElement && user) {
    userNameElement.textContent = user.name || user.email || 'Администратор';
  }

  if (userEmailElement && user) {
    userEmailElement.textContent = user.email || '';
  }
}
