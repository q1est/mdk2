// ===== BOOKING MODAL =====
const bookingModal = document.getElementById("bookingModal");
const ctaBookingBtn = document.getElementById("ctaBookingBtn");
const closeBooking = document.getElementById("closeBooking");

if (ctaBookingBtn) {
  ctaBookingBtn.addEventListener("click", () => {
    bookingModal.classList.add("active");
  });
}

if (closeBooking) {
  closeBooking.addEventListener("click", () => {
    bookingModal.classList.remove("active");
  });
}

if (bookingModal) {
  bookingModal.addEventListener("click", (e) => {
    if (e.target === bookingModal) {
      bookingModal.classList.remove("active");
    }
  });
}

// ===== BOOKING FORM =====
const bookingForm = document.getElementById("bookingForm");
const phoneInput = document.getElementById("phone");
const dateInput = document.getElementById("date");
const timeInput = document.getElementById("time");


// ввод только цифр (телефон)
if (phoneInput) {
  phoneInput.addEventListener("input", (e) => {
    e.target.value = e.target.value.replace(/\D/g, '').substring(0, 10);
  });
}

// запрет прошлых дат
if (dateInput) {
  const today = new Date().toISOString().split('T')[0];
  dateInput.setAttribute('min', today);
}

// отправка формы
if (bookingForm) {
  bookingForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const name = document.getElementById("name").value;
    const phone = phoneInput.value;
    const date = dateInput.value;
    const time = timeInput.value;
    const guests = document.getElementById("guests").value;

    // проверка телефона
    if (phone.length !== 10) {
      alert("⚠️ Номер телефона должен содержать 10 цифр! ⚠️");
      return;
    }

    // запрет прошедшего времени
    const selectedDateTime = new Date(date + 'T' + time);
    if (selectedDateTime < new Date()) {
      alert("⚠️ Нельзя выбрать прошедшее время! ⚠️");
      return;
    }

    // запрет понедельника
    const selectedDate = new Date(date);
    const day = selectedDate.getDay(); // 1 = понедельник

    if (day === 1) {
      alert("⚠️ По понедельникам бар закрыт ⚠️");
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
          guests: parseInt(guests)
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
      console.error(error);
    }
  });
}