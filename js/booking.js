import {
  openModal,
  closeModal,
  bindOverlayClose
} from "./modal.js";

const bookingModal =
  document.getElementById(
    "bookingModal"
  );

const bookingBtn =
  document.getElementById(
    "ctaBookingBtn"
  );

const closeBooking =
  document.getElementById(
    "closeBooking"
  );

if (bookingBtn && bookingModal) {

  bookingBtn.addEventListener(
    "click",
    () => {

      openModal(
        bookingModal
      );
    }
  );
}

if (closeBooking && bookingModal) {

  closeBooking.addEventListener(
    "click",
    () => {

      closeModal(
        bookingModal
      );
    }
  );
}

bindOverlayClose(
  bookingModal
);