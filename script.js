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
  "Тушёнка «Радость туриста»": 299,
  "Жаркое из мутантов": 699,
  "Чесночный суп": 1099,
  "Аномальный пирог": 499,
  "Радиационный коктейль": 249,
  "Мутационный чай": 99
};

const cart = [];
const cartList = document.getElementById("cartList");
const cartCount = document.getElementById("cartCount");
const totalPriceEl = document.getElementById("totalPrice");

// ===== ADD TO CART =====
document.querySelectorAll(".add-to-cart").forEach(btn => {
  btn.addEventListener("click", () => {                     
    const title = btn.parentElement.querySelector("h3").textContent;
    cart.push(title);
    updateCart();
  });
});

function updateCart() {
  if (!cartList) return;

  cartList.innerHTML = "";
  let total = 0;

  cart.forEach((item, index) => {
    const price = prices[item] || 0;
    total += price;

    const li = document.createElement("li");
    li.innerHTML = `
      <span>${item}</span>
      <span class="cart-price">${price} ₽</span>
      <button class="remove-item" data-index="${index}">✕</button>
    `;
    cartList.appendChild(li);
  });

  if (cartCount) cartCount.textContent = cart.length;
  if (totalPriceEl) totalPriceEl.textContent = total + " ₽";

  document.querySelectorAll(".remove-item").forEach(btn => {
    btn.onclick = () => {
      cart.splice(btn.dataset.index, 1);
      updateCart();
    };
  });
}

// ===== CART MODAL =====
const cartBtn = document.getElementById("cartBtn");
const cartModal = document.getElementById("cartModal");
const closeCart = document.getElementById("closeCart");
const orderBtn = document.getElementById("orderBtn");

if (cartBtn)
  cartBtn.onclick = () => cartModal.style.display = "flex";

if (closeCart)
  closeCart.onclick = () => cartModal.style.display = "none";

if (orderBtn)
  orderBtn.onclick = () => {
    if (!cart.length) {
      alert("Корзина пуста ☢️");
      return;
    }

    alert("Заказ оформлен! Сталкеры уже выдвинулись 🚶‍♂️");
    cart.length = 0;
    updateCart();
    cartModal.style.display = "none";
  };

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
      alert("⚠️ Номер телефона должен содержать 10 цифр!");
      return;
    }

    const selectedDateTime = new Date(date + 'T' + time);
    if (selectedDateTime < new Date()) {
      alert("⚠️ Нельзя выбрать прошедшее время!");
      return;
    }

    try {
      const response = await fetch("https://qwefsdfsdsg-mdk.hf.space", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },  
        body: JSON.stringify({
          name,
          phone: "+7" + phone,
          date,
          time,
          guests: parseInt(guests)
        })
      });

      if (!response.ok) {
        throw new Error("Ошибка сервера");
      }

      alert("Столик успешно забронирован! ☢️");
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
