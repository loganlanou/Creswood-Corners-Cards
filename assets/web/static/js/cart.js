(function () {
  const KEY = "creswood-go-cart";
  const COUPON_KEY = "creswood-go-coupon";
  const COUPONS = {
    WELCOME10: 0.1,
    LIVE5: 0.05,
  };

  function readCart() {
    try {
      return JSON.parse(localStorage.getItem(KEY) || "[]");
    } catch (_) {
      return [];
    }
  }

  function writeCart(items) {
    localStorage.setItem(KEY, JSON.stringify(items));
  }

  function readCoupon() {
    return (localStorage.getItem(COUPON_KEY) || "").toUpperCase().trim();
  }

  function writeCoupon(code) {
    if (code) {
      localStorage.setItem(COUPON_KEY, code.toUpperCase().trim());
    } else {
      localStorage.removeItem(COUPON_KEY);
    }
  }

  function escapeHTML(value) {
    return String(value || "")
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#039;");
  }

  function addItem(button) {
    const item = {
      product_id: Number(button.dataset.productId),
      slug: button.dataset.slug,
      title: button.dataset.title,
      summary: button.dataset.summary,
      price_cents: Number(button.dataset.price),
      accent: button.dataset.accent || "gold",
      team: button.dataset.team || "",
      quantity: 1,
    };
    const items = readCart();
    const existing = items.find((entry) => entry.product_id === item.product_id);
    if (existing) {
      existing.quantity += 1;
    } else {
      items.push(item);
    }
    writeCart(items);
    button.textContent = "Added";
    setTimeout(() => { button.textContent = button.classList.contains("small") ? "Add" : "Add to Cart"; }, 1000);
  }

  document.querySelectorAll(".add-to-cart").forEach((button) => {
    button.addEventListener("click", function () {
      addItem(button);
    });
  });

  const search = document.getElementById("catalog-search");
  const teamFilter = document.getElementById("catalog-team");
  const grid = document.getElementById("catalog-grid");
  if (grid) {
    const applyCatalogFilters = () => {
      const query = search ? search.value.toLowerCase().trim() : "";
      const team = teamFilter ? teamFilter.value.toLowerCase().trim() : "";
      grid.querySelectorAll(".product-card").forEach((card) => {
        const haystack = (card.dataset.search || "").toLowerCase();
        const cardTeam = (card.dataset.team || "").toLowerCase();
        const matchesSearch = !query || haystack.includes(query);
        const matchesTeam = !team || cardTeam === team;
        card.style.display = matchesSearch && matchesTeam ? "" : "none";
      });
    };
    if (search) search.addEventListener("input", applyCatalogFilters);
    if (teamFilter) teamFilter.addEventListener("change", applyCatalogFilters);
  }

  const cartItems = document.getElementById("cart-items");
  const cartSubtotal = document.getElementById("cart-subtotal");
  const cartDiscount = document.getElementById("cart-discount");
  const cartTax = document.getElementById("cart-tax");
  const cartTotal = document.getElementById("cart-total");
  const checkoutForm = document.getElementById("checkout-form");
  const cartError = document.getElementById("cart-error");
  const couponInput = document.getElementById("coupon-code");
  const applyCoupon = document.getElementById("apply-coupon");
  const couponMessage = document.getElementById("coupon-message");

  function formatCurrency(cents) {
    return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(cents / 100);
  }

  function renderCart() {
    if (!cartItems) return;
    const items = readCart();
    let subtotal = 0;
    if (items.length === 0) {
      cartItems.innerHTML = '<article class="panel"><h2>Your cart is empty.</h2><p>Add a few cards from the shop to test the checkout flow.</p></article>';
      updateTotals(0);
      return;
    }

    cartItems.innerHTML = "";
    items.forEach((item) => {
      subtotal += item.price_cents * item.quantity;
      const row = document.createElement("article");
      row.className = "cart-line";
      row.dataset.productId = item.product_id;
      row.innerHTML = `
        <div>
          <strong>${escapeHTML(item.title)}</strong>
          <p>${escapeHTML(item.summary)}</p>
          ${item.team ? `<span class="tag">${escapeHTML(item.team)}</span>` : ""}
        </div>
        <div class="cart-line-controls">
          <strong>${formatCurrency(item.price_cents * item.quantity)}</strong>
          <label>
            <span>Qty</span>
            <input class="cart-qty" type="number" min="1" step="1" value="${item.quantity}" data-product-id="${item.product_id}">
          </label>
          <button class="linkish cart-remove" type="button" data-product-id="${item.product_id}">Remove</button>
        </div>
      `;
      cartItems.appendChild(row);
    });
    updateTotals(subtotal);
  }

  function updateTotals(subtotal) {
    const code = readCoupon();
    const discountRate = COUPONS[code] || 0;
    const discount = Math.round(subtotal * discountRate);
    const tax = 0;
    const total = Math.max(subtotal - discount + tax, 0);

    if (cartSubtotal) cartSubtotal.textContent = formatCurrency(subtotal);
    if (cartDiscount) cartDiscount.textContent = discount ? `-${formatCurrency(discount)}` : "$0.00";
    if (cartTax) cartTax.textContent = formatCurrency(tax);
    if (cartTotal) cartTotal.textContent = formatCurrency(total);
    if (couponInput) couponInput.value = code;
  }

  if (cartItems) {
    cartItems.addEventListener("change", (event) => {
      if (!event.target.classList.contains("cart-qty")) return;
      const productID = Number(event.target.dataset.productId);
      const quantity = Math.max(Number(event.target.value) || 1, 1);
      const items = readCart().map((item) => (
        item.product_id === productID ? { ...item, quantity } : item
      ));
      writeCart(items);
      renderCart();
    });

    cartItems.addEventListener("click", (event) => {
      if (!event.target.classList.contains("cart-remove")) return;
      const productID = Number(event.target.dataset.productId);
      writeCart(readCart().filter((item) => item.product_id !== productID));
      renderCart();
    });
  }

  if (applyCoupon && couponInput) {
    applyCoupon.addEventListener("click", () => {
      const code = couponInput.value.toUpperCase().trim();
      if (!code) {
        writeCoupon("");
        couponMessage.textContent = "Coupon removed.";
        renderCart();
        return;
      }
      if (!COUPONS[code]) {
        couponMessage.textContent = "That demo code is not active.";
        return;
      }
      writeCoupon(code);
      couponMessage.textContent = `${code} applied.`;
      renderCart();
    });
  }

  if (checkoutForm) {
    checkoutForm.addEventListener("submit", async (event) => {
      event.preventDefault();
      cartError.textContent = "";
      const form = new FormData(checkoutForm);
      const items = readCart().map((item) => ({
        product_id: item.product_id,
        quantity: item.quantity,
      }));
      const response = await fetch("/checkout", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          customer_name: form.get("customer_name"),
          email: form.get("email"),
          items,
          coupon_code: readCoupon(),
        }),
      });

      if (!response.ok) {
        cartError.textContent = "Checkout failed. Make sure your email and cart items are valid.";
        return;
      }

      const data = await response.json();
      writeCart([]);
      window.location.href = data.redirect;
    });
  }

  renderCart();
})();
