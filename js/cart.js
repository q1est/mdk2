export const cart = {};

export function addToCart(item) {

    const id = item.id;

    if (!cart[id]) {

        cart[id] = {
            ...item,
            qty: 1,
        };

    } else {

        cart[id].qty++;
    }
}

export function removeFromCart(id) {

    if (!cart[id]) return;

    cart[id].qty--;

    if (cart[id].qty <= 0) {
        delete cart[id];
    }
}

export function clearCart() {

    Object.keys(cart).forEach(id => {
        delete cart[id];
    });
}

export function getCartData() {

    let total = 0;
    let count = 0;

    Object.values(cart).forEach(item => {

        total += item.price * item.qty;

        count += item.qty;
    });

    return {
        total,
        count,
    };
}