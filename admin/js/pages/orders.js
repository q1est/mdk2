document.addEventListener('DOMContentLoaded', async () => {
  Auth.protectRoute();
  await OrdersTable.init('orders-table', 'orders-table-body');
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

  if (userNameElement && user) {
    userNameElement.textContent = user.name || user.email || 'Администратор';
  }
}
