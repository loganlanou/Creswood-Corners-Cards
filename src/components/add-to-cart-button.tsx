"use client";

import { ShoppingCart } from "lucide-react";

import { useCart } from "@/components/cart-provider";

export function AddToCartButton({
  item,
  className = "",
}: {
  item: {
    productId: string;
    slug: string;
    title: string;
    summary: string;
    priceCents: number;
    accent?: string | null;
  };
  className?: string;
}) {
  const { addItem } = useCart();

  return (
    <button
      type="button"
      onClick={() => addItem(item)}
      className={`inline-flex items-center justify-center gap-2 rounded-full bg-[#f7b733] px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-[#ffc95a] ${className}`}
    >
      <ShoppingCart className="size-4" />
      Add to cart
    </button>
  );
}
