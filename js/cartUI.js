import { cart, addToCart, removeFromCart, getCartData } from "./cart.js";
import { prices } from "./data.js";

const cartList = document.getElementById("cartList");
const cartCount = document.getElementById("cartCount");
const totalPriceEl = document.getElementById("totalPrice");

document.querySelectorAll(".add-to-cart").forEach(btn => {
  btn.addEventListener("click", () => {
    const title = btn.parentElement.querySelector("h3").textContent;
    addToCart(title);
    updateCart();
  });
});

export function updateCart() {
  if (!cartList) return;

  cartList.innerHTML = "";

  Object.entries(cart).forEach(([item, qty]) => {
    const li = document.createElement("li");

    li.innerHTML = `
      <span>${item}</span>
      <div class="cart-controls">
        <button class="minus" data-item="${item}">−</button>
        <span>${qty}</span>
        <button class="plus" data-item="${item}">+</button>
      </div>
      <span>${prices[item] * qty} ₽</span>
    `;

    cartList.appendChild(li);
  });

  const { total, count } = getCartData();

  if (cartCount) cartCount.textContent = count;
  if (totalPriceEl) totalPriceEl.textContent = total + " ₽";

  bindControls();
}

function bindControls() {
  document.querySelectorAll(".plus").forEach(btn => {
    btn.onclick = () => {
      addToCart(btn.dataset.item);
      updateCart();
    };
  });

  document.querySelectorAll(".minus").forEach(btn => {
    btn.onclick = () => {
      removeFromCart(btn.dataset.item);
      updateCart();
    };
  });
}

updateCart();