import { prices } from "./data.js";

export const cart = {};

export function addToCart(item) {
  cart[item] = (cart[item] || 0) + 1;
}

export function removeFromCart(item) {
  if (!cart[item]) return;

  cart[item]--;

  if (cart[item] <= 0) delete cart[item];
}

export function getCartData() {
  let total = 0;
  let count = 0;

  Object.entries(cart).forEach(([item, qty]) => {
    total += (prices[item] || 0) * qty;
    count += qty;
  });

  return { total, count };
}