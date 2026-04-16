(function () {
  const KEY = "creswood-go-cart";

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

  function addItem(button) {
    const item = {
      product_id: Number(button.dataset.productId),
      slug: button.dataset.slug,
      title: button.dataset.title,
      summary: button.dataset.summary,
      price_cents: Number(button.dataset.price),
      accent: button.dataset.accent || "gold",
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
  const grid = document.getElementById("catalog-grid");
  if (search && grid) {
    search.addEventListener("input", () => {
      const query = search.value.toLowerCase().trim();
      grid.querySelectorAll(".product-card").forEach((card) => {
        const haystack = (card.dataset.search || "").toLowerCase();
        card.style.display = !query || haystack.includes(query) ? "" : "none";
      });
    });
  }

  const cartItems = document.getElementById("cart-items");
  const cartTotal = document.getElementById("cart-total");
  const checkoutForm = document.getElementById("checkout-form");
  const cartError = document.getElementById("cart-error");

  function formatCurrency(cents) {
    return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(cents / 100);
  }

  function renderCart() {
    if (!cartItems) return;
    const items = readCart();
    let total = 0;
    if (items.length === 0) {
      cartItems.innerHTML = '<article class="panel"><h2>Your cart is empty.</h2><p>Add a few cards from the shop to test the checkout flow.</p></article>';
      cartTotal.textContent = "$0.00";
      return;
    }

    cartItems.innerHTML = "";
    items.forEach((item) => {
      total += item.price_cents * item.quantity;
      const row = document.createElement("article");
      row.className = "cart-line";
      row.innerHTML = `
        <div>
          <strong>${item.title}</strong>
          <p>${item.summary}</p>
        </div>
        <div>
          <strong>${formatCurrency(item.price_cents * item.quantity)}</strong>
          <p>Qty ${item.quantity}</p>
        </div>
      `;
      cartItems.appendChild(row);
    });
    cartTotal.textContent = formatCurrency(total);
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
