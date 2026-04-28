import { cart } from "./cart.js";

// открыть модалку
export function openModal(modal) {
  if (!modal) return;
  modal.classList.add("active");
  modal.style.display = "flex";
}


// закрыть модалку
export function closeModal(modal) {
  if (!modal) return;
  modal.classList.remove("active");
  modal.style.display = "none";
}

// закрытие по клику вне окна
export function bindOverlayClose(modal) {
  if (!modal) return;

  modal.addEventListener("click", (e) => {
    if (e.target === modal) {
      closeModal(modal);
    }
  });
}

document.addEventListener("DOMContentLoaded", () => {
  const btn = document.getElementById("orderBtn");
  const modal = document.getElementById("orderModal");

  if (btn && modal) {
    btn.addEventListener("click", () => {
      modal.style.display = "flex";
    });
  }
})

const modal = document.getElementById("orderModal");
const closeBtn = document.getElementById("closeOrder");

if (closeBtn && modal) {
  closeBtn.addEventListener("click", () => {
    modal.style.display = "none";
  });
}

document.addEventListener("DOMContentLoaded", () => {
  const btn = document.getElementById("orderBtn");
  const modal = document.getElementById("orderModal");

  if (btn && modal) {
    btn.addEventListener("click", () => {

      const isEmpty = Object.values(cart).every(qty => qty === 0);

      if (isEmpty) {
        alert("Корзина пуста ☢️ Добавьте блюда перед заказом");
        return;
      }

      modal.style.display = "flex";
    });
  }
});