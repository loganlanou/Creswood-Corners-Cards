"use client";

import { useState, useTransition } from "react";
import Link from "next/link";

import { ProductVisual } from "@/components/product-visual";
import { useCart } from "@/components/cart-provider";
import { formatCurrency } from "@/lib/utils";

export function CartView() {
  const { items, subtotalCents, removeItem, updateQuantity, clear } = useCart();
  const [email, setEmail] = useState("");
  const [customerName, setCustomerName] = useState("");
  const [error, setError] = useState("");
  const [isPending, startTransition] = useTransition();

  const handleCheckout = () => {
    setError("");

    startTransition(async () => {
      const response = await fetch("/api/checkout", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email,
          customerName,
          items: items.map((item) => ({
            productId: item.productId,
            quantity: item.quantity,
          })),
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        setError(data.error ?? "Checkout failed.");
        return;
      }

      clear();
      window.location.href = data.url;
    });
  };

  if (items.length === 0) {
    return (
      <section className="rounded-[2rem] border border-dashed border-white/12 bg-white/4 p-10 text-center">
        <p className="text-lg font-semibold text-white">Your cart is empty.</p>
        <p className="mt-2 text-sm text-white/60">
          Add a few football singles and come back when you are ready to check out.
        </p>
        <Link
          href="/shop"
          className="mt-6 inline-flex rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950"
        >
          Browse inventory
        </Link>
      </section>
    );
  }

  return (
    <div className="grid gap-8 lg:grid-cols-[1.35fr_0.65fr]">
      <section className="space-y-4">
        {items.map((item) => (
          <article
            key={item.productId}
            className="grid gap-4 rounded-[2rem] border border-white/10 bg-white/5 p-4 md:grid-cols-[180px_1fr]"
          >
            <ProductVisual title={item.title} accent={item.accent} compact />
            <div className="flex flex-col justify-between gap-4">
              <div>
                <p className="text-xl font-semibold text-white">{item.title}</p>
                <p className="mt-2 max-w-xl text-sm leading-6 text-white/60">{item.summary}</p>
              </div>
              <div className="flex flex-wrap items-center justify-between gap-4">
                <div className="inline-flex items-center rounded-full border border-white/10 bg-[#0a1525] px-2 py-2">
                  <button
                    type="button"
                    onClick={() => updateQuantity(item.productId, item.quantity - 1)}
                    className="h-8 w-8 rounded-full text-white/80 transition hover:bg-white/8"
                  >
                    -
                  </button>
                  <span className="min-w-10 text-center text-sm font-medium">{item.quantity}</span>
                  <button
                    type="button"
                    onClick={() => updateQuantity(item.productId, item.quantity + 1)}
                    className="h-8 w-8 rounded-full text-white/80 transition hover:bg-white/8"
                  >
                    +
                  </button>
                </div>
                <div className="flex items-center gap-4">
                  <p className="text-lg font-semibold text-[#f7b733]">
                    {formatCurrency(item.quantity * item.priceCents)}
                  </p>
                  <button
                    type="button"
                    onClick={() => removeItem(item.productId)}
                    className="text-sm text-white/55 transition hover:text-white"
                  >
                    Remove
                  </button>
                </div>
              </div>
            </div>
          </article>
        ))}
      </section>
      <aside className="rounded-[2rem] border border-white/10 bg-white/5 p-6">
        <p className="text-sm font-semibold uppercase tracking-[0.24em] text-white/55">
          Checkout
        </p>
        <p className="mt-4 text-3xl font-semibold text-white">{formatCurrency(subtotalCents)}</p>
        <p className="mt-2 text-sm text-white/55">
          Stripe Checkout is used automatically when configured. Until then, the cart uses a local
          demo checkout flow so the full store can be tested on localhost.
        </p>
        <div className="mt-6 space-y-3">
          <input
            value={customerName}
            onChange={(event) => setCustomerName(event.target.value)}
            placeholder="Customer name"
            className="h-12 w-full rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm text-white outline-none placeholder:text-white/35"
          />
          <input
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            placeholder="Email address"
            className="h-12 w-full rounded-2xl border border-white/10 bg-[#0a1525] px-4 text-sm text-white outline-none placeholder:text-white/35"
          />
        </div>
        {error ? <p className="mt-4 text-sm text-rose-300">{error}</p> : null}
        <button
          type="button"
          disabled={isPending}
          onClick={handleCheckout}
          className="mt-6 inline-flex w-full items-center justify-center rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-[#ffc95a] disabled:cursor-not-allowed disabled:opacity-60"
        >
          {isPending ? "Processing..." : "Proceed to checkout"}
        </button>
      </aside>
    </div>
  );
}
