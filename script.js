// ===== MENU SWITCH =====
const buttons = document.querySelectorAll("nav button");
const categories = document.querySelectorAll(".category");

buttons.forEach(btn => {
  btn.addEventListener("click", () => {
    buttons.forEach(b => b.classList.remove("active"));
    categories.forEach(c => c.classList.remove("active"));
    btn.classList.add("active");
    document.getElementById(btn.dataset.target).classList.add("active");
  });
});

// ===== PRICES =====
const prices = {
  "Оливье": 350,
  "Кукуруза": 250,
  "Красная икра": 1099,
  "Гороховая каша с мясом": 650,
  "Макароны по-флотски": 600,
  "Пшёнка с мясом": 620,
  "Солянка": 700,
  "Мороженое": 250,
  "Пломбир": 300,
  "Шоколадный пломбир": 320,
  "Эскимо": 280,
  "Брусничная водка": 400,
  "Детское шампанское": 200,
  "Пунш": 300,
  "Шампанское": 350,
  "Гороховый суп": 550,

};

const cart = {};
const cartList = document.getElementById("cartList");
const cartCount = document.getElementById("cartCount");
const totalPriceEl = document.getElementById("totalPrice");

// ===== ADD TO CART =====
document.querySelectorAll(".add-to-cart").forEach(btn => {
  btn.addEventListener("click", () => {
    const title = btn.parentElement.querySelector("h3").textContent;

    if (cart[title]) {
      cart[title]++;
    } else {
      cart[title] = 1;
    }

    updateCart();
  });
});

function updateCart() {
  if (!cartList) return;

  cartList.innerHTML = "";
  let total = 0;
  let count = 0;

  Object.entries(cart).forEach(([item, qty]) => {
    const price = prices[item] || 0;
    total += price * qty;
    count += qty;

    const li = document.createElement("li");

    li.innerHTML = `
      <span>${item}</span>

      <div class="cart-controls">
        <button class="minus" data-item="${item}">−</button>
        <span class="qty">${qty}</span>
        <button class="plus" data-item="${item}">+</button>
      </div>

      <span class="cart-price">${price * qty} ₽</span>
    `;

    cartList.appendChild(li);
  });

  if (cartCount) cartCount.textContent = count;
  if (totalPriceEl) totalPriceEl.textContent = total + " ₽";

  // ➕ добавить
  document.querySelectorAll(".plus").forEach(btn => {
    btn.onclick = () => {
      cart[btn.dataset.item]++;
      updateCart();
    };
  });

  // ➖ убрать
  document.querySelectorAll(".minus").forEach(btn => {
    btn.onclick = () => {
      const item = btn.dataset.item;
      cart[item]--;

      if (cart[item] <= 0) {
        delete cart[item];
      }

      updateCart();
    };
  });
}



// ===== CART MODAL =====
const cartBtn = document.getElementById("cartBtn");
const cartModal = document.getElementById("cartModal");
const closeCart = document.getElementById("closeCart");
const orderBtn = document.getElementById("orderBtn");

// ===== CLEAR CART =====
const clearCartBtn = document.getElementById("clearCart");

if (clearCartBtn) {
  clearCartBtn.addEventListener("click", () => {

    if (Object.keys(cart).length === 0) {
      alert("Корзина уже пуста ☢️");
      return;
    }

    if (confirm("Очистить корзину?")) {
      for (let item in cart) {
      delete cart[item];
      }   // очищаем массив
      updateCart();      // обновляем интерфейс
    }
  });
}
if (cartBtn)
  cartBtn.onclick = () => cartModal.style.display = "flex";

if (closeCart)
  closeCart.onclick = () => cartModal.style.display = "none";

const orderModal = document.getElementById("orderModal");
const closeOrder = document.getElementById("closeOrder");
const orderForm = document.getElementById("orderForm");
const orderPhone = document.getElementById("orderPhone");
const orderTelegram = document.getElementById("orderTelegram"); 

if (orderBtn)
  orderBtn.onclick = () => {
    if (Object.keys(cart).length === 0) {
      alert("☢️ Корзина пуста ☢️");
      return;
    }

    cartModal.style.display = "none";
    orderModal.classList.add("active");
  };

if (closeOrder)
  closeOrder.onclick = () => {
    orderModal.classList.remove("active");
  };

if (orderModal)
  orderModal.addEventListener("click", (e) => {
    if (e.target === orderModal) {
      orderModal.classList.remove("active");
    }
  });


// ===== BOOKING MODAL =====
const bookingModal = document.getElementById("bookingModal");
const ctaBookingBtn = document.getElementById("ctaBookingBtn");
const closeBooking = document.getElementById("closeBooking");

if (ctaBookingBtn)
  ctaBookingBtn.addEventListener("click", () => {
    bookingModal.classList.add("active");
  });

if (closeBooking)
  closeBooking.addEventListener("click", () => {
    bookingModal.classList.remove("active");
  });

if (bookingModal)
  bookingModal.addEventListener("click", (e) => {
    if (e.target === bookingModal) {
      bookingModal.classList.remove("active");
    }
  });

// ===== BOOKING FORM =====
const bookingForm = document.getElementById("bookingForm");
const phoneInput = document.getElementById("phone");
const dateInput = document.getElementById("date");
const timeInput = document.getElementById("time");


if (phoneInput) {
  phoneInput.addEventListener("input", (e) => {
    e.target.value = e.target.value.replace(/\D/g, '').substring(0, 10);
  });
}
if (orderPhone) { 
  orderPhone.addEventListener("input", (e) => {
    e.target.value = e.target.value.replace(/\D/g, '').substring(0, 10);
  });
}


if (dateInput) {
  const today = new Date().toISOString().split('T')[0];
  dateInput.setAttribute('min', today);
}

if (bookingForm) {
  bookingForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const name = document.getElementById("name").value;
    const phone = phoneInput.value;
    const date = dateInput.value;
    const time = timeInput.value;
    const guests = document.getElementById("guests").value;

    if (phone.length !== 10) {
      alert("⚠️ Номер телефона должен содержать 10 цифр! ⚠️");
      return;
    }

    const selectedDateTime = new Date(date + 'T' + time);
    if (selectedDateTime < new Date()) {
      alert("⚠️ Нельзя выбрать прошедшее время! ⚠️");
      return;
    }

    try {
      const response = await fetch("https://qwefsdfsdsg-mdk.hf.space/api/reservations", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },  
        body: JSON.stringify({
          name,
          phone: "+7" + phone,
          date,
          time,
          guests: parseInt(guests),
       
        })
      });

      if (!response.ok) {
        throw new Error("Ошибка сервера");
      }

      alert("☢️ Столик успешно забронирован! ☢️");
      bookingForm.reset();
      bookingModal.classList.remove("active");

    } catch (error) {
      alert("Ошибка при отправке бронирования");
    }
  });
}

// ===== SCROLL =====
const scrollToContacts = document.getElementById("scrollToContacts");
if (scrollToContacts) {
  scrollToContacts.addEventListener("click", () => {
    document.getElementById("contacts").scrollIntoView({
      behavior: "smooth", 
      block: "center"
    });
  });
}

const scrollToBooking = document.getElementById("scrollToBooking");
if (scrollToBooking) {
  scrollToBooking.addEventListener("click", () => {
    document.getElementById("booking").scrollIntoView({
      behavior: "smooth",
      block: "center"
    });
  });
}
if (orderPhone) {
  orderPhone.addEventListener("input", (e) => {
    e.target.value = e.target.value.replace(/\D/g, '').substring(0, 10);
  });
}

if (orderForm) {
  orderForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const name = document.getElementById("orderName").value;
    const phone = orderPhone.value;
    const address = document.getElementById("orderAddress").value;
    const tg = orderTelegram.value.trim();

    if (phone.length !== 10) {
      alert("⚠️ Номер телефона должен содержать 10 цифр! ⚠️");
      return;
    }

    if (!address.trim()) {
      alert("⚠️ Укажите адрес доставки! ⚠️");
      return;
    }

    // ✅ СНАЧАЛА проверка tg
    if (!tg.startsWith("@")) {
      alert("⚠️ Telegram username должен начинаться с @ ⚠️");
      return;
    }

    // подготовим данные для отправки
    const orderItems = Object.entries(cart).map(([name, qty]) => ({
    name,
    qty
    }));
    let total = 0;
    orderItems.forEach(item => {
    total += (prices[item.name] || 0) * item.qty;
    });

    try {
      const response = await fetch("https://qwefsdfsdsg-mdk.hf.space/api/orders", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          name,
          phone: "+7" + phone,
          address,
          telegram: tg,
          items: orderItems,
          total
        })
      });

      if (!response.ok) {
        throw new Error("Ошибка сервера");
      }

      alert("🚶‍♂️ Заказ оформлен! Сталкеры уже выдвинулись 🚶‍♂️");
      for (let item in cart) {
      delete cart[item];
      }
      updateCart();
      orderForm.reset();
      orderModal.classList.remove("active");

    } catch (error) {
      alert("Ошибка при отправке заказа");
      console.error(error);
    }
  });
}

updateCart();
// ===== reveal animation =====
const reveals = document.querySelectorAll(".reveal");

function revealOnScroll() {
  const windowHeight = window.innerHeight;

  reveals.forEach(el => {
    const top = el.getBoundingClientRect().top;

    if (top < windowHeight - 100) {
      el.classList.add("active");
    }
  });
}

window.addEventListener("scroll", revealOnScroll);
window.addEventListener("load", revealOnScroll);